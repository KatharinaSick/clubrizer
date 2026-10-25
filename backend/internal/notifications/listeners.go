package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/katharinasick/clubrizer/internal/events"
	"github.com/katharinasick/clubrizer/internal/rbac"
	"github.com/katharinasick/clubrizer/internal/users"
)

// This file turns what happened in the events and users packages into push notifications
// and emails: who gets told, with which text and link. The person who did something is never
// notified about it themselves.
//
// Emails, like push notifications, are sent inside the request rather than in a goroutine —
// see notify for why, and https://github.com/KatharinaSick/clubrizer/issues/53 for making
// both async.

func (s *Service) EventCancelled(ctx context.Context, change events.EventChanged) {
	s.notifyAttendees(ctx, change, typeEventCancelled, eventUrl(change.Event.ID))
}

func (s *Service) EventUncancelled(ctx context.Context, change events.EventChanged) {
	s.notifyAttendees(ctx, change, typeEventUncancelled, eventUrl(change.Event.ID))
}

func (s *Service) EventDeleted(ctx context.Context, change events.EventChanged) {
	s.notifyAttendees(ctx, change, typeEventDeleted, "/events")
}

// CommentCreated tells the event's creator and everyone who said yes.
func (s *Service) CommentCreated(ctx context.Context, created events.CommentCreated) {
	recipients := append(created.Attendees, created.Event.CreatedBy)
	s.notify(ctx, excluding(recipients, created.AuthorID), message{
		Type: typeEventComment,
		Params: map[string]string{
			"name":    created.AuthorName,
			"title":   created.Event.Title,
			"comment": created.Comment.Body,
		},
		URL: eventUrl(created.Event.ID),
	})
}

// ApprovalRequested tells everyone who can approve members, by push and email. The push
// notification names nobody: applicants only enter their name after approval, and their
// email address must not show up on lock screens. The email names them (falling back to the
// email address), as approvers need to know who is asking.
func (s *Service) ApprovalRequested(ctx context.Context, requested users.ApprovalRequested) {
	approvers, err := s.rbac.GetUserIDsWithPermission(ctx, rbac.PermissionUsersApprove)
	if err != nil {
		s.log.Error("failed to load approvers for approval request notification", "error", err)
		return
	}
	if len(approvers) == 0 {
		return
	}

	s.notify(ctx, approvers, message{Type: typeMemberRequest, URL: "/manage-members"})
	s.sendEmails(ctx, func(ctx context.Context) {
		emails, err := s.store.getEmails(ctx, approvers)
		if err != nil {
			s.log.Error("failed to load approver emails for approval request notification", "error", err)
			return
		}
		if err := s.email.SendNewApprovalRequest(ctx, emails, applicantName(requested.Applicant)); err != nil {
			s.log.Error("failed to send new-approval-request email", "error", err)
		}
	})
}

// ApprovalsDecided tells each decided account, and the parent of each kid decided on its own,
// by push and email. Kids decided together with their account are covered by the account's
// message, so a whole-family decision sends one of each.
func (s *Service) ApprovalsDecided(ctx context.Context, decided users.ApprovalsDecided) {
	approved := decided.Decision == users.StatusApproved
	accountType, kidType := typeAccountApproved, typeKidApproved
	if !approved {
		accountType, kidType = typeAccountRejected, typeKidRejected
	}

	accountIds := make([]uuid.UUID, 0, len(decided.Accounts))
	for _, a := range decided.Accounts {
		accountIds = append(accountIds, a.UserID)
	}
	s.notify(ctx, accountIds, message{Type: accountType, URL: "/"})
	for _, kid := range decided.Kids {
		s.notify(ctx, []uuid.UUID{kid.ParentID}, message{
			Type:   kidType,
			Params: map[string]string{"name": kid.Name},
			URL:    "/my-kids",
		})
	}

	s.sendEmails(ctx, func(ctx context.Context) {
		for _, a := range decided.Accounts {
			send := s.email.SendAccountApproved
			if !approved {
				send = s.email.SendAccountRejected
			}
			if err := send(ctx, a.Email); err != nil {
				s.log.Error("failed to send account decision email", "decision", decided.Decision, "error", err)
			}
		}
		for _, kid := range decided.Kids {
			send := s.email.SendKidApproved
			if !approved {
				send = s.email.SendKidRejected
			}
			if err := send(ctx, kid.ParentEmail, kid.Name); err != nil {
				s.log.Error("failed to send kid decision email", "decision", decided.Decision, "error", err)
			}
		}
	})
}

// sendEmails runs send on a context detached from the request and bounded by notifyTimeout.
// Errors are logged inside send — the action that triggered it has already happened.
func (s *Service) sendEmails(ctx context.Context, send func(ctx context.Context)) {
	ctx, cancel := detached(ctx)
	defer cancel()
	send(ctx)
}

// applicantName picks the friendliest name available for an applicant, falling back to the
// email when no name has been set yet.
func applicantName(u *users.User) string {
	switch {
	case u.NickName != nil && *u.NickName != "":
		return *u.NickName
	case u.GivenName != nil && *u.GivenName != "":
		return *u.GivenName
	default:
		return u.Email
	}
}

func (s *Service) notifyAttendees(ctx context.Context, change events.EventChanged, msgType string, url string) {
	s.notify(ctx, excluding(change.Attendees, change.ActorID), message{
		Type: msgType,
		Params: map[string]string{
			"title": change.Event.Title,
			"date":  change.Event.StartTime.Format(time.RFC3339),
		},
		URL: url,
	})
}

func eventUrl(id uuid.UUID) string {
	return fmt.Sprintf("/events/%s", id)
}

func excluding(ids []uuid.UUID, exclude uuid.UUID) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if id != exclude {
			out = append(out, id)
		}
	}
	return out
}
