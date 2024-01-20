package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/anchamber-studios/hevonen/services/club/db"
)

type ClubService struct {
	repo db.ClubRepo
}

func NewClubService(repo db.ClubRepo) *ClubService {
	return &ClubService{repo: repo}
}

func (s *ClubService) List(ctx context.Context) ([]client.Club, error) {
	return s.repo.List(ctx)
}

func (s *ClubService) ListForIdentity(ctx context.Context, identity string) ([]client.ClubMember, error) {
	return s.repo.ListForIdentity(ctx, identity)
}
