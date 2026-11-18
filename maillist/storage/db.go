package storage

import (
	"context"

	"github.com/titpetric/platform"

	"database/sql"
)

// DB will return the *sql.DB in use for the maillist module.
func DB(ctx context.Context) (*sql.DB, error) {
	return platform.Database.Connect(ctx, "maillist")
}
