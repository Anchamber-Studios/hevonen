package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/labstack/echo/v4"
)

func Migrate(ctx context.Context, logger echo.Logger, service string, conn *pgx.Conn, migrationFiles embed.FS) error {
	fmt.Printf("Check database for migrations\n")
	m, err := migrate.NewMigrator(ctx, conn, fmt.Sprintf("public.%s_db_version", service))
	if err != nil {
		return err
	}

	migrationRoot, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	m.LoadMigrations(migrationRoot)

	// get the current migration status
	cur, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return err
	}
	var last int32
	for _, thisMigration := range m.Migrations {
		last = thisMigration.Sequence
	}

	if cur < last {
		logger.Infof("migration needed:  %d -> %d\n", cur, last)
		fmt.Printf("migration needed:  %d -> %d\n", cur, last)
		println(cur)
		err := m.Migrate(ctx)
		if err != nil {
			return err
		}
	} else {
		logger.Infof("no migration needed:  %d -> %d\n", cur, last)
		fmt.Printf("no migration needed:  %d -> %d\n", cur, last)
	}
	return nil
}
