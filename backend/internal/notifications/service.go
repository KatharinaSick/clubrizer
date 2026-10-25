package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/users"
)

// allowedPushHosts are the push services of the browsers we support. The backend POSTs to
// whatever endpoint a client registers, so accepting arbitrary URLs would let anyone make
// the server call internal or third-party addresses. Subdomains of these hosts are allowed
// too (e.g. Edge uses <region>.notify.windows.com).
var allowedPushHosts = []string{
	"fcm.googleapis.com",        // Chrome, Edge on Android, Samsung Internet, Opera, ...
	"push.services.mozilla.com", // Firefox
	"push.apple.com",            // Safari (macOS and iOS home-screen apps)
	"notify.windows.com",        // Edge on Windows
	"android.googleapis.com",    // older Chrome versions
}

// pushTTL is how long a push service keeps a message for an offline device.
const pushTTL = 24 * time.Hour

// notifyTimeout bounds how long Notify may hold up the request that triggered it.
const notifyTimeout = 10 * time.Second

// notifyConcurrency is how many users Notify sends to at the same time.
const notifyConcurrency = 10

// Service is responsible for registering devices and sending Web Push notifications. It
// listens to the events and users packages (see listeners.go) and decides who to notify.
type Service struct {
	log        app.Logger
	cfg        *app.Config
	store      *store
	rbac       rbacService
	httpClient *http.Client
}

type rbacService interface {
	GetUserIDsWithPermission(ctx context.Context, permissionKey string) ([]uuid.UUID, error)
}

func NewService(log app.Logger, cfg *app.Config, conn *pgxpool.Pool, rbac rbacService) *Service {
	if cfg.WebPush.PublicKey == "" || cfg.WebPush.PrivateKey == "" {
		log.Warn("VAPID keys are not configured, push notifications will not work. Generate them with `go run ./cmd/vapidkeys`.")
	}

	store := newStore(log, cfg, conn)
	return &Service{log, cfg, store, rbac, &http.Client{Timeout: 10 * time.Second}}
}

func (s *Service) GetVapidPublicKey(_ context.Context) (*VapidPublicKeyResponse, error) {
	if s.cfg.WebPush.PublicKey == "" {
		return nil, errors.New("push notifications are not configured")
	}
	return &VapidPublicKeyResponse{PublicKey: s.cfg.WebPush.PublicKey}, nil
}

func (s *Service) Subscribe(ctx context.Context, req SubscribeRequest) error {
	if !isAllowedEndpoint(req.Endpoint) {
		return apperrors.NewBadRequest("unsupported push service", nil)
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	return s.store.upsertSubscription(ctx, userId, req.Endpoint, req.Keys.P256dh, req.Keys.Auth)
}

// GetSubscriptionStatus tells a browser whether its subscription is still registered for the
// caller. The browser alone can't know: the device may have been turned off from another
// device ("turn off on all devices"), or someone else may have signed in on it since.
func (s *Service) GetSubscriptionStatus(ctx context.Context, req SubscriptionStatusRequest) (*SubscriptionStatusResponse, error) {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	subscribed, err := s.store.hasSubscription(ctx, userId, req.Endpoint)
	if err != nil {
		return nil, err
	}
	return &SubscriptionStatusResponse{Subscribed: subscribed}, nil
}

func (s *Service) Unsubscribe(ctx context.Context, req UnsubscribeRequest) error {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	return s.store.deleteSubscription(ctx, userId, req.Endpoint)
}

func (s *Service) UnsubscribeAll(ctx context.Context) error {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	return s.store.deleteAllSubscriptions(ctx, userId)
}

// SendTest sends a test notification to all of the caller's devices.
func (s *Service) SendTest(ctx context.Context) error {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	sent, err := s.sendToUser(ctx, userId, message{Type: typeTest, URL: "/notifications"})
	if err != nil {
		return err
	}
	if sent == 0 {
		return apperrors.NewNotFound("no device could be reached")
	}
	return nil
}

// notify sends the message to every device of the given users. It is best-effort and never
// fails the caller: the action that triggered it (a comment, an approval, ...) has already
// happened, so errors are only logged.
//
// It runs inside the request rather than in a goroutine that outlives it: on Cloud Run with
// request-based billing the CPU is throttled once the response is sent, so background work
// can stall or get lost. Making this truly async is tracked in
// https://github.com/KatharinaSick/clubrizer/issues/53. The context is detached from the
// request so a client hanging up doesn't cancel the sends, and bounded by notifyTimeout.
func (s *Service) notify(ctx context.Context, userIds []uuid.UUID, msg message) {
	if len(userIds) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), notifyTimeout)
	defer cancel()

	var wg sync.WaitGroup
	limit := make(chan struct{}, notifyConcurrency)
	for _, userId := range dedupe(userIds) {
		wg.Add(1)
		limit <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-limit }()
			if _, err := s.sendToUser(ctx, userId, msg); err != nil {
				s.log.Warn("failed to send push notification", "user", userId, "type", msg.Type, "error", err)
			}
		}()
	}
	wg.Wait()
}

// sendToUser delivers the message to every device of the user and returns how many devices
// accepted it. Devices the push service reports as gone are removed; other failures are
// logged and skipped so one broken device doesn't block the others.
func (s *Service) sendToUser(ctx context.Context, userId uuid.UUID, msg message) (int, error) {
	subs, err := s.store.getSubscriptionsByUser(ctx, userId)
	if err != nil {
		return 0, err
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to encode push message: %s", err.Error()))
	}

	// Send one device after another: webpush-go v1.4.0 writes into the payload's spare
	// capacity, so sharing payload across goroutines races. Fixed upstream (#83) but not
	// released yet — only parallelize once a release includes it. (notify sending to several
	// users in parallel is fine: each sendToUser call encodes its own payload.)
	sent := 0
	for _, sub := range subs {
		if s.send(ctx, sub, payload) {
			sent++
		}
	}
	return sent, nil
}

func (s *Service) send(ctx context.Context, sub *Subscription, payload []byte) bool {
	resp, err := webpush.SendNotificationWithContext(ctx, payload, &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys:     webpush.Keys{P256dh: sub.P256dh, Auth: sub.Auth},
	}, &webpush.Options{
		HTTPClient:      s.httpClient,
		Subscriber:      s.cfg.WebPush.Subject,
		VAPIDPublicKey:  s.cfg.WebPush.PublicKey,
		VAPIDPrivateKey: s.cfg.WebPush.PrivateKey,
		TTL:             int(pushTTL.Seconds()),
		Urgency:         webpush.UrgencyNormal,
	})
	if err != nil {
		s.log.Warn("failed to send push notification", "subscription", sub.ID, "error", err)
		return false
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone:
		// The browser unsubscribed or the subscription expired — it will never work again.
		if err := s.store.deleteSubscriptionById(ctx, sub.ID); err != nil {
			s.log.Warn("failed to remove expired push subscription", "subscription", sub.ID, "error", err)
		}
		return false
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		if err := s.store.touchSubscription(ctx, sub.ID); err != nil {
			s.log.Warn("failed to update push subscription", "subscription", sub.ID, "error", err)
		}
		return true
	default:
		s.log.Warn("push service rejected notification", "subscription", sub.ID, "status", resp.StatusCode)
		return false
	}
}

func isAllowedEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Port() != "" || u.User != nil {
		return false
	}

	host := strings.ToLower(u.Hostname())
	for _, allowed := range allowedPushHosts {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}

func dedupe(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]bool, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}
