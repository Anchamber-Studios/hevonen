package db

import (
	"embed"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/lib/db"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

//go:embed migrations/*.sql
var MigrationFiles embed.FS

func SetupDB(conf config.Config, logger echo.Logger) *pgx.Conn {
	return db.SetupDB("auth", conf, logger, MigrationFiles)
}
