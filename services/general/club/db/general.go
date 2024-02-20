package db

import (
	"embed"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/anchamber-studios/hevonen/lib/db"
)

//go:embed migrations/*.sql
var MigrationFiles embed.FS

func SetupDB(conf config.Config) error {
	return db.SetupDB("clubs", conf, MigrationFiles)
}
