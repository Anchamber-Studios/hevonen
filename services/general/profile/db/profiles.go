package db

import (
	"context"
	"errors"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/sqids/sqids-go"
)

type ProfileRepo interface {
	GetByIdentityID(ctx context.Context, identityID string) (client.ProfileResponse, error)
	Create(ctx context.Context, profile client.ProfileCreateRequest) (string, error)
	Delete(ctx context.Context, id string) error
	UpdateByIdentityID(ctx context.Context, id string, profile client.ProfileUpdateRequest) error

	ListAddresses(ctx context.Context, profileId string) ([]client.AddressResponse, error)
	CreateAddress(ctx context.Context, profileId string, address client.AddressCreateRequest) (string, error)
	DeleteAddress(ctx context.Context, profileId string, addressId string) error

	ListContactInfo(ctx context.Context, profileId string) ([]client.ContactInformationResponse, error)
	CreateContactInfo(ctx context.Context, profileId string, contactInfo client.ContactInformationCreateRequest) (string, error)
	DeleteContactInfo(ctx context.Context, profileId string, contactInfoId string) error
}

type ProfileRepoPostgre struct {
	Logger       echo.Logger
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
}

const (
	IdOffset uint64 = 2345678901
)

func (r *ProfileRepoPostgre) GetByIdentityID(ctx context.Context, identityID string) (client.ProfileResponse, error) {
	var profile client.ProfileResponse
	err := r.DB.QueryRow(ctx, `
			SELECT id, first_name, middle_name, last_name, height, weight, birthday
			FROM profile.profiles WHERE identity_id = $1;
			`, identityID).
		Scan(&profile.ID, &profile.FirstName, &profile.MiddleName, &profile.LastName, &profile.Height, &profile.Weight, &profile.Birthday)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return client.ProfileResponse{}, lib.NewNotFoundError().WithMessage("profile not found")
		}
		r.Logger.Errorf("unable to get profile for identity '%s': %v\n", identityID, err)
		return client.ProfileResponse{}, err
	}
	r.Logger.Infof("profile for identity '%s' retrieved", identityID)
	return profile, nil
}

func (r *ProfileRepoPostgre) Create(ctx context.Context, profile client.ProfileCreateRequest) (string, error) {
	var id string
	err := r.DB.QueryRow(ctx, `
			INSERT INTO profile.profiles (identity_id, first_name, middle_name, last_name, height, weight, birthday)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id;
			`, profile.IdentityID, profile.FirstName, profile.MiddleName, profile.LastName, profile.Height, profile.Weight, profile.Birthday).
		Scan(&id)
	if err != nil {
		r.Logger.Errorf("unable to create profile: %v\n", err)
		return "", err
	}
	r.Logger.Infof("profile '%s' created", id)
	return id, nil
}

func (r *ProfileRepoPostgre) Delete(ctx context.Context, id string) error {
	cId := r.IdConversion.Decode(id)
	_, err := r.DB.Exec(ctx, "DELETE FROM profile.profiles WHERE identity_id = $1;", cId[0])
	if err != nil {
		return err
	}
	r.Logger.Info("profile '%s' deleted", id)
	return nil
}

func (r *ProfileRepoPostgre) UpdateByIdentityID(ctx context.Context, id string, profile client.ProfileUpdateRequest) error {
	_, err := r.DB.Exec(ctx, `
			UPDATE profile.profiles
			SET first_name = $1, middle_name = $2, last_name = $3, height = $4, weight = $5, birthday = $6
			WHERE identity_id = $7;
			`, profile.FirstName, profile.MiddleName, profile.LastName, profile.Height, profile.Weight, profile.Birthday, id)
	if err != nil {
		r.Logger.Errorf("unable to update profile '%s': %v\n", id, err)
		return err
	}
	r.Logger.Infof("profile '%s' updated", id)
	return nil
}

func (r *ProfileRepoPostgre) ListAddresses(ctx context.Context, profileId string) ([]client.AddressResponse, error) {
	cId := r.IdConversion.Decode(profileId)
	rows, err := r.DB.Query(ctx, `
			SELECT id, profile_id, address_line_1, address_line_2, address_line_3, city, state, zip_code, country
			FROM profile.addresses WHERE profile_id = $1;
			`, cId[0])
	if err != nil {
		r.Logger.Errorf("unable to list addresses for profile '%s': %v\n", profileId, err)
		return nil, err
	}
	defer rows.Close()
	var addresses []client.AddressResponse
	for rows.Next() {
		var address client.AddressResponse
		err := rows.Scan(&address.ID, &address.ProfileID, &address.AddressLine1, &address.AddressLine2, &address.AddressLine3, &address.City, &address.State, &address.Zip, &address.Country)
		if err != nil {
			r.Logger.Errorf("unable to scan address for profile '%s': %v\n", profileId, err)
			return nil, err
		}
		addresses = append(addresses, address)
	}
	r.Logger.Infof("addresses for profile '%s' retrieved", profileId)
	return addresses, nil
}

