package users

import "context"

// MemberKid is an approved kid shown nested under its parent in the members roster.
type MemberKid struct {
	GivenName  *string `json:"givenName"`
	FamilyName *string `json:"familyName"`
	Picture    *string `json:"picture"`
}

// Member is one approved account in the club roster: the account holder plus their assigned
// roles (the implicit 'member' role is excluded — see getRolesByUserID) and their approved
// kids. SelfParticipates distinguishes a participating member from a guardian who only
// manages kids.
type Member struct {
	Email            string       `json:"email"`
	GivenName        *string      `json:"givenName"`
	FamilyName       *string      `json:"familyName"`
	NickName         *string      `json:"nickName"`
	Picture          *string      `json:"picture"`
	SelfParticipates bool         `json:"selfParticipates"`
	Roles            []*Role      `json:"roles"`
	Kids             []*MemberKid `json:"kids"`
}

// ListMembers returns the full roster of approved accounts with their roles and approved
// kids. The route is gated on the users.approve permission (see requirePermission), so this
// trusts the caller.
func (s *Service) ListMembers(ctx context.Context) ([]*Member, error) {
	members, err := s.store.listMembers(ctx)
	if err != nil {
		return nil, err
	}
	if members == nil {
		members = []*Member{}
	}
	return members, nil
}
