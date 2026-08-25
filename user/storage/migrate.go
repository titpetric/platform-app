package storage

import (
	"context"
	"io/fs"
	"log/slog"

	"github.com/go-bridget/mig/migrate"
	"github.com/jmoiron/sqlx"
)

// Migrate applies SQL migrations from the given filesystem to the database.
func Migrate(ctx context.Context, db *sqlx.DB, schema fs.FS) error {
	m, err := migrate.NewManager(db, schema, "user")
	if err != nil {
		return err
	}

	applied, err := m.Apply(ctx)
	for _, item := range applied {
		slog.Default().Info("migration", "file", item.Filename, "status", item.Status)
	}
	return err
}
