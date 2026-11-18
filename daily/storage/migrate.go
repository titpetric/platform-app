package storage

import (
	"context"
	"io/fs"
	"log/slog"

	"github.com/go-bridget/mig/migrate"

	"database/sql"
)

func Migrate(ctx context.Context, db *sql.DB, schema fs.FS) error {
	m, err := migrate.NewManager(db, schema, "daily")
	if err != nil {
		return err
	}

	applied, err := m.Apply(ctx)
	for _, item := range applied {
		slog.Default().Info("migration", "file", item.Filename, "status", item.Status)
	}
	return err
}
