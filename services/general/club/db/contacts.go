package db

import (
	"context"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/jackc/pgx/v5"
	"github.com/sqids/sqids-go"
)

type ContactRepo interface {
	List(ctx context.Context) ([]types.Contact, error)
	ListForClub(ctx context.Context, clubIdEncoded string) ([]types.Contact, error)
	Create(ctx context.Context, member types.ContactCreate) (string, error)
	Get(ctx context.Context, memberIdEncoded string) (types.Contact, error)
}

type ContactRepoPostgre struct {
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
}

const (
	IdOffset uint64 = 1234567890
)

func (r *ContactRepoPostgre) List(ctx context.Context) ([]types.Contact, error) {
	rows, err := r.DB.Query(ctx, `
	SELECT id, first_name, middle_name, last_name, email, phone, height, weight FROM clubs.contacts;
	`)
	if err != nil {
		return nil, err
	}

	members := make([]types.Contact, 0)
	for rows.Next() {
		member := types.Contact{}
		var id uint64
		err := rows.Scan(&id, &member.FirstName, &member.MiddleName, &member.LastName, &member.Email, &member.Phone, &member.Height, &member.Weight)
		if err != nil {
			return nil, err
		}
		cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
		if err != nil {
			return nil, err
		}
		member.ID = cId
		members = append(members, member)
	}
	return members, nil
}

func (r *ContactRepoPostgre) ListForClub(ctx context.Context, clubIdEncoded string) ([]types.Contact, error) {
	clubId := r.IdConversion.Decode(clubIdEncoded)[0]
	rows, err := r.DB.Query(ctx, `
		SELECT id, 
			COALESCE(first_name, '') as first_name, 
			COALESCE(middle_name, '') as middle_name,
			COALESCE(last_name, '') as last_name,
			COALESCE(email, '') as email,
			COALESCE(phone, '') as phone,
			height, weight 
		FROM clubs.contacts 
		WHERE club_id = $1;
		`, clubId)
	if err != nil {
		return nil, err
	}

	members := make([]types.Contact, 0)
	for rows.Next() {
		member := types.Contact{}
		var id uint64
		err := rows.Scan(&id, &member.FirstName, &member.MiddleName, &member.LastName, &member.Email, &member.Phone, &member.Height, &member.Weight)
		if err != nil {
			return nil, err
		}
		cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
		if err != nil {
			return nil, err
		}
		member.ID = cId
		members = append(members, member)
	}
	return members, nil
}

func (r *ContactRepoPostgre) Create(ctx context.Context, member types.ContactCreate) (string, error) {
	var id uint64
	err := r.DB.QueryRow(ctx, `
	INSERT INTO clubs.members (identity_id, club_id, email) 
	VALUES ($1, $2, $3) 
	RETURNING id;
	`, member.IdentityID, member.ClubID, member.Email).Scan(&id)
	if err != nil {
		return "", err
	}
	cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
	if err != nil {
		return "", err
	}
	return cId, nil
}

func (r *ContactRepoPostgre) Get(ctx context.Context, memberIdEncoded string) (types.Contact, error) {
	memberId := r.IdConversion.Decode(memberIdEncoded)[0]
	var id uint64
	var member types.Contact
	r.DB.QueryRow(ctx, "SELECT id, first_name, last_name, email, phone FROM clubs.members WHERE id = $1;", memberId).Scan(&id, &member.FirstName, &member.LastName, &member.Email, &member.Phone)
	if id == 0 {
		return member, lib.ErrNotFound
	}
	member.ID = memberIdEncoded
	return member, nil
}

type ContactRepoMock struct {
	Contacts []types.Contact
}

func (r *ContactRepoMock) List(ctx context.Context) ([]types.Contact, error) {
	return r.Contacts, nil
}

func (r *ContactRepoMock) ListForClub(ctx context.Context, clubIdEncoded string) ([]types.Contact, error) {
	return r.Contacts, nil
}

func (r *ContactRepoMock) Create(ctx context.Context, member types.ContactCreate) (string, error) {
	for _, m := range r.Contacts {
		if m.Email == member.Email {
			return "", lib.ErrAlreadyExists
		}
	}
	id := fmt.Sprintf("%d", len(r.Contacts)+1)
	r.Contacts = append(r.Contacts, types.Contact{
		ID:     id,
		ClubID: member.ClubID,
		Email:  member.Email,
	})
	return id, nil
}

func (r *ContactRepoMock) Get(ctx context.Context, memberIdEncoded string) (types.Contact, error) {
	for _, m := range r.Contacts {
		if m.ID == memberIdEncoded {
			return m, nil
		}
	}
	return types.Contact{}, lib.ErrNotFound
}
