package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"golang.org/x/xerrors"
)

type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	SSLMode    string
	Migrations string
}

// NewPostgresConnection initializes the DB connection to the Postgres instance.
func NewPostgresConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, xerrors.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr != nil {
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}

	if pingErr != nil {
		return nil, xerrors.Errorf("failed to ping PostgreSQL: %w", pingErr)
	}

	if err := runMigrations(db, cfg.Migrations); err != nil {
		return nil, xerrors.Errorf("running migrations failed: %w", err)
	}

	return db, nil
}

// runMigrations ensures that migrations in the migrationsFolder are applied to the Postgres DB during connection init.
func runMigrations(db *sql.DB, migrationsFolder string) error {
	if err := goose.Up(db, migrationsFolder); err != nil {
		if goosErr := goose.Reset(db, migrationsFolder); goosErr != nil {
			return xerrors.Errorf("failed to reset migrations: %w", goosErr)
		}

		if dbErr := db.Close(); dbErr != nil {
			return xerrors.Errorf("failed to close db connection: %w", dbErr)
		}

		return xerrors.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
