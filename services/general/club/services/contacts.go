package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/db"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
)

// ContactService represents a service for managing contacts in a club.
type ContactService struct {
	repo db.ContactRepo
}

// NewContactService creates a new instance of ContactService.
func NewContactService(repo db.ContactRepo) *ContactService {
	return &ContactService{repo: repo}
}

// List returns a list of all contacts for a club.
func (s *ContactService) List(ctx context.Context, cId string) ([]types.Contact, error) {
	return s.repo.ListForClub(ctx, cId)
}

// Create creates a new contact.
func (s *ContactService) Create(ctx context.Context, contact types.ContactCreate) (string, error) {
	return s.repo.Create(ctx, contact)
}
