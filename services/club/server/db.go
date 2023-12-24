package server

import (
	"context"
	"embed"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/lib/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func setupDb(conf config.Config, logger echo.Logger) *pgxpool.Pool {
	logger.Infof("Setup database\n")
	fmt.Printf("Setup database\n")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DB.User, conf.DB.Password, conf.DB.Url, conf.DB.Port, conf.DB.Database)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatalf("Unable to ping database: %v\n", err)
	}

	migrationCtx := context.Background()
	con, err := pgx.Connect(migrationCtx, dsn)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	err = db.Migrate(migrationCtx, logger, con, migrationFiles)
	if err != nil {
		logger.Fatalf("Migration failed: %v\n", err)
	}
	err = con.Close(migrationCtx)
	if err != nil {
		logger.Fatalf("Migration failed: %v\n", err)
	}
	return pool
}
