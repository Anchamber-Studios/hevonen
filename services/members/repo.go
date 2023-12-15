package main

import (
	"context"
	"fmt"

	"github.com/anchamber-studios/hevonen/services/members/client"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqids/sqids-go"
)

type MemberRepo interface {
	List(ctx context.Context) ([]client.Member, error)
	Create(ctx context.Context, member client.MemberCreate) (string, error)
	Get(ctx context.Context, memberIdEncoded string) (client.Member, error)
}

type MemberRepoPostgre struct {
	DB           *pgxpool.Conn
	IdConversion *sqids.Sqids
}

const (
	IdOffset uint64 = 1234567890
)

func (r *MemberRepoPostgre) List(ctx context.Context) ([]client.Member, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, first_name, last_name, email, phone FROM members.members;")
	if err != nil {
		return nil, err
	}

	members := make([]client.Member, 0)
	for rows.Next() {
		member := client.Member{}
		var id uint64
		err := rows.Scan(&id, &member.FirstName, &member.LastName, &member.Email, &member.Phone)
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

func (r *MemberRepoPostgre) Create(ctx context.Context, member client.MemberCreate) (string, error) {
	var id uint64
	err := r.DB.QueryRow(ctx, "INSERT INTO members.members (first_name, last_name, email, phone) VALUES ($1, $2, $3, $4) RETURNING id;", member.FirstName, member.LastName, member.Email, member.Phone).Scan(&id)
	if err != nil {
		return "", err
	}
	cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
	if err != nil {
		return "", err
	}
	return cId, nil
}

func (r *MemberRepoPostgre) Get(ctx context.Context, memberIdEncoded string) (client.Member, error) {
	memberId := r.IdConversion.Decode(memberIdEncoded)[0]
	var id uint64
	var member client.Member
	r.DB.QueryRow(ctx, "SELECT id, first_name, last_name, email, phone FROM members.members WHERE id = $1;", memberId).Scan(&id, &member.FirstName, &member.LastName, &member.Email, &member.Phone)
	if id == 0 {
		return member, ErrNotFound
	}
	member.ID = memberIdEncoded
	return member, nil
}

type MemberRepoMock struct {
	Members []client.Member
}

func (r *MemberRepoMock) List(ctx context.Context) ([]client.Member, error) {
	return r.Members, nil
}

func (r *MemberRepoMock) Create(ctx context.Context, member client.MemberCreate) (string, error) {
	for _, m := range r.Members {
		if m.Email == member.Email {
			return "", ErrAlreadyExists
		}
	}
	id := fmt.Sprintf("%d", len(r.Members)+1)
	r.Members = append(r.Members, client.Member{
		ID:        id,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     member.Email,
		Phone:     member.Phone,
		Height:    170,
		Weight:    75,
	})
	return id, nil
}

func (r *MemberRepoMock) Get(ctx context.Context, memberIdEncoded string) (client.Member, error) {
	for _, m := range r.Members {
		if m.ID == memberIdEncoded {
			return m, nil
		}
	}
	return client.Member{}, ErrNotFound
}
