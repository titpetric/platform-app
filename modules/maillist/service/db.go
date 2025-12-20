package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform"
)

// DB will return the sqlx.DB in use for the maillist module.
func DB(ctx context.Context) (*sqlx.DB, error) {
	return platform.Database.Connect(ctx, "maillist")
}
