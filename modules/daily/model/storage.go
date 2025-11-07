package model

import "context"

// Storage defines operations against the todo store.
//
// Note: Add and Update accept a Todo value so we can pass title/completed
// and timestamps from Go. Methods always take a context.
type Storage interface {
	// List returns todos.
	List(ctx context.Context) ([]Todo, error)

	// Get returns a todo by id.
	Get(ctx context.Context, id string) (Todo, error)

	// Add inserts a new todo. The returned Todo contains the generated ID and timestamps.
	Add(ctx context.Context, t Todo) (Todo, error)

	// Update updates title/completed/updated_at of the todo.
	Update(ctx context.Context, t Todo) error

	// Complete marks the todo identified by id as completed and updates updated_at.
	Complete(ctx context.Context, id string) error

	// Delete soft-deletes the todo (sets deleted_at).
	Delete(ctx context.Context, id string) error
}
