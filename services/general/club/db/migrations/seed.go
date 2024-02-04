package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/services/club/shared/types"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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
	clubs := []types.ClubCreate{
		{Name: "Speedy Stallions", Website: "https://www.speedystallions.net"},
		{Name: "Whinny Warriors", Website: "https://www.whinnywarriors.club"},
		{Name: "Buckle Up", Website: "https://www.buckleup.club"},
	}
	var ids = make([]string, len(clubs))
	for i, club := range clubs {
		var id uint64
		err := conn.QueryRow(ctx, "INSERT INTO club.clubs (name, website) VALUES ($1, $2) RETURNING id;", club.Name, club.Website).Scan(&id)
		if err != nil {
			fmt.Printf("error inserting club %d: %v\n", i, err)
		}
		ids[i] = fmt.Sprintf("%d", id)
	}
}
