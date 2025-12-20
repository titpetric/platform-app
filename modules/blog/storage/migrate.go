package storage

import (
	"context"
	"embed"
	"io/fs"
	"path"

	"github.com/go-bridget/mig/migrate"
	"github.com/jmoiron/sqlx"
)

func Migrate(ctx context.Context, db *sqlx.DB, schema embed.FS) error {
	entries, err := fs.Glob(schema, "*.sql")
	if err != nil {
		return err
	}

	migrations := make(map[string][]byte, len(entries))
	for _, name := range entries {
		contents, _ := schema.ReadFile(name)
		migrations[path.Base(name)] = contents
	}

	return migrate.RunWithFS(
		db,
		migrations,
		&migrate.Options{
			Project: "blog",
			Apply:   true,
		},
	)
}
