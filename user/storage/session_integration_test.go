//go:build integration

package storage_test

import (
	"database/sql"
	"testing"

	_ "github.com/titpetric/platform-app/internal/drivers"

	"github.com/titpetric/platform-app/internal/assert"
	"github.com/titpetric/platform-app/user"
	"github.com/titpetric/platform-app/user/model"
	"github.com/titpetric/platform-app/user/schema"
	"github.com/titpetric/platform-app/user/storage"
)

func TestNewSessionStorage_integration(t *testing.T) {
	ctx := t.Context()
	ctx = user.SetSessionUser(ctx, &model.User{
		ID: "test",
	})

	db := NewTestDB(t)
	assert.NoError(t, storage.Migrate(ctx, db, schema.Migrations()))

	s := storage.NewSessionStorage(db)
	assert.NotNil(t, s)

	{
		err := s.Delete(ctx, "non-existant")
		assert.NoError(t, err)
	}

	{
		user, err := s.Get(ctx, "non-existant")
		assert.Nil(t, user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
