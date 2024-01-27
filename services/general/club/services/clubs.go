package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
)

type ClubService struct {
	repo db.ClubRepo
}

func NewClubService(repo db.ClubRepo) *ClubService {
	return &ClubService{repo: repo}
}

func (s *ClubService) List(ctx context.Context) ([]types.Club, error) {
	return s.repo.List(ctx)
}

func (s *ClubService) ListForIdentity(ctx context.Context, identity string) ([]types.ClubMember, error) {
	return s.repo.ListForIdentity(ctx, identity)
}

func (s *ClubService) Create(ctx context.Context, club types.ClubCreate) (string, error) {
	return s.repo.Create(ctx, club)
}
