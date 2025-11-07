package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/platform/pkg/ulid"

	"github.com/titpetric/platform-app/modules/daily/model"
)

// Storage implements model.Storage backed by sqlite
type Storage struct {
	db *sqlx.DB
}

// NewStorage creates a Storage using an existing sqlx.DB
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// List returns all non-deleted todos
func (s *Storage) List(ctx context.Context) ([]model.Todo, error) {
	var todos []model.Todo
	err := s.db.SelectContext(ctx, &todos, `
		SELECT id, title, completed, created_at, updated_at, deleted_at
		FROM `+model.TodoTable+`
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	return todos, err
}

// Get returns a single todo by id
func (s *Storage) Get(ctx context.Context, id string) (model.Todo, error) {
	var t model.Todo
	err := s.db.GetContext(ctx, &t, `
		SELECT id, title, completed, created_at, updated_at, deleted_at
		FROM `+model.TodoTable+`
		WHERE id = ?
		LIMIT 1
	`, id)
	return t, err
}

// Add inserts a new todo
func (s *Storage) Add(ctx context.Context, t model.Todo) (model.Todo, error) {
	now := time.Now().UTC()
	t.SetCreatedAt(now)
	t.SetUpdatedAt(now)
	t.ID = ulid.String()

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO `+model.TodoTable+` (id, title, completed, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, t.ID, t.Title, t.Completed, t.CreatedAt, t.UpdatedAt, t.DeletedAt)
	if err != nil {
		return model.Todo{}, err
	}
	return s.Get(ctx, t.ID)
}

// Update modifies title/completed/updated_at for a todo
func (s *Storage) Update(ctx context.Context, id string, t model.Todo) (model.Todo, error) {
	if id == "" {
		return model.Todo{}, errors.New("id required")
	}
	now := time.Now().UTC()
	t.SetUpdatedAt(now)

	res, err := s.db.ExecContext(ctx, `
		UPDATE `+model.TodoTable+`
		SET title = ?, completed = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`, t.Title, t.Completed, t.UpdatedAt, id)
	if err != nil {
		return model.Todo{}, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return model.Todo{}, sql.ErrNoRows
	}
	return s.Get(ctx, id)
}

// Complete marks a todo completed and updates updated_at
func (s *Storage) Complete(ctx context.Context, id string) (model.Todo, error) {
	if id == "" {
		return model.Todo{}, errors.New("id required")
	}
	now := time.Now().UTC()

	res, err := s.db.ExecContext(ctx, `
		UPDATE `+model.TodoTable+`
		SET completed = 1, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`, now, id)
	if err != nil {
		return model.Todo{}, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return model.Todo{}, sql.ErrNoRows
	}
	return s.Get(ctx, id)
}

// Delete soft-deletes a todo by setting deleted_at
func (s *Storage) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	now := time.Now().UTC()
	_, err := s.db.ExecContext(ctx, `
		UPDATE `+model.TodoTable+`
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`, now, id)
	if err != nil {
		return err
	}
	return nil
}

// helper: bool -> int for sqlite
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
