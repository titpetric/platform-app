//go:build integration

package storage_test

import (
	"testing"
	"time"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/titpetric/platform/pkg/require"

	"github.com/titpetric/platform-app/user/schema"
	"github.com/titpetric/platform-app/user/storage"
)

func TestRevokedTokenStorage_integration(t *testing.T) {
	ctx := t.Context()

	db := NewTestDB(t)
	require.NoError(t, storage.Migrate(ctx, db, schema.Migrations()))

	s := storage.NewRevokedTokenStorage(db)
	require.NotNil(t, s)

	t.Run("unknown jti is not revoked", func(t *testing.T) {
		revoked, err := s.IsRevoked(ctx, "01HX0000000000000000000000")
		require.NoError(t, err)
		require.False(t, revoked)
	})

	t.Run("empty jti is treated as not revoked", func(t *testing.T) {
		revoked, err := s.IsRevoked(ctx, "")
		require.NoError(t, err)
		require.False(t, revoked)
	})

	t.Run("revoke then IsRevoked returns true", func(t *testing.T) {
		jti := "01HX0000000000000000000001"
		require.NoError(t, s.Revoke(ctx, jti, "user-1", time.Now().Add(time.Hour)))

		revoked, err := s.IsRevoked(ctx, jti)
		require.NoError(t, err)
		require.True(t, revoked)
	})

	t.Run("Revoke is idempotent", func(t *testing.T) {
		jti := "01HX0000000000000000000002"
		exp := time.Now().Add(time.Hour)
		require.NoError(t, s.Revoke(ctx, jti, "user-2", exp))
		require.NoError(t, s.Revoke(ctx, jti, "user-2", exp))
	})

	t.Run("PurgeExpired removes only past entries", func(t *testing.T) {
		past := "01HX0000000000000000000003"
		future := "01HX0000000000000000000004"
		require.NoError(t, s.Revoke(ctx, past, "user-3", time.Now().Add(-time.Hour)))
		require.NoError(t, s.Revoke(ctx, future, "user-3", time.Now().Add(time.Hour)))

		n, err := s.PurgeExpired(ctx)
		require.NoError(t, err)
		require.True(t, n >= 1)

		// The future entry must survive.
		stillRevoked, err := s.IsRevoked(ctx, future)
		require.NoError(t, err)
		require.True(t, stillRevoked)

		// The past entry must be gone.
		stillRevoked, err = s.IsRevoked(ctx, past)
		require.NoError(t, err)
		require.False(t, stillRevoked)
	})
}
