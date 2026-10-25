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

// This file turns what happened in the events and users packages into push notifications:
// who gets told, with which text and link. The person who did something is never notified
// about it themselves.

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

// ApprovalRequested tells everyone who can approve members. It names nobody: applicants only
// enter their name after approval, and their email address must not show up on lock screens.
func (s *Service) ApprovalRequested(ctx context.Context, _ users.ApprovalRequested) {
	approvers, err := s.rbac.GetUserIDsWithPermission(ctx, rbac.PermissionUsersApprove)
	if err != nil {
		s.log.Warn("failed to load approvers for push notification", "error", err)
		return
	}
	s.notify(ctx, approvers, message{Type: typeMemberRequest, URL: "/manage-members"})
}

// ApprovalsDecided tells each decided account, and the parent of each kid decided on its own.
func (s *Service) ApprovalsDecided(ctx context.Context, decided users.ApprovalsDecided) {
	accountType, kidType := typeAccountApproved, typeKidApproved
	if decided.Decision == users.StatusRejected {
		accountType, kidType = typeAccountRejected, typeKidRejected
	}

	s.notify(ctx, decided.AccountIDs, message{Type: accountType, URL: "/"})
	for _, kid := range decided.Kids {
		s.notify(ctx, []uuid.UUID{kid.ParentID}, message{
			Type:   kidType,
			Params: map[string]string{"name": kid.Name},
			URL:    "/my-kids",
		})
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
