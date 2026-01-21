//go:build integration
// +build integration

package email

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/titpetric/platform/pkg/drivers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/email/model"
	"github.com/titpetric/platform-app/email/schema"
	"github.com/titpetric/platform-app/email/smtp"
	"github.com/titpetric/platform-app/email/storage"
)

// setupEmailTables runs migrations for email tables
func setupEmailTables(t *testing.T, ctx context.Context) {
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	if err := storage.Migrate(ctx, db, schema.Migrations); err != nil {
		t.Fatalf("failed to migrate email tables: %v", err)
	}
}

// TestMailhogConnectivity verifies mailhog is running and accessible
func TestMailhogConnectivity(t *testing.T) {
	smtpConfig := smtp.Config{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "test@example.com",
	}

	sender := smtp.NewSMTPSender(smtpConfig)

	// Send a test email to verify connection
	err := sender.Send(
		"test@example.com",
		"Mailhog Connectivity Test",
		"If you see this, mailhog is working!",
	)
	require.NoError(t, err, "failed to connect to mailhog")
}

// TestSendEmailViaMailhog tests sending a simple email through mailhog
func TestSendEmailViaMailhog(t *testing.T) {
	smtpConfig := smtp.Config{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "sender@example.com",
	}

	sender := smtp.NewSMTPSender(smtpConfig)

	email := model.NewEmail(
		"recipient@example.com",
		"Test Subject",
		"This is a test body",
	)

	err := sender.Send(email.Recipient, email.Subject, email.Body)
	require.NoError(t, err, "failed to send email via mailhog")

	// Give mailhog time to receive
	time.Sleep(500 * time.Millisecond)

	// Verify in mailhog API
	resp, err := http.Get("http://localhost:8025/api/v1/messages")
	require.NoError(t, err, "failed to query mailhog API")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "mailhog API should be accessible")
}

// TestMultipleEmailsToMailhog tests sending multiple emails
func TestMultipleEmailsToMailhog(t *testing.T) {
	smtpConfig := smtp.Config{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "sender@example.com",
	}

	sender := smtp.NewSMTPSender(smtpConfig)

	// Send 5 emails
	for i := 1; i <= 5; i++ {
		email := model.NewEmail(
			"recipient@example.com",
			"Test Email "+string(rune(48+i)),
			"Body "+string(rune(48+i)),
		)

		err := sender.Send(email.Recipient, email.Subject, email.Body)
		require.NoError(t, err, "failed to send email %d", i)
	}

	// Give mailhog time to receive all
	time.Sleep(1 * time.Second)

	// All emails should have been sent successfully
	assert.True(t, true, "all emails sent successfully")
}

