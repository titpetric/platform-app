package storage

import (
	"context"

	"github.com/titpetric/platform"

	"database/sql"
)

// DB will return the *sql.DB in use for the blog module.
// This enables reuse from outside without exposing implementation detail.
// Uses the "blog" named connection from platform.Database.
func DB(ctx context.Context) (*sql.DB, error) {
	return platform.Database.Connect(ctx, "blog")
}
