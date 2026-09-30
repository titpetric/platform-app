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

func TestNewUserStorage_integration(t *testing.T) {
	ctx := t.Context()
	ctx = user.SetSessionUser(ctx, &model.User{
		ID: "test",
	})

	db := NewTestDB(t)
	assert.NoError(t, storage.Migrate(ctx, db, schema.Migrations()))

	s := storage.NewUserStorage(db)
	assert.NotNil(t, s)

	{
		user, err := s.Authenticate(ctx, model.UserAuth{})
		assert.Nil(t, user)
		assert.ErrorContains(t, err, "missing authentication info:")
	}

	{
		user, err := s.Authenticate(ctx, model.UserAuth{
			Email:    "me@titpetric.com",
			Password: "horse battery staple",
		})
		assert.Nil(t, user)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}

	{
		user, err := s.Create(ctx, &model.UserCreateRequest{
			FullName: "Tit Petric",
			Email:    "me@titpetric.com",
			Password: "horse battery staple",
			Username: "titpetric",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, user)
	}

	{
		user, err := s.Create(ctx, &model.UserCreateRequest{
			FullName: "No Username",
			Email:    "nousername@titpetric.com",
			Password: "horse battery staple",
		})
		assert.Nil(t, user)
		assert.ErrorContains(t, err, "username is required")
	}

	{
		userlist, err := s.List(ctx)
		assert.NoError(t, err)
		assert.True(t, len(userlist) == 1)
	}
}
