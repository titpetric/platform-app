//go:build integration
// +build integration

package storage_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/modules/email/model"
	"github.com/titpetric/platform-app/modules/email/schema"
	"github.com/titpetric/platform-app/modules/email/storage"
)

func setupTestDB(t *testing.T, ctx context.Context) {
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	if err := storage.Migrate(ctx, db, schema.Migrations); err != nil {
		t.Fatalf("failed to migrate email tables: %v", err)
	}
}

func TestEmailStorageCreate(t *testing.T) {
	ctx := context.Background()
	setupTestDB(t, ctx)

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	emailStorage := storage.NewEmailStorage(db)

	email := model.NewEmail("test@example.com", "Subject", "Body")
	err = emailStorage.Create(ctx, email)
	require.NoError(t, err)

	assert.NotEmpty(t, email.ID)
	assert.Equal(t, "test@example.com", email.Recipient)
	assert.Equal(t, model.StatusPending, email.Status)
}

func TestEmailStorageGet(t *testing.T) {
	ctx := context.Background()
	setupTestDB(t, ctx)

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	emailStorage := storage.NewEmailStorage(db)

	// Create an email
	email := model.NewEmail("test@example.com", "Subject", "Body")
	err = emailStorage.Create(ctx, email)
	require.NoError(t, err)

	// Retrieve it
	retrieved, err := emailStorage.Get(ctx, email.ID)
	require.NoError(t, err)

	assert.Equal(t, email.ID, retrieved.ID)
	assert.Equal(t, email.Recipient, retrieved.Recipient)
}

func TestEmailStorageGetPending(t *testing.T) {
	ctx := context.Background()
	setupTestDB(t, ctx)

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	emailStorage := storage.NewEmailStorage(db)

	// Create multiple emails
	email1 := model.NewEmail("test1@example.com", "Subject 1", "Body 1")
	email2 := model.NewEmail("test2@example.com", "Subject 2", "Body 2")

	err = emailStorage.Create(ctx, email1)
	require.NoError(t, err)
	err = emailStorage.Create(ctx, email2)
	require.NoError(t, err)

	// Retrieve pending emails
	pending, err := emailStorage.GetPending(ctx, 10)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(pending), 2)
}

func TestEmailStorageUpdate(t *testing.T) {
	ctx := context.Background()
	setupTestDB(t, ctx)

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	emailStorage := storage.NewEmailStorage(db)

	// Create and update an email
	email := model.NewEmail("test@example.com", "Subject", "Body")
	err = emailStorage.Create(ctx, email)
	require.NoError(t, err)

	now := time.Now()
	email.Status = model.StatusSent
	email.SentAt = &now

	err = emailStorage.Update(ctx, email)
	require.NoError(t, err)

	// Verify email was removed from queue and stored in email_sent table
	_, err = emailStorage.Get(ctx, email.ID)
	assert.Error(t, err, "email should be removed from queue after successful send")

	// Check email_sent table
	sent, err := emailStorage.GetSent(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sent), 1, "email should be in email_sent table")
}

func TestEmailStorageUpdateFailed(t *testing.T) {
	ctx := context.Background()
	setupTestDB(t, ctx)

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	emailStorage := storage.NewEmailStorage(db)

	// Create and mark as failed
	email := model.NewEmail("test@example.com", "Subject", "Body")
	err = emailStorage.Create(ctx, email)
	require.NoError(t, err)

	errMsg := "connection timeout"
	email.Status = model.StatusFailed
	email.Error = &errMsg
	email.RetryCount = 3

	err = emailStorage.Update(ctx, email)
	require.NoError(t, err)

	// Verify email was removed from queue and stored in email_failed table
	_, err = emailStorage.Get(ctx, email.ID)
	assert.Error(t, err, "email should be removed from queue after failure")

	// Check email_failed table
	failed, err := emailStorage.GetFailed(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(failed), 1, "email should be in email_failed table")
	assert.Equal(t, model.StatusFailed, failed[0].Status)
	assert.Equal(t, 3, failed[0].RetryCount)
}

func TestEmailStorageGetFailed(t *testing.T) {
	ctx := context.Background()
	setupTestDB(t, ctx)

	db, err := storage.DB(ctx)
	require.NoError(t, err)

	emailStorage := storage.NewEmailStorage(db)

	// Create and fail multiple emails
	email1 := model.NewEmail("test1@example.com", "Subject 1", "Body 1")
	email1.Status = model.StatusFailed
	errMsg := "error1"
	email1.Error = &errMsg
	err = emailStorage.Create(ctx, email1)
	require.NoError(t, err)

	email2 := model.NewEmail("test2@example.com", "Subject 2", "Body 2")
	email2.Status = model.StatusFailed
	errMsg2 := "error2"
	email2.Error = &errMsg2
	err = emailStorage.Create(ctx, email2)
	require.NoError(t, err)

	// Move to failed table
	emailStorage.Update(ctx, email1)
	emailStorage.Update(ctx, email2)

	// Retrieve failed emails
	failed, err := emailStorage.GetFailed(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(failed), 2)
}
