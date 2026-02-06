//go:build integration
// +build integration

package storage_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/pkg/assert"
)

func NewTestDB(t *testing.T) *sqlx.DB {
	dbh, err := sqlx.Connect("sqlite", ":memory:")
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, dbh.Close())
	})

	return dbh
}
