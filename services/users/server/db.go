package server

import (
	"context"
	"embed"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/lib/db"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func setupDb(conf config.Config, logger echo.Logger) *pgx.Conn {
	logger.Infof("Setup database\n")
	fmt.Printf("Setup database\n")
	conn, err := openConnection(conf, logger)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	migrationCtx := context.Background()
	err = db.Migrate(migrationCtx, logger, "users", conn, migrationFiles)
	if err != nil {
		logger.Fatalf("Migration failed: %v\n", err)
	}
	return conn
}

func openConnection(conf config.Config, logger echo.Logger) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DB.User, conf.DB.Password, conf.DB.Url, conf.DB.Port, conf.DB.Database)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		logger.Errorf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		logger.Errorf("Unable to ping database: %v\n", err)
		return nil, err
	}
	return conn, nil
}
