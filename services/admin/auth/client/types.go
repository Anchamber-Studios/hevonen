package client

import (
	"time"

	"github.com/google/uuid"
)

type Authorization struct {
	ID          uuid.UUID `json:"id"`
	ServiceId   uuid.UUID `json:"serviceId"`
	ServiceName string    `json:"serviceName"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Service struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	AuthorizationEndpoint string    `json:"authorizationEndpoint"`
	UpdatedAt             time.Time `json:"updatedAt"`
	CreatedAt             time.Time `json:"createdAt"`
}

type ServiceKey struct {
	ID        uuid.UUID `json:"id"`
	ServiceId uuid.UUID `json:"service_id"`
	SecretKey string    `json:"key"`
	ValidTill time.Time `json:"validTill"`
	CreatedAt time.Time `json:"createdAt"`
}

type Group struct {
	ID            uuid.UUID       `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
	Users         []User          `json:"users"`
	Authorization []Authorization `json:"authorizations"`
	Parent        *Group          `json:"parent"`
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
