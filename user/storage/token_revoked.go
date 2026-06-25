package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/pkg/telemetry"
)

// RevokedTokenStorage records JWT IDs (jti) that have been explicitly
// revoked before their natural expiry. It is consulted by the auth
// middleware and the refresh-token endpoint.
type RevokedTokenStorage struct {
	db *sqlx.DB
}

// NewRevokedTokenStorage returns a new RevokedTokenStorage.
func NewRevokedTokenStorage(db *sqlx.DB) *RevokedTokenStorage {
	return &RevokedTokenStorage{db: db}
}

// Revoke marks the given jti as revoked. expiresAt is the JWT's original
// `exp` claim and is stored so expired rows can be cleaned up.
// Calling Revoke on an already-revoked jti is a no-op.
func (s *RevokedTokenStorage) Revoke(ctx context.Context, jti, userID string, expiresAt time.Time) error {
	if s == nil || s.db == nil {
		return nil
	}
	ctx, span := telemetry.StartAuto(ctx, s.Revoke)
	defer span.End()

	if jti == "" {
		return errors.New("revoke: empty jti")
	}

	query := `INSERT OR IGNORE INTO user_token_revoked (jti, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	if _, err := s.db.ExecContext(ctx, query, jti, userID, expiresAt, time.Now()); err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	return nil
}

// IsRevoked reports whether the given jti has been revoked. An empty jti
// is treated as not revoked so tokens issued before JTI rollout continue
// to validate.
func (s *RevokedTokenStorage) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if s == nil || s.db == nil || jti == "" {
		return false, nil
	}
	ctx, span := telemetry.StartAuto(ctx, s.IsRevoked)
	defer span.End()

	var found string
	err := s.db.GetContext(ctx, &found, `SELECT jti FROM user_token_revoked WHERE jti=? LIMIT 1`, jti)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("is revoked: %w", err)
	}
	return true, nil
}

// PurgeExpired deletes revocation rows whose stored exp has already passed.
// Such rows can be removed because the underlying JWT can no longer
// validate regardless of revocation state.
func (s *RevokedTokenStorage) PurgeExpired(ctx context.Context) (int64, error) {
	if s == nil || s.db == nil {
		return 0, nil
	}
	ctx, span := telemetry.StartAuto(ctx, s.PurgeExpired)
	defer span.End()

	res, err := s.db.ExecContext(ctx, `DELETE FROM user_token_revoked WHERE expires_at IS NOT NULL AND expires_at < ?`, time.Now())
	if err != nil {
		return 0, fmt.Errorf("purge revoked: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}
