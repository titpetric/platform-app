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

// UserCount holds a user's total keystroke count.
type UserCount struct {
	UserID string `db:"user_id"`
	Count  int64  `db:"count"`
}

// ListUserCounts returns the total keystroke count per user for the last 30 days.
func (s *Storage) ListUserCounts(ctx context.Context) ([]UserCount, error) {
	var counts []UserCount
	query := `SELECT user_id, SUM(count) as count FROM pulse_daily WHERE stamp >= date('now', '-30 days') GROUP BY user_id ORDER BY count DESC`
	if err := s.db.SelectContext(ctx, &counts, query); err != nil {
		return nil, fmt.Errorf("list user counts: %w", err)
	}
	return counts, nil
}

// HourlyCount holds hourly keystroke data.
type HourlyCount struct {
	Hour  int   `db:"hour" json:"hour"`
	Count int64 `db:"count" json:"count"`
}

// GetUserHourly returns hourly keystroke distribution for a user (aggregated across all days).
func (s *Storage) GetUserHourly(ctx context.Context, userID string) ([]HourlyCount, error) {
	var counts []HourlyCount
	query := `
		SELECT
			CAST(strftime('%H', stamp) AS INTEGER) as hour,
			SUM(count) as count
		FROM pulse_hourly
		WHERE user_id = ?
		GROUP BY hour
		ORDER BY hour`
	if err := s.db.SelectContext(ctx, &counts, query, userID); err != nil {
		return nil, fmt.Errorf("get user hourly: %w", err)
	}
	return counts, nil
}

// DailyHostCount holds daily keystroke data per host.
type DailyHostCount struct {
	Hostname string `db:"hostname" json:"hostname"`
	Stamp    string `db:"stamp" json:"stamp"`
	Count    int64  `db:"count" json:"count"`
}

// GetUserDaily returns daily keystroke counts per host for a user over the last 30 days.
func (s *Storage) GetUserDaily(ctx context.Context, userID string) ([]DailyHostCount, error) {
	var counts []DailyHostCount
	query := `
		SELECT hostname, stamp, count
		FROM pulse_daily
		WHERE user_id = ? AND stamp >= date('now', '-30 days')
		ORDER BY hostname, stamp`
	if err := s.db.SelectContext(ctx, &counts, query, userID); err != nil {
		return nil, fmt.Errorf("get user daily: %w", err)
	}
	return counts, nil
}

// GetUserHourlyByHost returns hourly keystroke counts per host for a user over the last 48 hours.
func (s *Storage) GetUserHourlyByHost(ctx context.Context, userID string, hostname string) ([]DailyHostCount, error) {
	var counts []DailyHostCount
	query := `
		SELECT hostname, stamp, count
		FROM pulse_hourly
		WHERE user_id = ? AND hostname = ? AND stamp >= datetime('now', '-48 hours')
		ORDER BY stamp`
	if err := s.db.SelectContext(ctx, &counts, query, userID, hostname); err != nil {
		return nil, fmt.Errorf("get user hourly by host: %w", err)
	}
	return counts, nil
}

// GetUserHosts returns all hostnames for a user.
func (s *Storage) GetUserHosts(ctx context.Context, userID string) ([]string, error) {
	var hosts []string
	query := `SELECT hostname FROM pulse_hosts WHERE user_id = ? ORDER BY hostname`
	if err := s.db.SelectContext(ctx, &hosts, query, userID); err != nil {
		return nil, fmt.Errorf("get user hosts: %w", err)
	}
	return hosts, nil
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
