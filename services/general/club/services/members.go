package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
)

type MemberService struct {
	repo db.MemberRepo
}

func NewMemberService(repo db.MemberRepo) *MemberService {
	return &MemberService{repo: repo}
}

func (s *MemberService) List(ctx context.Context) ([]types.Member, error) {
	return s.repo.List(ctx)
}

func (s *MemberService) ListForClub(ctx context.Context, cId string) ([]types.Member, error) {
	return s.repo.ListForClub(ctx, cId)
}

func (s *MemberService) Create(ctx context.Context, member types.MemberCreate) (string, error) {
	return s.repo.Create(ctx, member)
}
