package storage

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"

	"github.com/go-bridget/mig/migrate"
)

func Migrate(ctx context.Context, schema embed.FS) error {
	db, err := DB(ctx)
	if err != nil {
		return err
	}

	// The embed holds the schema directory, and mig records the filename
	// relative to the root it is given, so hand it the directory.
	root, err := fs.Sub(schema, "schema")
	if err != nil {
		return err
	}

	m, err := migrate.NewManager(db, root, "maillist")
	if err != nil {
		return err
	}

	applied, err := m.Apply(ctx)
	for _, item := range applied {
		slog.Default().Info("migration", "file", item.Filename, "status", item.Status)
	}
	return err
}
