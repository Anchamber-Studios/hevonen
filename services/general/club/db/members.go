package db

import (
	"context"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/jackc/pgx/v5"
	"github.com/sqids/sqids-go"
)

type MemberRepo interface {
	List(ctx context.Context) ([]types.Member, error)
	ListForClub(ctx context.Context, clubIdEncoded string) ([]types.Member, error)
	Create(ctx context.Context, member types.MemberCreate) (string, error)
	Get(ctx context.Context, memberIdEncoded string) (types.Member, error)
}

type MemberRepoPostgre struct {
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
}

const (
	IdOffset uint64 = 1234567890
)

func (r *MemberRepoPostgre) List(ctx context.Context) ([]types.Member, error) {
	rows, err := r.DB.Query(ctx, `
	SELECT id, first_name, middle_name, last_name, email, phone, height, weight FROM clubs.contacts;
	`)
	if err != nil {
		return nil, err
	}

	members := make([]types.Member, 0)
	for rows.Next() {
		member := types.Member{}
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

func (r *MemberRepoPostgre) ListForClub(ctx context.Context, clubIdEncoded string) ([]types.Member, error) {
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

	members := make([]types.Member, 0)
	for rows.Next() {
		member := types.Member{}
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

func (r *MemberRepoPostgre) Create(ctx context.Context, member types.MemberCreate) (string, error) {
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

func (r *MemberRepoPostgre) Get(ctx context.Context, memberIdEncoded string) (types.Member, error) {
	memberId := r.IdConversion.Decode(memberIdEncoded)[0]
	var id uint64
	var member types.Member
	r.DB.QueryRow(ctx, "SELECT id, first_name, last_name, email, phone FROM clubs.members WHERE id = $1;", memberId).Scan(&id, &member.FirstName, &member.LastName, &member.Email, &member.Phone)
	if id == 0 {
		return member, lib.ErrNotFound
	}
	member.ID = memberIdEncoded
	return member, nil
}

type MemberRepoMock struct {
	Members []types.Member
}

func (r *MemberRepoMock) List(ctx context.Context) ([]types.Member, error) {
	return r.Members, nil
}

func (r *MemberRepoMock) ListForClub(ctx context.Context, clubIdEncoded string) ([]types.Member, error) {
	return r.Members, nil
}

func (r *MemberRepoMock) Create(ctx context.Context, member types.MemberCreate) (string, error) {
	for _, m := range r.Members {
		if m.Email == member.Email {
			return "", lib.ErrAlreadyExists
		}
	}
	id := fmt.Sprintf("%d", len(r.Members)+1)
	r.Members = append(r.Members, types.Member{
		ID:     id,
		ClubID: member.ClubID,
		Email:  member.Email,
	})
	return id, nil
}

func (r *MemberRepoMock) Get(ctx context.Context, memberIdEncoded string) (types.Member, error) {
	for _, m := range r.Members {
		if m.ID == memberIdEncoded {
			return m, nil
		}
	}
	return types.Member{}, lib.ErrNotFound
}
