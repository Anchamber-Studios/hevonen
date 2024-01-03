package db

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/admin/auth/client"
	"github.com/google/uuid"
)

const (
	ServiceUUID = "018cd07b-83c8-7931-8110-9cb2b9c62d4a"
)

type AuthRepository interface {
	GetAuthorizations(ctx context.Context) ([]client.Authorization, error)
	GetAuthorizationsForService(ctx context.Context, serviceId string) ([]client.Authorization, error)
	GetGroups(ctx context.Context) ([]client.Group, error)
}

type AuthRepositoryPostgre struct{}

func (a *AuthRepositoryPostgre) GetAuthorizations(ctx context.Context) ([]client.Authorization, error) {
	return systemAuthorizations(), nil
}

func (a *AuthRepositoryPostgre) GetAuthorizationsForService(ctx context.Context, serviceId string) ([]client.Authorization, error) {
	auths, err := a.GetAuthorizations(ctx)
	if err != nil {
		return nil, err
	}
	filtered := []client.Authorization{}
	for _, auth := range auths {
		if auth.ServiceId == uuid.Must(uuid.Parse(serviceId)) {
			filtered = append(filtered, auth)
		}
	}
	return filtered, nil
}

func (a *AuthRepositoryPostgre) GetGroups(ctx context.Context) ([]client.Group, error) {
	return systemGroups(), nil
}

func NewAuthRepositoryPostgre() AuthRepository {
	return &AuthRepositoryPostgre{}
}

func systemAuthorizations() []client.Authorization {
	serviceID := uuid.Must(uuid.Parse(ServiceUUID))
	return []client.Authorization{
		{
			ID:          uuid.Must(uuid.Parse("018cd07b-2e89-7daa-a35c-e64b0c8373cb")),
			ServiceId:   serviceID,
			ServiceName: "admin/auth",
			Name:        "users/get",
		},
		{
			ID:          uuid.Must(uuid.Parse("018cd080-8957-7416-b1b9-56e98c65cd54")),
			ServiceId:   serviceID,
			ServiceName: "admin/auth",
			Name:        "users/edit/:orgId",
		},
		{
			ID:          uuid.Must(uuid.Parse("018cd080-6157-734f-80cc-a52bac568e19")),
			ServiceId:   serviceID,
			ServiceName: "admin/auth",
			Name:        "services/get",
		},
		{
			ID:          uuid.Must(uuid.Parse("018cd080-7505-7403-ab03-2279d0078552")),
			ServiceId:   serviceID,
			ServiceName: "admin/auth",
			Name:        "services/edit",
		},
	}
}

func systemGroups() []client.Group {
	return []client.Group{
		{
			ID:            uuid.Must(uuid.Parse("018cd07e-4f63-7de9-b5bc-eddd8302ce06")),
			Name:          "admin",
			Description:   "Admin group for the auth service",
			Users:         []client.User{},
			Authorization: systemAuthorizations(),
			Parent:        nil,
		},
	}
}
