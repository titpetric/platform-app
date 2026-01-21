package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/pkg/telemetry"
	"github.com/titpetric/platform/pkg/ulid"

	"github.com/titpetric/platform-app/daily/model"
	"github.com/titpetric/platform-app/daily/schema"
	"github.com/titpetric/platform-app/user"
)

// Storage implements model.Storage backed by sqlite
type Storage struct {
	db *sqlx.DB
}

// New creates a new Storage instance, runs migrations.
func New(ctx context.Context) (*Storage, error) {
	db, err := DB(ctx)
	if err != nil {
		return nil, err
	}
	return NewStorage(ctx, db)
}

// NewStorage creates a new Storage instance, rungs migrations.
func NewStorage(ctx context.Context, db *sqlx.DB) (*Storage, error) {
	if err := Migrate(ctx, db, schema.Migrations); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

var empty model.Todo

// List returns all non-deleted todos
func (s *Storage) List(ctx context.Context) ([]model.Todo, error) {
	u, ok := user.GetSessionUser(ctx)
	if !ok {
		return nil, user.ErrLoginRequired
	}

	var todos []model.Todo
	err := s.db.SelectContext(ctx, &todos, `
		SELECT id, title, completed, created_at, updated_at, deleted_at
		FROM `+model.TodoTable+`
		WHERE user_id=? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, u.ID)
	if err != nil {
		return nil, err
	}

	return todos, err
}

// Get returns a todo by ID
func (s *Storage) Get(ctx context.Context, id string) (model.Todo, error) {
	u, ok := user.GetSessionUser(ctx)
	if !ok {
		return empty, user.ErrLoginRequired
	}

	var todo model.Todo
	err := s.db.GetContext(ctx, &todo, `
		SELECT id, title, completed, created_at, updated_at, deleted_at
		FROM `+model.TodoTable+`
		WHERE id=? AND user_id=? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, id, u.ID)
	if err != nil {
		return empty, err
	}

	return todo, err
}

// Add inserts a new todo
func (s *Storage) Add(ctx context.Context, t model.Todo) (model.Todo, error) {
	u, ok := user.GetSessionUser(ctx)
	if !ok {
		return empty, user.ErrLoginRequired
	}

	now := time.Now().UTC()
	t.ID = ulid.String()
	t.UserID = u.ID
	t.SetCreatedAt(now)
	t.SetUpdatedAt(now)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO `+model.TodoTable+` (id, user_id, title, completed, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, t.ID, t.UserID, t.Title, t.Completed, t.CreatedAt, t.UpdatedAt, t.DeletedAt)
	if err != nil {
		return empty, err
	}
	return s.Get(ctx, t.ID)
}

// Update modifies title/completed/updated_at for a todo
func (s *Storage) Update(ctx context.Context, t model.Todo) error {
	u, ok := user.GetSessionUser(ctx)
	if !ok {
		return user.ErrLoginRequired
	}
	if t.ID == "" {
		return errors.New("id required")
	}
	now := time.Now().UTC()
	t.SetUpdatedAt(now)

	res, err := s.db.ExecContext(ctx, `
		UPDATE `+model.TodoTable+`
		SET title = ?, completed = ?, updated_at = ?
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, t.Title, t.Completed, t.UpdatedAt, t.ID, u.ID)

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		if err != nil {
			telemetry.CaptureError(ctx, err)
		}
		return sql.ErrNoRows
	}

	return err
}

// Complete marks a todo completed and updates updated_at
func (s *Storage) Complete(ctx context.Context, id string) error {
	u, ok := user.GetSessionUser(ctx)
	if !ok {
		return user.ErrLoginRequired
	}
	if id == "" {
		return errors.New("id required")
	}
	now := time.Now().UTC()

	res, err := s.db.ExecContext(ctx, `
		UPDATE `+model.TodoTable+`
		SET completed = 1, updated_at = ?
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, now, id, u.ID)

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		if err != nil {
			telemetry.CaptureError(ctx, err)
		}
		return sql.ErrNoRows
	}

	return err
}

// Delete soft-deletes a todo by setting deleted_at
func (s *Storage) Delete(ctx context.Context, id string) error {
	u, ok := user.GetSessionUser(ctx)
	if !ok {
		return user.ErrLoginRequired
	}
	if id == "" {
		return errors.New("id required")
	}
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE `+model.TodoTable+`
		SET deleted_at = ?
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`, now, id, u.ID)

	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		if err != nil {
			telemetry.CaptureError(ctx, err)
		}
		return sql.ErrNoRows
	}

	return err
}

// helper: bool -> int for sqlite
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
