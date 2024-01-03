package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/admin/auth/client"
	"github.com/anchamber-studios/hevonen/services/admin/auth/db"
)

type AuthService interface {
	GetAuthorizations(ctx context.Context) ([]client.Authorization, error)
	GetAuthorizationsForService(ctx context.Context, serviceId string) ([]client.Authorization, error)
	GetGroups(ctx context.Context) ([]client.Group, error)
}

type AuthServiceImpl struct {
	AuthRepo db.AuthRepository
}

func (a *AuthServiceImpl) GetAuthorizations(ctx context.Context) ([]client.Authorization, error) {
	return a.AuthRepo.GetAuthorizations(ctx)
}

func (a *AuthServiceImpl) GetAuthorizationsForService(ctx context.Context, serviceId string) ([]client.Authorization, error) {
	return a.AuthRepo.GetAuthorizationsForService(ctx, serviceId)
}

func (a *AuthServiceImpl) GetGroups(ctx context.Context) ([]client.Group, error) {
	return a.AuthRepo.GetGroups(ctx)
}
