package notifications

import (
	"github.com/google/uuid"
)

// Subscription is a stored Web Push subscription of one browser/device.
type Subscription struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Endpoint string
	P256dh   string
	Auth     string
}

// SubscribeRequest mirrors the browser's PushSubscription.toJSON() shape.
type SubscribeRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url,max=2048"`
	Keys     struct {
		P256dh string `json:"p256dh" validate:"required,max=255"`
		Auth   string `json:"auth" validate:"required,max=255"`
	} `json:"keys"`
}

type UnsubscribeRequest struct {
	Endpoint string `json:"endpoint" validate:"required,max=2048"`
}

type SubscriptionStatusRequest struct {
	Endpoint string `json:"endpoint" validate:"required,max=2048"`
}

type SubscriptionStatusResponse struct {
	Subscribed bool `json:"subscribed"`
}

type VapidPublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
}

// Message is the payload delivered to the service worker. It only describes what happened —
// the service worker turns Type into a translated text, so all user-facing wording stays in
// the frontend's i18n.
type Message struct {
	Type string `json:"type"`
	// URL is the in-app path opened when the notification is tapped.
	URL string `json:"url,omitempty"`
}

const (
	MessageTypeTest = "test"
)
