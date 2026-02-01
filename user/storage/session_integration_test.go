//go:build integration
// +build integration

package storage_test

import (
	"database/sql"
	"testing"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform/pkg/require"

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

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	require.NoError(t, storage.Migrate(ctx, db, schema.Migrations))

	s := storage.NewSessionStorage(db)
	require.NotNil(t, s)

	{
		err := s.Delete(ctx, "non-existant")
		require.NoError(t, err)
	}

	{
		user, err := s.Get(ctx, "non-existant")
		require.Nil(t, user)
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}
