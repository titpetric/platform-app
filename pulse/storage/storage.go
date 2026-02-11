// Package storage provides database persistence for pulse data.
package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/user"
)

// Storage provides pulse data persistence.
type Storage struct {
	db *sqlx.DB
}

// NewStorage creates a new storage backed by the given database.
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Pulse records keystroke activity for the authenticated user.
func (s *Storage) Pulse(ctx context.Context, count int64, hostname string) error {
	user, active := user.GetSessionUser(ctx)
	if !active {
		return errors.New("invalid user auth")
	}

	return platform.Transaction(ctx, s.db, s.pulseFn(user.ID, count, hostname))
}

// https://github.com/jmoiron/sqlx/issues/368
// - uses :: to escape non-named params
// - replaces :count to literal int
const updatePulseHourly = `
INSERT INTO
  pulse_hourly (user_id, hostname, stamp, count)
VALUES
  (:user_id, :hostname, strftime('%Y-%m-%d %H::00::00', 'now'), :count)
ON
  CONFLICT(user_id, hostname, stamp)
DO
  UPDATE SET count = count + :count`

const updatePulseDaily = `
INSERT INTO
  pulse_daily (user_id, hostname, stamp, count)
VALUES
  (:user_id, :hostname, strftime('%Y-%m-%d', 'now'), :count)
ON
  CONFLICT(user_id, hostname, stamp)
DO
  UPDATE SET count = count + :count`

const updatePulseHosts = `INSERT OR IGNORE INTO pulse_hosts (user_id, hostname, created_at) VALUES (:user_id, :hostname, CURRENT_TIMESTAMP)`

func (s *Storage) pulseFn(userID string, count int64, hostname string) func(context.Context, *sqlx.Tx) error {
	params := map[string]any{
		"user_id":  userID,
		"hostname": hostname,
	}

	return func(ctx context.Context, tx *sqlx.Tx) error {
		query := strings.ReplaceAll(updatePulseHourly, ":count", fmt.Sprint(count))
		if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
			return fmt.Errorf("error in %s: %w", query, err)
		}

		query = strings.ReplaceAll(updatePulseDaily, ":count", fmt.Sprint(count))
		if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
			return fmt.Errorf("error in %s: %w", query, err)
		}

		query = updatePulseHosts
		if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
			return fmt.Errorf("error in %s: %w", query, err)
		}

		return nil
	}
}
