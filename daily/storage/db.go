package storage

import (
	"context"

	"github.com/titpetric/platform"

	"database/sql"
)

// DB will return the *sql.DB in use for the user module.
// This enables reuse from outside without exposing implementation detail.
func DB(ctx context.Context) (*sql.DB, error) {
	return platform.Database.Connect(ctx, "daily")
}
