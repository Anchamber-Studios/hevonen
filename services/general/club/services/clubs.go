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

// List all clubs
func (s *ClubService) List(ctx context.Context) ([]types.Club, error) {
	return s.repo.List(ctx)
}

// List all clubs for a given identity
func (s *ClubService) ListForIdentity(ctx context.Context, identity string) ([]types.ClubContact, error) {
	return s.repo.ListForIdentity(ctx, identity)
}

// Create a new club
func (s *ClubService) Create(ctx context.Context, club types.ClubCreate) (string, error) {
	return s.repo.Create(ctx, club)
}

// Create a new club with the current user as admin contact
func (s *ClubService) CreateWithAdminContact(ctx context.Context, club types.ClubCreate, admin types.ContactCreate) (string, error) {
	return s.repo.CreateWithAdminContact(ctx, club, admin)
}

// Delete a club
func (s *ClubService) Delete(ctx context.Context, identity string, id string) error {
	return s.repo.Delete(ctx, identity, id)
}
