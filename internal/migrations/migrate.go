package migrations

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
)

func Migrate(dsn string) error {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get caller info")
	}
	migrationsDir := filepath.Join(filepath.Dir(file), ".")

	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		return fmt.Errorf("goose failed to open DB: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}

	return nil
}
