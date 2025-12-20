package service

import (
	"context"
	"embed"
	"io/fs"
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

	return migrate.RunWithFS(
		db,
		migrations,
		&migrate.Options{
			Project: "maillist",
			Apply:   true,
		},
	)
}
