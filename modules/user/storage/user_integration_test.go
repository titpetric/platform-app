//go:build integration
// +build integration

package storage_test

import (
	"database/sql"
	"testing"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform/pkg/require"

	"github.com/titpetric/platform-app/modules/user"
	"github.com/titpetric/platform-app/modules/user/model"
	"github.com/titpetric/platform-app/modules/user/storage"
)

func TestNewUserStorage_integration(t *testing.T) {
	ctx := t.Context()
	ctx = user.SetSessionUser(ctx, &model.User{
		ID: "test",
	})

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	s := storage.NewUserStorage(db)
	require.NotNil(t, s)

	{
		user, err := s.Authenticate(ctx, model.UserAuth{})
		require.Nil(t, user)
		require.ErrorContains(t, err, "missing authentication info:")
	}

	{
		user, err := s.Authenticate(ctx, model.UserAuth{
			Email:    "me@titpetric.com",
			Password: "correct horse battery staple",
		})
		require.Nil(t, user)
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}
