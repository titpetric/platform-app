//go:build integration
// +build integration

package storage_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/titpetric/platform/pkg/assert"

	_ "github.com/titpetric/platform-app/autoload"
	"github.com/titpetric/platform-app/modules/daily"
	"github.com/titpetric/platform-app/modules/daily/model"
	"github.com/titpetric/platform-app/modules/daily/storage"
)

func Must[T any](t *testing.T, ctor func(context.Context) (T, error)) T {
	res, err := ctor(t.Context())
	assert.NoError(t, err)
	return res
}

func TestStorage(t *testing.T) {
	ctx := t.Context()
	db := Must[*sqlx.DB](t, storage.DB)

	assert.NoError(t, daily.Migrate(ctx, db))

	repo := storage.NewStorage(db)

	{
		_, err := repo.Get(ctx, "test")

		require.Error(t, err)
	}

	{
		todo := model.Todo{
			Title: "Hello world",
		}
		got, err := repo.Add(ctx, todo)

		assert.NoError(t, err)
		assert.NotEmpty(t, got.ID)

		todos, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, todos)

		got2, err := repo.Get(ctx, got.ID)

		assert.NoError(t, err)
		assert.Equal(t, got.ID, got2.ID)

		assert.NoError(t, repo.Delete(ctx, got.ID))
	}
}
