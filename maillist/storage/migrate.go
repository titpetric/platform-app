package storage

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"path"

	"github.com/go-bridget/mig/migrate"
)

func Migrate(ctx context.Context, schema embed.FS) error {
	db, err := DB(ctx)
	if err != nil {
		return err
	}

	entries, err := fs.Glob(schema, "schema/*.sql")
	if err != nil {
		return err
	}

	migrations := make(map[string][]byte, len(entries))
	for _, name := range entries {
		contents, _ := schema.ReadFile(name)
		migrations[path.Base(name)] = contents
	}

	options := migrate.NewOptions(slog.Default())
	options.Project = "maillist"
	options.Apply = true

	return migrate.RunWithFS(ctx, db, migrations, options)
}
