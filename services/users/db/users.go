package db

import (
	"context"

	"github.com/anchamber-studios/hevonen/lib"
	"github.com/anchamber-studios/hevonen/services/users/client"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/sqids/sqids-go"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	List(ctx context.Context) ([]client.User, error)
	Get(ctx context.Context, id string) (client.User, error)
	Create(ctx context.Context, user client.UserCreate) (string, error)
	Login(ctx context.Context, user client.UserLogin) (client.User, error)
}

type UserRepoPostgre struct {
	Logger       echo.Logger
	DB           *pgx.Conn
	IdConversion *sqids.Sqids
}

const (
	IdOffset uint64 = 2345678901
)

func (r *UserRepoPostgre) List(ctx context.Context) ([]client.User, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, email, email_confirmed, active, updated_at, created_at FROM users.users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []client.User
	for rows.Next() {
		var id uint64
		var user client.User
		err = rows.Scan(&id, &user.Email, &user.EmailConfirmed, &user.Active, &user.UpdatedAt, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
		if err != nil {
			r.Logger.Errorf("id conversion for '%d' failed: %v\n", id, err)
			return nil, err
		}
		user.ID = cId
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepoPostgre) Get(ctx context.Context, id string) (client.User, error) {
	var user client.User
	cId := r.IdConversion.Decode(id)
	var dbId uint64
	err := r.DB.QueryRow(ctx, "SELECT id, email, email_confirmed, active, updated_at, created_at FROM users.users WHERE id = $1;", cId[0]).Scan(&dbId, &user.Email, &user.EmailConfirmed, &user.Active, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		return client.User{}, err
	}
	user.ID = id
	return user, nil
}

func (r *UserRepoPostgre) Create(ctx context.Context, user client.UserCreate) (string, error) {
	var id uint64
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	err = r.DB.QueryRow(ctx, "INSERT INTO users.users (email, password) VALUES ($1, $2) RETURNING id;", user.Email, password).Scan(&id)
	if err != nil {
		r.Logger.Errorf("insertion of user failed: %v\n", id, err)
		return "", err
	}
	cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
	if err != nil {
		r.Logger.Errorf("id conversion for '%d' failed: %v\n", id, err)
		return "", err
	}
	r.Logger.Info("user '%s' created", cId)
	return cId, nil
}

func (r *UserRepoPostgre) Login(ctx context.Context, login client.UserLogin) (client.User, error) {
	var user client.User
	var id uint64
	var hashedPassword string
	err := r.DB.QueryRow(
		ctx,
		"SELECT id, email, password FROM users.users WHERE email = $1;",
		login.Email).Scan(&id, &user.Email, &hashedPassword)
	if err != nil {
		return client.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(login.Password))
	if err != nil {
		r.Logger.Errorf("hash comparison of password failed: %v\n", err)
		return client.User{}, lib.ErrLoginFailed
	}

	cId, err := r.IdConversion.Encode([]uint64{id, IdOffset})
	if err != nil {
		r.Logger.Errorf("id conversion for '%d' failed: %v\n", id, err)
		return client.User{}, err
	}
	user.ID = cId
	r.Logger.Info("user '%s' logged in", cId)
	return user, nil
}
