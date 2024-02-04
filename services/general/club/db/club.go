package db

import (
	"context"
	"strings"

	"github.com/anchamber-studios/hevonen/lib/logger"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/jackc/pgx/v5"
	"github.com/sqids/sqids-go"
)

type ClubRepo interface {
	List(ctx context.Context) ([]types.Club, error)
	ListForIdentity(ctx context.Context, identity string) ([]types.ClubMember, error)
	CreateWithAdminMember(ctx context.Context, club types.ClubCreate, admin types.MemberCreate) (string, error)
	Create(ctx context.Context, club types.ClubCreate) (string, error)
	Get(ctx context.Context, clubIdEncoded string) (types.Club, error)
}

type ClubRepoPostgre struct {
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
}

const (
	idOffsetClub uint64 = 1234567890
)

func (r *ClubRepoPostgre) List(ctx context.Context) ([]types.Club, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name, website, email, phone FROM clubs.clubs;")
	if err != nil {
		return nil, err
	}
	var clubs []types.Club
	for rows.Next() {
		var id uint64
		var club types.Club
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

func (r *ClubRepoPostgre) ListForIdentity(ctx context.Context, identity string) ([]types.ClubMember, error) {
	rows, err := r.DB.Query(ctx, `
	SELECT c.id, c.name, string_agg(mr.role_name, ',')
	FROM clubs.clubs c 
	INNER JOIN clubs.members m 
		ON c.id = m.club_id 
	INNER JOIN clubs.member_roles mr
		ON m.id = mr.member_id
	WHERE m.identity_id = $1
	GROUP BY c.id, c.name;
	`, identity)
	if err != nil {
		return nil, err
	}
	var clubs []types.ClubMember
	for rows.Next() {
		var club types.ClubMember
		var roles string
		err := rows.Scan(&club.ID, &club.Name, &roles)
		if err != nil {
			return nil, err
		}
		club.Roles = strings.Split(roles, ",")
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (r *ClubRepoPostgre) Create(ctx context.Context, club types.ClubCreate) (string, error) {
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

func (r *ClubRepoPostgre) CreateWithAdminMember(ctx context.Context, club types.ClubCreate, admin types.MemberCreate) (string, error) {
	log := logger.FromContext(ctx)
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	var clubID string
	err = tx.QueryRow(ctx, `
		INSERT INTO clubs.clubs (name, website, email, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
		`, club.Name, club.Website, club.Email, club.Phone).Scan(&clubID)
	if err != nil {
		log.Sugar().Errorf("Failed to create club: %v", err)
		tx.Rollback(ctx)
		return "", err
	}
	admin.ClubID = clubID
	var memberID string
	err = tx.QueryRow(ctx, `
	INSERT INTO clubs.members (identity_id, club_id, email) 
	VALUES ($1, $2, $3)
	RETURNING id;
	`, admin.IdentityID, clubID, admin.Email).Scan(&memberID)
	if err != nil {
		log.Sugar().Errorf("Failed to create admin member: %v", err)
		tx.Rollback(ctx)
		return "", err
	}
	_, err = tx.Exec(ctx, `
	INSERT INTO clubs.member_roles (member_id, role_name)
	VALUES ($1, 'admin');
	`, memberID)
	if err != nil {
		log.Sugar().Errorf("Failed to assign admin role to member: %v", err)
		tx.Rollback(ctx)
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Sugar().Errorf("Failed to commit transaction: %v", err)
		tx.Rollback(ctx)
		return "", err
	}
	return clubID, nil
}

func (r *ClubRepoPostgre) Get(ctx context.Context, clubIdEncoded string) (types.Club, error) {
	cId := r.IdConversion.Decode(clubIdEncoded)
	var club types.Club
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