func (r *ProfileRepoPostgre) CreateAddress(ctx context.Context, profileId string, address client.AddressCreateRequest) (string, error) {
	var id string
	cId := r.IdConversion.Decode(profileId)
	err := r.DB.QueryRow(ctx, `
			INSERT INTO profile.addresses (profile_id, address_line_1, address_line_2, address_line_3, city, state, zip_code, country)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id;
			`, cId[0], address.AddressLine1, address.AddressLine2, address.AddressLine3, address.City, address.State, address.Zip, address.Country).
		Scan(&id)
	if err != nil {
		r.Logger.Errorf("unable to create address for profile '%s': %v\n", profileId, err)
		return "", err
	}
	r.Logger.Infof("address '%s' created for profile '%s'", id, profileId)
	return id, nil
}

func (r *ProfileRepoPostgre) DeleteAddress(ctx context.Context, profileId string, addressId string) error {
	cId := r.IdConversion.Decode(profileId)
	cAddressId := r.IdConversion.Decode(addressId)
	_, err := r.DB.Exec(ctx, "DELETE FROM profile.addresses WHERE profile_id = $1 AND id = $2;", cId[0], cAddressId[0])
	if err != nil {
		r.Logger.Errorf("unable to delete address '%s' for profile '%s': %v\n", addressId, profileId, err)
		return err
	}
	r.Logger.Infof("address '%s' deleted for profile '%s'", addressId, profileId)
	return nil
}

func (r *ProfileRepoPostgre) ListContactInfo(ctx context.Context, profileId string) ([]client.ContactInformationResponse, error) {
	cId := r.IdConversion.Decode(profileId)
	rows, err := r.DB.Query(ctx, `
			SELECT id, profile_id, contact_type, contact_value
			FROM profile.contact_information WHERE profile_id = $1;
			`, cId[0])
	if err != nil {
		r.Logger.Errorf("unable to list contact information for profile '%s': %v\n", profileId, err)
		return nil, err
	}
	defer rows.Close()
	var contactInfo []client.ContactInformationResponse
	for rows.Next() {
		var ci client.ContactInformationResponse
		err := rows.Scan(&ci.ID, &ci.ProfileID, &ci.ContactType, &ci.ContactValue)
		if err != nil {
			r.Logger.Errorf("unable to scan contact information for profile '%s': %v\n", profileId, err)
			return nil, err
		}
		contactInfo = append(contactInfo, ci)
	}
	r.Logger.Infof("contact information for profile '%s' retrieved", profileId)
	return contactInfo, nil
}

func (r *ProfileRepoPostgre) CreateContactInfo(ctx context.Context, profileId string, contactInfo client.ContactInformationCreateRequest) (string, error) {
	var id string
	cId := r.IdConversion.Decode(profileId)
	err := r.DB.QueryRow(ctx, `
			INSERT INTO profile.contact_information (profile_id, contact_type, contact_value)
			VALUES ($1, $2, $3)
			RETURNING id;
			`, cId[0], contactInfo.ContactType, contactInfo.ContactValue).
		Scan(&id)
	if err != nil {
		r.Logger.Errorf("unable to create contact information for profile '%s': %v\n", profileId, err)
		return "", err
	}
	r.Logger.Infof("contact information '%s' created for profile '%s'", id, profileId)
	return id, nil
}

func (r *ProfileRepoPostgre) DeleteContactInfo(ctx context.Context, profileId string, contactInfoId string) error {
	cId := r.IdConversion.Decode(profileId)
	cContactInfoId := r.IdConversion.Decode(contactInfoId)
	_, err := r.DB.Exec(ctx, "DELETE FROM profile.contact_information WHERE profile_id = $1 AND id = $2;", cId[0], cContactInfoId[0])
	if err != nil {
		r.Logger.Errorf("unable to delete contact information '%s' for profile '%s': %v\n", contactInfoId, profileId, err)
		return err
	}
	r.Logger.Infof("contact information '%s' deleted for profile '%s'", contactInfoId, profileId)
	return nil
}
