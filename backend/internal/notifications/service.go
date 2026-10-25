package notifications

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/rbac"
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

// notifyTimeout bounds how long notify (push) and sendEmails may each hold up the request
// that triggered them.
const notifyTimeout = 10 * time.Second

// notifyConcurrency is how many devices notify sends to at the same time.
const notifyConcurrency = 10

// Service is responsible for registering devices and sending Web Push notifications. It
// listens to the events and users packages (see listeners.go) and decides who to notify.
type Service struct {
	log        app.Logger
	cfg        *app.Config
	store      *store
	rbac       rbacService
	email      emailSender
	httpClient *http.Client
}

type rbacService interface {
	IsAuthorized(ctx context.Context, userID uuid.UUID, permissionKey string) (bool, error)
	GetUserIDsWithPermission(ctx context.Context, permissionKey string) ([]uuid.UUID, error)
}

// emailSender is the part of the email client used for notification emails.
type emailSender interface {
	SendAccountApproved(ctx context.Context, to string) error
	SendAccountRejected(ctx context.Context, to string) error
	SendKidApproved(ctx context.Context, to, kidName string) error
	SendKidRejected(ctx context.Context, to, kidName string) error
	SendNewApprovalRequest(ctx context.Context, to []string, applicantName string) error
}

func NewService(log app.Logger, cfg *app.Config, conn *pgxpool.Pool, rbac rbacService, email emailSender) *Service {
	if cfg.WebPush.PublicKey == "" || cfg.WebPush.PrivateKey == "" {
		log.Warn("VAPID keys are not configured, push notifications will not work. Generate them with `go run ./cmd/vapidkeys`.")
	}

	store := newStore(log, cfg, conn)
	return &Service{log, cfg, store, rbac, email, &http.Client{Timeout: 10 * time.Second}}
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
// caller. The browser alone can't know: someone else may have signed in on it since.
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

// ReleaseDevice removes this device's subscription whoever it belongs to. The app calls it on
// logout — also when the session has already expired — and when a different person signs in
// on the device, so the previous person's notifications stop arriving there. Needs no
// session: the subscription's secret proves the caller is on that device.
func (s *Service) ReleaseDevice(ctx context.Context, req ReleaseDeviceRequest) error {
	return s.store.deleteSubscriptionOfDevice(ctx, req.Endpoint, req.Auth)
}

// SendTest sends a test notification to all of the caller's devices and reports how many
// accepted it. None is not an error: the app tells the user in its own words.
func (s *Service) SendTest(ctx context.Context) (*SendTestResponse, error) {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	subs, err := s.store.getSubscriptionsByUsers(ctx, []uuid.UUID{userId})
	if err != nil {
		return nil, err
	}
	delivered, err := s.sendToDevices(ctx, subs, message{Type: typeTest, URL: "/notifications"})
	if err != nil {
		return nil, err
	}
	return &SendTestResponse{Delivered: delivered}, nil
}

// GetSettings lists the notification settings screen's switches for the caller. Member
// requests only show up for people who can approve members, as nobody else gets them.
func (s *Service) GetSettings(ctx context.Context) ([]*Setting, error) {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	categories, err := s.availableCategories(ctx, userId)
	if err != nil {
		return nil, err
	}
	optOuts, err := s.store.getOptOuts(ctx, userId)
	if err != nil {
		return nil, err
	}

	settings := make([]*Setting, 0, len(categories))
	for _, c := range categories {
		settings = append(settings, &Setting{Category: c, Enabled: !optOuts[c]})
	}
	return settings, nil
}

// UpdateSetting switches one category on or off for the caller. Only categories listed by
// GetSettings can be changed — checked here because it depends on the request body.
func (s *Service) UpdateSetting(ctx context.Context, req UpdateSettingRequest) error {
	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	categories, err := s.availableCategories(ctx, userId)
	if err != nil {
		return err
	}
	if !slices.Contains(categories, req.Category) {
		return apperrors.NewBadRequest("unknown notification category", nil)
	}
	return s.store.setOptOut(ctx, userId, req.Category, !*req.Enabled)
}

func (s *Service) availableCategories(ctx context.Context, userId uuid.UUID) ([]string, error) {
	categories := []string{categoryNewEvents, categoryResponses, categoryComments, categoryEventChanges}
	canApprove, err := s.rbac.IsAuthorized(ctx, userId, rbac.PermissionUsersApprove)
	if err != nil {
		return nil, err
	}
	if canApprove {
		categories = append(categories, categoryMemberRequests)
	}
	return categories, nil
}

// notify sends the message to every device of the given users, except those who switched its
// category off on the notification settings screen. It is best-effort and never
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

	ctx, cancel := detached(ctx)
	defer cancel()

	userIds, err := s.withoutOptedOut(ctx, typeCategories[msg.Type], userIds)
	if err != nil {
		// Rather send nothing than notify people who asked not to be.
		s.log.Error("failed to load notification settings", "type", msg.Type, "error", err)
		return
	}
	if len(userIds) == 0 {
		return
	}

	subs, err := s.store.getSubscriptionsByUsers(ctx, userIds)
	if err != nil {
		s.log.Error("failed to load push subscriptions", "type", msg.Type, "error", err)
		return
	}
	if _, err := s.sendToDevices(ctx, subs, msg); err != nil {
		s.log.Error("failed to send push notifications", "type", msg.Type, "error", err)
	}
}

// sendToDevices delivers the message to the given devices and returns how many accepted it.
// Devices the push service reports as gone are removed; other failures are logged and
// skipped so one broken device doesn't block the others.
func (s *Service) sendToDevices(ctx context.Context, subs []*Subscription, msg message) (int, error) {
	payload, err := encode(msg)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to encode push message: %s", err.Error()))
	}

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		delivered []uuid.UUID
	)
	limit := make(chan struct{}, notifyConcurrency)
	for _, sub := range subs {
		wg.Add(1)
		limit <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-limit }()
			// Each send gets its own copy: webpush-go v1.4.0 writes into the payload's spare
			// capacity, so sharing one payload across goroutines races. Fixed upstream (#83)
			// but not released yet.
			if s.send(ctx, sub, bytes.Clone(payload)) {
				mu.Lock()
				delivered = append(delivered, sub.ID)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if err := s.store.touchSubscriptions(ctx, delivered); err != nil {
		s.log.Warn("failed to update push subscriptions", "error", err)
	}
	return len(delivered), nil
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
		return true
	default:
		s.log.Warn("push service rejected notification", "subscription", sub.ID, "status", resp.StatusCode)
		return false
	}
}

// withoutOptedOut returns the users who didn't switch the category off. An empty category
// means the message is always sent.
func (s *Service) withoutOptedOut(ctx context.Context, category string, userIds []uuid.UUID) ([]uuid.UUID, error) {
	if category == "" || len(userIds) == 0 {
		return userIds, nil
	}
	optedOut, err := s.store.getOptedOutUserIds(ctx, category, userIds)
	if err != nil {
		return nil, err
	}
	// On a copy: callers may still need the full list (e.g. for emails).
	return slices.DeleteFunc(slices.Clone(userIds), func(id uuid.UUID) bool { return optedOut[id] }), nil
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

// detached returns a context that outlives the request (a client hanging up must not cancel
// sends for an action that already happened), bounded by notifyTimeout.
func detached(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), notifyTimeout)
}
