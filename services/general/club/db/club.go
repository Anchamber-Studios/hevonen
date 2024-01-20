package db

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/jackc/pgx/v5"
	"github.com/sqids/sqids-go"
)

type ClubRepo interface {
	List(ctx context.Context) ([]client.Club, error)
	ListForIdentity(ctx context.Context, identity string) ([]client.ClubMember, error)
	Create(ctx context.Context, club client.ClubCreate) (string, error)
	Get(ctx context.Context, clubIdEncoded string) (client.Club, error)
}

type ClubRepoPostgre struct {
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
}

const (
	idOffsetClub uint64 = 1234567890
)

func (r *ClubRepoPostgre) List(ctx context.Context) ([]client.Club, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name, website, email, phone FROM clubs.clubs;")
	if err != nil {
		return nil, err
	}
	var clubs []client.Club
	for rows.Next() {
		var id uint64
		var club client.Club
		err := rows.Scan(&id, &club.Name, &club.Website, &club.Email, &club.Phone)
		if err != nil {
			return nil, err
		}
		club.ID, err = r.IdConversion.Encode([]uint64{id, idOffsetClub})
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (r *ClubRepoPostgre) ListForIdentity(ctx context.Context, identity string) ([]client.ClubMember, error) {
	rows, err := r.DB.Query(ctx, `
	SELECT c.id, c.name 
	FROM clubs.clubs c 
	INNER JOIN clubs.members m 
		ON c.id = m.club_id 
	WHERE m.identity_id = $1;
	`, identity)
	if err != nil {
		return nil, err
	}
	var clubs []client.ClubMember
	for rows.Next() {
		var club client.ClubMember
		err := rows.Scan(&club.ID, &club.Name)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (r *ClubRepoPostgre) Create(ctx context.Context, club client.ClubCreate) (string, error) {
	var id string
	err := r.DB.QueryRow(ctx, `
		INSERT INTO clubs.clubs (name, website, email, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
		`, club.Name, club.Website, club.Email, club.Phone).
		Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *ClubRepoPostgre) Get(ctx context.Context, clubIdEncoded string) (client.Club, error) {
	cId := r.IdConversion.Decode(clubIdEncoded)
	var club client.Club
	err := r.DB.QueryRow(ctx, `
		SELECT id, name, website, email, phone
		FROM clubs.clubs
		WHERE id = $1;
		`, cId[0]).
		Scan(&club.ID, &club.Name, &club.Website, &club.Email, &club.Phone)
	if err != nil {
		return club, err
	}
	return club, nil
}
