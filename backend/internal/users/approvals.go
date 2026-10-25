package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/katharinasick/clubrizer/internal/apperrors"
)

// ApprovalKid is a pending kid awaiting an admin decision, shown under its parent account
// in the approval queue.
type ApprovalKid struct {
	ID         uuid.UUID `json:"id"`
	GivenName  *string   `json:"givenName"`
	FamilyName *string   `json:"familyName"`
	Picture    *string   `json:"picture"`
}

// ApprovalRequest is one item in the approval queue: an account together with any of its
// pending kids. Status is the account holder's own status — StatusPending when the adult
// themselves is awaiting approval, or StatusApproved when the adult is already in and only
// newly-added kids need a decision. SelfParticipates distinguishes a guardian ("only my
// kids") from a participating member. The frontend batches an item's ids into one
// approve/reject call — and, guided by Status, must not include an already-approved UserID
// in an approve batch.
type ApprovalRequest struct {
	UserID           uuid.UUID      `json:"userId"`
	Email            string         `json:"email"`
	GivenName        *string        `json:"givenName"`
	FamilyName       *string        `json:"familyName"`
	NickName         *string        `json:"nickName"`
	Picture          *string        `json:"picture"`
	Status           Status         `json:"status"`
	SelfParticipates bool           `json:"selfParticipates"`
	PendingKids      []*ApprovalKid `json:"pendingKids"`
}

// ApprovalDecisionRequest is the payload for approving or rejecting one or many accounts
// and/or kids in a single all-or-nothing batch. The frontend collects a family's pending
// ids into one request so a whole family is decided with one click; a per-member decision
// is just a smaller batch.
type ApprovalDecisionRequest struct {
	UserIDs []uuid.UUID `json:"userIds"`
	KidIDs  []uuid.UUID `json:"kidIds"`
}

// ListApprovals returns the approval queue: every account that is itself awaiting approval,
// plus every approved account that has newly-added kids awaiting approval. The route is
// gated on the users.approve permission (see requirePermission), so this trusts the caller.
func (s *Service) ListApprovals(ctx context.Context) ([]*ApprovalRequest, error) {
	requests, err := s.store.listApprovals(ctx)
	if err != nil {
		return nil, err
	}
	if requests == nil {
		requests = []*ApprovalRequest{}
	}
	return requests, nil
}

// ApproveApprovals approves the given accounts and/or kids in one transaction, then tells
// the listener (which notifies the account holders and the parents of individually-approved
// kids).
func (s *Service) ApproveApprovals(ctx context.Context, req ApprovalDecisionRequest) error {
	return s.decide(ctx, req, StatusApproved)
}

// RejectApprovals rejects the given accounts and/or kids in one transaction, then tells the
// listener (which notifies the account holders and the parents of individually-rejected
// kids).
func (s *Service) RejectApprovals(ctx context.Context, req ApprovalDecisionRequest) error {
	return s.decide(ctx, req, StatusRejected)
}

func (s *Service) decide(ctx context.Context, req ApprovalDecisionRequest, decision Status) error {
	// Authorization is enforced at the route (requirePermission with users.approve), so
	// this only validates the payload.
	if len(req.UserIDs) == 0 && len(req.KidIDs) == 0 {
		return apperrors.NewBadRequest("select at least one account or kid", nil)
	}

	accounts, kids, err := s.store.decideApprovals(ctx, req.UserIDs, req.KidIDs, decision)
	if err != nil {
		return err
	}

	s.listener.ApprovalsDecided(ctx, ApprovalsDecided{Decision: decision, Accounts: accounts, Kids: kids})
	return nil
}