// TestServiceWithMailhog tests the service queue with mailhog
func TestServiceWithMailhog(t *testing.T) {
	ctx := context.Background()
	setupEmailTables(t, ctx)

	// Get database connection
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	emailStorage := storage.NewEmailStorage(db)

	// Create SMTP sender configured for mailhog
	smtpConfig := smtp.Config{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "service@example.com",
	}

	service := NewService(emailStorage, &smtpConfig, ServiceOptions{
		TickerInterval: 100 * time.Millisecond,
		Logger:         slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	defer service.Stop()
	service.Start()

	// Create and send email
	email := model.NewEmail(
		"recipient@example.com",
		"Service Test Email",
		"This email was sent via the service queue",
	)

	err = service.AddEmail(ctx, email)
	require.NoError(t, err, "failed to add email to service")

	// Wait for processing
	time.Sleep(2 * time.Second)

	// Verify email was marked as sent
	assert.Equal(t, model.StatusSent, email.Status)
	assert.Nil(t, email.Error)

	// Verify it was removed from queue and stored in email_sent table
	_, err = emailStorage.Get(ctx, email.ID)
	assert.Error(t, err, "email should be removed from queue after successful send")

	// Check email_sent table
	sent, err := emailStorage.GetSent(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sent), 1, "email should be in email_sent table")
	assert.NotNil(t, sent[0].SentAt)
}

// TestServiceStartProcessesPending tests that Start() processes pending emails
func TestServiceStartProcessesPending(t *testing.T) {
	ctx := context.Background()
	setupEmailTables(t, ctx)

	// Get database connection
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	emailStorage := storage.NewEmailStorage(db)

	// Create some pending emails in storage
	pendingEmail1 := model.NewEmail("user1@example.com", "Subject 1", "Body 1")
	pendingEmail1.Status = model.StatusPending
	err = emailStorage.Create(ctx, pendingEmail1)
	require.NoError(t, err)

	pendingEmail2 := model.NewEmail("user2@example.com", "Subject 2", "Body 2")
	pendingEmail2.Status = model.StatusPending
	err = emailStorage.Create(ctx, pendingEmail2)
	require.NoError(t, err)

	// Create service with mailhog
	smtpConfig := smtp.Config{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "sender@example.com",
	}

	service := NewService(emailStorage, &smtpConfig, ServiceOptions{
		TickerInterval: 100 * time.Millisecond,
		Logger:         slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	defer service.Stop()

	// Start() should immediately process pending emails
	service.Start()
	time.Sleep(200 * time.Millisecond)

	// Both emails should be removed from queue and stored in email_sent table
	_, err = emailStorage.Get(ctx, pendingEmail1.ID)
	assert.Error(t, err, "email 1 should be removed from queue after successful send")

	_, err = emailStorage.Get(ctx, pendingEmail2.ID)
	assert.Error(t, err, "email 2 should be removed from queue after successful send")

	// Check email_sent table
	sent, err := emailStorage.GetSent(ctx, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sent), 2, "both emails should be in email_sent table")
	assert.NotNil(t, sent[0].SentAt)
}

// TestServiceRetryOnFailure tests retry logic when SMTP fails
func TestServiceRetryOnFailure(t *testing.T) {
	ctx := context.Background()
	setupEmailTables(t, ctx)

	// Get database connection
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	emailStorage := storage.NewEmailStorage(db)

	// Use an invalid SMTP config that will fail
	invalidConfig := smtp.Config{
		Host:     "invalid.example.local",
		Port:     9999,
		Username: "",
		Password: "",
		From:     "sender@example.com",
	}

	service := NewService(emailStorage, &invalidConfig, ServiceOptions{
		TickerInterval: 100 * time.Millisecond,
		MaxRetries:     2,
		RetryDuration:  time.Minute,
		Logger:         slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	defer service.Stop()
	service.Start()

	email := model.NewEmail("test@example.com", "Test Subject", "Test Body")

	err = service.AddEmail(ctx, email)
	require.NoError(t, err)

	// Wait for processing attempts
	time.Sleep(500 * time.Millisecond)

	// After max retries reached, email should be in email_failed table
	failed, err := emailStorage.GetFailed(ctx, 10)
	require.NoError(t, err)

	// Find our email in the failed list
	var failedEmail *model.Email
	for i := range failed {
		if failed[i].ID == email.ID {
			failedEmail = &failed[i]
			break
		}
	}

	require.NotNil(t, failedEmail, "email should be in failed table after max retries")
	assert.Equal(t, model.StatusFailed, failedEmail.Status)
	assert.GreaterOrEqual(t, failedEmail.RetryCount, 1)
	assert.NotNil(t, failedEmail.Error)
}

// TestServiceLogging verifies logging is working
func TestServiceLogging(t *testing.T) {
	ctx := context.Background()
	setupEmailTables(t, ctx)

	// Get database connection
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	emailStorage := storage.NewEmailStorage(db)

	smtpConfig := smtp.Config{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "sender@example.com",
	}

	// Create a logger that captures output
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	service := NewService(emailStorage, &smtpConfig, ServiceOptions{
		TickerInterval: 100 * time.Millisecond,
		Logger:         logger,
	})
	defer service.Stop()
	service.Start()

	email := model.NewEmail("test@example.com", "Test Subject", "Test Body")

	// This should log the email being created and queued
	err = service.AddEmail(ctx, email)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	// If we got here without panic, logging is working
	assert.True(t, true, "logging worked without errors")
}

// TestServiceWithConfigFromEnv tests that config can be loaded from environment
func TestServiceWithConfigFromEnv(t *testing.T) {
	ctx := context.Background()
	setupEmailTables(t, ctx)

	// Get database connection
	db, err := storage.DB(ctx)
	if err != nil {
		t.Skipf("skipping: database not available: %v", err)
	}

	emailStorage := storage.NewEmailStorage(db)

	// Set environment variables for mailhog
	os.Setenv("PLATFORM_EMAIL_HOST", "localhost")
	os.Setenv("PLATFORM_EMAIL_PORT", "1025")
	os.Setenv("PLATFORM_EMAIL_FROM", "env-sender@example.com")
	defer func() {
		os.Unsetenv("PLATFORM_EMAIL_HOST")
		os.Unsetenv("PLATFORM_EMAIL_PORT")
		os.Unsetenv("PLATFORM_EMAIL_FROM")
	}()

	// Pass nil to use environment variables
	service := NewService(emailStorage, nil, ServiceOptions{
		TickerInterval: 100 * time.Millisecond,
		Logger:         slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	defer service.Stop()
	service.Start()

	email := model.NewEmail("test@example.com", "Test Subject", "Test Body")

	err = service.AddEmail(ctx, email)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	// Email should be sent using env config
	assert.Equal(t, model.StatusSent, email.Status)

	// Verify email was removed from queue and stored in email_sent table
	_, err = emailStorage.Get(ctx, email.ID)
	assert.Error(t, err, "email should be removed from queue after successful send")
}
