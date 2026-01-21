package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/pkg/ulid"

	"github.com/titpetric/platform-app/email/model"
)

// EmailStorage handles database operations for emails
type EmailStorage struct {
	db *sqlx.DB
}

// NewEmailStorage creates a new email storage instance
func NewEmailStorage(db *sqlx.DB) *EmailStorage {
	return &EmailStorage{db: db}
}

// Create inserts a new email into the email table
func (s *EmailStorage) Create(ctx context.Context, email *model.Email) error {
	email.ID = ulid.String()
	email.SetCreatedAt(time.Now())

	query := email.Insert()
	_, err := s.db.NamedExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an email by ID from email
func (s *EmailStorage) Get(ctx context.Context, id string) (*model.Email, error) {
	email := &model.Email{}
	query := `SELECT * FROM email WHERE id=?`

	err := s.db.GetContext(ctx, email, query, id)
	if err != nil {
		return nil, err
	}

	return email, nil
}

// GetPending retrieves all pending emails from email
func (s *EmailStorage) GetPending(ctx context.Context, limit int) ([]model.Email, error) {
	var emails []model.Email
	query := `SELECT * FROM email WHERE status='pending' ORDER BY created_at ASC LIMIT ?`

	err := s.db.SelectContext(ctx, &emails, query, limit)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

// Update updates an email's status in email
// If email is marked as sent, it moves to email_sent table
// If email is marked as failed, it moves to email_failed table
func (s *EmailStorage) Update(ctx context.Context, email *model.Email) error {
	// If sent, insert into email_sent for audit trail
	if email.Status == model.StatusSent && email.SentAt != nil {
		insertSentQuery := email.Insert(model.WithTable(model.EmailSentTable).WithStatement("INSERT OR IGNORE"))
		_, err := s.db.NamedExecContext(ctx, insertSentQuery, email)
		if err != nil {
			// Log but don't fail - audit trail insertion is secondary
		}

		// Remove from email table after successful send
		deleteQuery := `DELETE FROM email WHERE id=:id`
		_, err = s.db.NamedExecContext(ctx, deleteQuery, email)
		return err
	}

	// If failed, insert into email_failed for audit trail
	if email.Status == model.StatusFailed {
		insertFailedQuery := email.Insert(model.WithTable(model.EmailFailedTable).WithStatement("INSERT OR IGNORE"))
		_, err := s.db.NamedExecContext(ctx, insertFailedQuery, email)
		if err != nil {
			// Log but don't fail - audit trail insertion is secondary
		}

		// Remove from email table after failure recorded
		deleteQuery := `DELETE FROM email WHERE id=:id`
		_, err = s.db.NamedExecContext(ctx, deleteQuery, email)
		return err
	}

	// For pending emails, update email table with retry info
	query := `UPDATE email SET retry_count=:retry_count, retry_error=:retry_error, retry_at=:retry_at WHERE id=:id`

	_, err := s.db.NamedExecContext(ctx, query, email)
	return err
}

// GetSent retrieves sent emails from email_sent table (for audit/records)
func (s *EmailStorage) GetSent(ctx context.Context, limit int) ([]model.Email, error) {
	var emails []model.Email
	query := `SELECT * FROM email_sent ORDER BY sent_at DESC LIMIT ?`

	err := s.db.SelectContext(ctx, &emails, query, limit)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

// GetFailed retrieves failed emails from email_failed table (for audit/records)
func (s *EmailStorage) GetFailed(ctx context.Context, limit int) ([]model.Email, error) {
	var emails []model.Email
	query := `SELECT * FROM email_failed ORDER BY created_at DESC LIMIT ?`

	err := s.db.SelectContext(ctx, &emails, query, limit)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
