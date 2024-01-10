package client

import "time"

type ProfileCreateRequest struct {
	IdentityID string    `json:"identityId"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Height     uint      `json:"height"`
	Weight     uint      `json:"weight"`
	Birthday   time.Time `json:"birthDate"`
}

type ProfileUpdateRequest struct {
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Height     uint      `json:"height"`
	Weight     uint      `json:"weight"`
	Birthday   time.Time `json:"birthDate"`
}

type ContactInformationCreateRequest struct {
	ContactType  string `json:"contactType"`
	ContactValue string `json:"contactValue"`
}

type AddressCreateRequest struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
}

type ProfileResponse struct {
	ID         string    `json:"id"`
	UpdatedAt  time.Time `json:"updatedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Height     uint      `json:"height"`
	Weight     uint      `json:"weight"`
	Birthday   time.Time `json:"birthDate"`
}

type ContactInformationResponse struct {
	ID           string    `json:"id"`
	ProfileID    string    `json:"profileId"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	ContactType  string    `json:"contactType"`
	ContactValue string    `json:"contactValue"`
}

type AddressResponse struct {
	ID           string    `json:"id"`
	ProfileID    string    `json:"profileId"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
	AddressLine1 string    `json:"addressLine1"`
	AddressLine2 string    `json:"addressLine2"`
	AddressLine3 string    `json:"addressLine3"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Zip          string    `json:"zip"`
	Country      string    `json:"country"`
}
