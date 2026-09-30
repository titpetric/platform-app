//go:build integration

package storage_test

import (
	"testing"
	"time"

	_ "github.com/titpetric/platform-app/internal/drivers"

	"github.com/titpetric/platform-app/internal/assert"
	"github.com/titpetric/platform-app/user/schema"
	"github.com/titpetric/platform-app/user/storage"
)

func TestRevokedTokenStorage_integration(t *testing.T) {
	ctx := t.Context()

	db := NewTestDB(t)
	assert.NoError(t, storage.Migrate(ctx, db, schema.Migrations()))

	s := storage.NewRevokedTokenStorage(db)
	assert.NotNil(t, s)

	t.Run("unknown jti is not revoked", func(t *testing.T) {
		revoked, err := s.IsRevoked(ctx, "01HX0000000000000000000000")
		assert.NoError(t, err)
		assert.False(t, revoked)
	})

	t.Run("empty jti is treated as not revoked", func(t *testing.T) {
		revoked, err := s.IsRevoked(ctx, "")
		assert.NoError(t, err)
		assert.False(t, revoked)
	})

	t.Run("revoke then IsRevoked returns true", func(t *testing.T) {
		jti := "01HX0000000000000000000001"
		assert.NoError(t, s.Revoke(ctx, jti, "user-1", time.Now().Add(time.Hour)))

		revoked, err := s.IsRevoked(ctx, jti)
		assert.NoError(t, err)
		assert.True(t, revoked)
	})

	t.Run("Revoke is idempotent", func(t *testing.T) {
		jti := "01HX0000000000000000000002"
		exp := time.Now().Add(time.Hour)
		assert.NoError(t, s.Revoke(ctx, jti, "user-2", exp))
		assert.NoError(t, s.Revoke(ctx, jti, "user-2", exp))
	})

	t.Run("PurgeExpired removes only past entries", func(t *testing.T) {
		past := "01HX0000000000000000000003"
		future := "01HX0000000000000000000004"
		assert.NoError(t, s.Revoke(ctx, past, "user-3", time.Now().Add(-time.Hour)))
		assert.NoError(t, s.Revoke(ctx, future, "user-3", time.Now().Add(time.Hour)))

		n, err := s.PurgeExpired(ctx)
		assert.NoError(t, err)
		assert.True(t, n >= 1)

		// The future entry must survive.
		stillRevoked, err := s.IsRevoked(ctx, future)
		assert.NoError(t, err)
		assert.True(t, stillRevoked)

		// The past entry must be gone.
		stillRevoked, err = s.IsRevoked(ctx, past)
		assert.NoError(t, err)
		assert.False(t, stillRevoked)
	})
}
