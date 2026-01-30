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

type Storage struct {
	db *sqlx.DB
}

func NewStorage(ctx context.Context) (*Storage, error) {
	db, err := DB(ctx)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Pulse(ctx context.Context, count int64) error {
	user, active := user.GetSessionUser(ctx)
	if !active {
		return errors.New("invalid user auth")
	}

	return platform.Transaction(ctx, s.db, s.pulseFn(user.ID, count))
}

const updatePulseHourly = `
INSERT INTO
  pulse_hourly (user_id, stamp, count)
VALUES
  (:user_id, strftime('%Y-%m-%d %H:00:00', 'now'), :count)
ON
  CONFLICT(user_id, stamp)
DO
  UPDATE SET count = count + :count;`

const updatePulseDaily = `
INSERT INTO
  pulse_daily (user_id, stamp, count)
VALUES
  (:user_id, strftime('%Y-%m-%d', 'now'), :count)
ON
  CONFLICT(user_id, stamp)
DO
  UPDATE SET count = count + :count;
`

func (s *Storage) pulseFn(userID string, count int64) func(context.Context, *sqlx.Tx) error {
	params := map[string]any{
		"user_id": userID,
	}

	return func(ctx context.Context, tx *sqlx.Tx) error {
		query := strings.ReplaceAll(updatePulseHourly, ":count", fmt.Sprint(count))
		if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
			return err
		}

		query = strings.ReplaceAll(updatePulseDaily, ":count", fmt.Sprint(count))
		if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
			return err
		}
		return nil
	}
}
