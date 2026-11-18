//go:build integration

package storage_test

import (
	"testing"

	"github.com/titpetric/platform-app/internal/assert"

	"database/sql"
)

func NewTestDB(t *testing.T) *sql.DB {
	dbh, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, dbh.Close())
	})

	return dbh
}
