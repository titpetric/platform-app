package storage

import (
	"context"

	"github.com/titpetric/platform"

	"database/sql"
)

// DB will return the *sql.DB in use for the pulse module.
// This enables reuse of the storage package from outside
// the app without exposing implementation detail.
func DB(ctx context.Context) (*sql.DB, error) {
	return platform.Database.Connect(ctx, "pulse")
}
