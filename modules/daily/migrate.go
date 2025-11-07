package daily

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"path"

	"github.com/go-bridget/mig/migrate"
	"github.com/jmoiron/sqlx"
)

//go:embed model/*.up.sql
var schema embed.FS

func Migrate(ctx context.Context, db *sqlx.DB) error {
	migrations, err := loadMigrations(schema)
	if err != nil {
		return err
	}

	return migrate.RunWithFS(
		db,
		migrations,
		&migrate.Options{
			Project: "daily",
			Apply:   true,
			Verbose: true,
		},
	)
}

func loadMigrations(filesystem fs.FS) (migrate.FS, error) {
	entries, err := fs.Glob(schema, "model/*.up.sql")
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, sql.ErrNoRows
	}

	migrations := make(map[string][]byte, len(entries))

	for _, name := range entries {
		contents, _ := schema.ReadFile(name)
		migrations[path.Base(name)] = contents
	}

	return migrations, nil
}
