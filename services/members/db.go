package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/labstack/echo/v4"
)

func setupDb(config Config, logger echo.Logger) *pgxpool.Pool {
	logger.Infof("Setup database\n")
	fmt.Printf("Setup database\n")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DB.User, config.DB.Password, config.DB.Url, config.DB.Port, config.DB.Database)

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
	err = runMigrations(migrationCtx, con, logger)
	if err != nil {
		logger.Fatalf("Migration failed: %v\n", err)
	}
	err = con.Close(migrationCtx)
	if err != nil {
		logger.Fatalf("Migration failed: %v\n", err)
	}
	return pool
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func runMigrations(ctx context.Context, conn *pgx.Conn, logger echo.Logger) error {
	fmt.Printf("Check database for migrations\n")
	m, err := migrate.NewMigrator(ctx, conn, "members.db_version")
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
		m.Migrate(ctx)
	} else {
		logger.Infof("no migration needed:  %d -> %d\n", cur, last)
		fmt.Printf("no migration needed:  %d -> %d\n", cur, last)
	}
	return nil
}
