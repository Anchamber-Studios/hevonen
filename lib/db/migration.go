package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/anchamber-studios/hevonen/lib/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

func Migrate(ctx context.Context, service string, conn *pgx.Conn, migrationFiles embed.FS) error {
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
		fmt.Printf("migration needed:  %d -> %d\n", cur, last)
		println(cur)
		err := m.Migrate(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResetDB(ctx context.Context, service string, conf config.Config) error {
	conn, err := OpenConnection(conf)
	m, err := migrate.NewMigrator(ctx, conn, fmt.Sprintf("public.%s_db_version", service))
	if err != nil {
		return err
	}
	m.MigrateTo(ctx, 0)
	return nil
}

func SetupDB(service string, conf config.Config, migrationFiles embed.FS) error {
	conn, err := OpenConnection(conf)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	migrationCtx := context.Background()
	err = Migrate(migrationCtx, service, conn, migrationFiles)
	if err != nil {
		return err
	}
	return nil
}

func OpenConnection(conf config.Config) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.DB.User, conf.DB.Password, conf.DB.Url, conf.DB.Port, conf.DB.Database)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}
