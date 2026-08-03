package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/katharinasick/clubrizer/internal/apperrors"
)

// notificationTimeout bounds the fire-and-forget email goroutines below, which run on a
// fresh background context (not the request context, which is cancelled once the handler
// returns).
const notificationTimeout = 30 * time.Second

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

// decisionRecipients is who to notify after a batch decision. accountEmails are the account
// holders whose own status changed — one email each. kidParents are the parents of decided
// kids whose account was NOT part of the batch, i.e. a later-added kid on an already-approved
// account: those get a kid-specific email. Kids decided alongside their own account (a
// whole-family batch) are intentionally absent — the single account email covers the family,
// so a family approval sends exactly one email.
type decisionRecipients struct {
	accountEmails []string
	kidParents    []kidParentNotification
}

type kidParentNotification struct {
	parentEmail string
	kidName     string
}

// ApproveApprovals approves the given accounts and/or kids in one transaction, then emails
// the affected account holders (and the parents of any individually-approved kids).
func (s *Service) ApproveApprovals(ctx context.Context, req ApprovalDecisionRequest) error {
	return s.decide(ctx, req, StatusApproved)
}

// RejectApprovals rejects the given accounts and/or kids in one transaction, then emails
// the affected account holders (and the parents of any individually-rejected kids).
func (s *Service) RejectApprovals(ctx context.Context, req ApprovalDecisionRequest) error {
	return s.decide(ctx, req, StatusRejected)
}

func (s *Service) decide(ctx context.Context, req ApprovalDecisionRequest, decision Status) error {
	// Authorization is enforced at the route (requirePermission with users.approve), so
	// this only validates the payload.
	if len(req.UserIDs) == 0 && len(req.KidIDs) == 0 {
		return apperrors.NewBadRequest("select at least one account or kid", nil)
	}

	recipients, err := s.store.decideApprovals(ctx, req.UserIDs, req.KidIDs, decision)
	if err != nil {
		return err
	}

	s.sendDecisionEmails(recipients, decision)
	return nil
}

// sendDecisionEmails notifies the affected account holders and kid parents off the request
// path. The decision is already committed, so delivery is best-effort and fire-and-forget:
// it must neither block nor fail the admin's response, and errors are only logged. Runs on a
// fresh background context so it survives the handler returning.
func (s *Service) sendDecisionEmails(recipients decisionRecipients, decision Status) {
	if len(recipients.accountEmails) == 0 && len(recipients.kidParents) == 0 {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), notificationTimeout)
		defer cancel()

		for _, e := range recipients.accountEmails {
			var mailErr error
			if decision == StatusApproved {
				mailErr = s.email.SendAccountApproved(ctx, e)
			} else {
				mailErr = s.email.SendAccountRejected(ctx, e)
			}
			if mailErr != nil {
				s.log.Error("failed to send account decision email", "decision", decision, "error", mailErr)
			}
		}

		for _, k := range recipients.kidParents {
			var mailErr error
			if decision == StatusApproved {
				mailErr = s.email.SendKidApproved(ctx, k.parentEmail, k.kidName)
			} else {
				mailErr = s.email.SendKidRejected(ctx, k.parentEmail, k.kidName)
			}
			if mailErr != nil {
				s.log.Error("failed to send kid decision email", "decision", decision, "error", mailErr)
			}
		}
	}()
}

// notifyAdminsOfNewApprovalRequest lets admins know a freshly-submitted account is waiting
// in the queue. Fire-and-forget and best-effort: it runs off the request path on a fresh
// background context so it never blocks or fails onboarding — the account is already pending
// and will show up in the queue regardless.
func (s *Service) notifyAdminsOfNewApprovalRequest(applicant *User) {
	name := applicantDisplayName(applicant)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), notificationTimeout)
		defer cancel()

		admins, err := s.store.getAdminEmails(ctx)
		if err != nil {
			s.log.Error("failed to load admin emails for approval notification", "error", err)
			return
		}
		if len(admins) == 0 {
			return
		}
		if err := s.email.SendNewApprovalRequest(ctx, admins, name); err != nil {
			s.log.Error("failed to send new-approval-request email", "error", err)
		}
	}()
}

// applicantDisplayName picks the friendliest name available for an applicant, falling back
// to the email when no name has been set yet.
func applicantDisplayName(u *User) string {
	switch {
	case u.NickName != nil && *u.NickName != "":
		return *u.NickName
	case u.GivenName != nil && *u.GivenName != "":
		return *u.GivenName
	default:
		return u.Email
	}
}
