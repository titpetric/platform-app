package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/pkg/telemetry"
	"github.com/titpetric/platform/pkg/ulid"

	"github.com/titpetric/platform-app/user/model"
)

// PasskeyStorage implements passkey persistence using the database.
type PasskeyStorage struct {
	db *sqlx.DB
}

// NewPasskeyStorage creates a new PasskeyStorage.
func NewPasskeyStorage(db *sqlx.DB) *PasskeyStorage {
	return &PasskeyStorage{
		db: db,
	}
}

// Create inserts a new passkey record.
func (s *PasskeyStorage) Create(ctx context.Context, passkey *model.UserPasskey) (*model.UserPasskey, error) {
	ctx, span := telemetry.StartAuto(ctx, s.Create)
	defer span.End()

	passkey.ID = ulid.String()
	now := time.Now()
	passkey.SetCreatedAt(now)

	query := `INSERT INTO user_passkey (id, user_id, credential_id, public_key, attestation_type, transport, sign_count, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, passkey.ID, passkey.UserID, passkey.CredentialID, passkey.PublicKey, passkey.AttestationType, passkey.Transport, passkey.SignCount, passkey.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create passkey: %w", err)
	}

	return passkey, nil
}

// Delete removes a passkey by ID.
func (s *PasskeyStorage) Delete(ctx context.Context, id string) error {
	ctx, span := telemetry.StartAuto(ctx, s.Delete)
	defer span.End()

	query := `DELETE FROM user_passkey WHERE id=?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete passkey: %w", err)
	}
	return nil
}

// ListByUser returns all passkeys for a given user.
func (s *PasskeyStorage) ListByUser(ctx context.Context, userID string) ([]model.UserPasskey, error) {
	ctx, span := telemetry.StartAuto(ctx, s.ListByUser)
	defer span.End()

	var passkeys []model.UserPasskey
	query := `SELECT * FROM user_passkey WHERE user_id=?`
	if err := s.db.SelectContext(ctx, &passkeys, query, userID); err != nil {
		return nil, fmt.Errorf("list passkeys: %w", err)
	}
	return passkeys, nil
}

// GetByCredentialID finds a passkey by its WebAuthn credential ID.
func (s *PasskeyStorage) GetByCredentialID(ctx context.Context, credentialID []byte) (*model.UserPasskey, error) {
	ctx, span := telemetry.StartAuto(ctx, s.GetByCredentialID)
	defer span.End()

	passkey := &model.UserPasskey{}
	query := `SELECT * FROM user_passkey WHERE credential_id=?`
	if err := s.db.GetContext(ctx, passkey, query, credentialID); err != nil {
		return nil, fmt.Errorf("get passkey by credential_id: %w", err)
	}
	return passkey, nil
}

// UpdateSignCount updates the sign count for a passkey.
func (s *PasskeyStorage) UpdateSignCount(ctx context.Context, id string, signCount int64) error {
	ctx, span := telemetry.StartAuto(ctx, s.UpdateSignCount)
	defer span.End()

	query := `UPDATE user_passkey SET sign_count=? WHERE id=?`
	_, err := s.db.ExecContext(ctx, query, signCount, id)
	if err != nil {
		return fmt.Errorf("update passkey sign_count: %w", err)
	}
	return nil
}

var _ model.PasskeyStorage = (*PasskeyStorage)(nil)
