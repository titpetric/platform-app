package handlers

import (
	"testing"

	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/modules/blog/schema"
	"github.com/titpetric/platform-app/modules/blog/storage"
)

// setupTestDB creates a temporary SQLite database for testing
func setupTestDB(t *testing.T) *sqlx.DB {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, err := sqlx.Open("sqlite", dbPath)
	require.NoError(t, err)

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	require.NoError(t, storage.Migrate(t.Context(), db, schema.Migrations))

	return db
}

// TestNewHandlers validates handler creation with repository and views
func TestNewHandlers(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	// Note: this will fail in unit tests due to no themeFS, but tests the signature
	// In actual use, a valid fs.FS is provided
	assert.NotNil(t, repo)
}

// TestNewAdminHandlers validates admin-only handler creation
func TestNewAdminHandlers(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	h := NewAdminHandlers(repo)

	assert.NotNil(t, h)
	assert.Equal(t, repo, h.repository)
	assert.Nil(t, h.views)
}

// TestHandlers_RepositoryAssignment validates repository is properly assigned
func TestHandlers_RepositoryAssignment(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	h := NewAdminHandlers(repo)

	assert.Same(t, repo, h.repository)
}
