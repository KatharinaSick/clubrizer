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

// ReleaseDeviceRequest identifies a browser by its subscription. Auth is the subscription's
// secret, which only that browser and our database know — it proves the caller is on that
// device, whoever the subscription currently belongs to.
type ReleaseDeviceRequest struct {
	Endpoint string `json:"endpoint" validate:"required,max=2048"`
	Auth     string `json:"auth" validate:"required,max=255"`
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
