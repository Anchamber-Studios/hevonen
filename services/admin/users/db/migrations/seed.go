package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/admin/users/client"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	conf := config.LoadConfig()
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DB.User, conf.DB.Password, conf.DB.Url, conf.DB.Port, conf.DB.Database)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db")
	}

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("failed to ping db")
	}

	seed(ctx, conn)
}

func seed(ctx context.Context, conn *pgx.Conn) {
	users := []client.UserCreate{
		{
			Email:    "admin@hevonen.io",
			Password: "123!QWEqwe",
			Apps:     []client.AppConnection{},
		},
		{
			Email:    "rene@hevonen.io",
			Password: "123!QWEqwe",
			Apps:     []client.AppConnection{},
		},
		{
			Email:    "lea@hevonen.io",
			Password: "123!QWEqwe",
			Apps:     []client.AppConnection{},
		},
		{
			Email:    "charlotte@hevonen.io",
			Password: "123!QWEqwe",
			Apps:     []client.AppConnection{},
		},
	}

	for _, user := range users {
		var id uint64
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to generate password")
		}
		err = conn.QueryRow(ctx, "INSERT INTO users.users (email, password, email_confirmed) VALUES ($1, $2, $3) RETURNING id;", user.Email, password, true).Scan(&id)
		if err != nil {
			log.Fatalf("failed to insert user: %v", err)
		}
		log.Printf("user '%s' created", user.Email)
	}
}
