package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
)

// MemberService represents a service for managing members in a club.
type MemberService struct {
	repo db.MemberRepo
}

// NewMemberService creates a new instance of MemberService.
func NewMemberService(repo db.MemberRepo) *MemberService {
	return &MemberService{repo: repo}
}

// List returns a list of all members for a club.
func (s *MemberService) List(ctx context.Context, cId string) ([]types.Member, error) {
	return s.repo.List(ctx)
}

// Create creates a new member.
func (s *MemberService) Create(ctx context.Context, member types.MemberCreate) (string, error) {
	return s.repo.Create(ctx, member)
}
