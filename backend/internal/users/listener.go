package users

import (
	"context"

	"github.com/google/uuid"
)

// listener is told about approval changes once they have happened, e.g. to send
// notifications. It only gets the facts — deciding who to tell what is up to the listener.
// Implementations are best-effort and never fail the action that triggered them.
type listener interface {
	ApprovalRequested(ctx context.Context, requested ApprovalRequested)
	ApprovalsDecided(ctx context.Context, decided ApprovalsDecided)
}

// ApprovalRequested describes an account that was just submitted for approval.
type ApprovalRequested struct {
	Applicant *User
}

// ApprovalsDecided describes one approve/reject batch. Accounts are the accounts whose own
// status changed; kids decided together with their account are covered by it. Kids lists
// kids decided on their own, i.e. added later to an already-approved account.
type ApprovalsDecided struct {
	Decision Status
	Accounts []DecidedAccount
	Kids     []DecidedKid
}

type DecidedAccount struct {
	UserID uuid.UUID
	Email  string
}

type DecidedKid struct {
	ParentID    uuid.UUID
	ParentEmail string
	Name        string
}
