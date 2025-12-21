package handlers

import (
	"io/fs"

	"github.com/titpetric/vuego"

	"github.com/titpetric/platform-app/modules/blog/storage"
	"github.com/titpetric/platform-app/modules/blog/view"
)

// Handlers handles HTTP requests for the blog module
type Handlers struct {
	repository *storage.Storage
	views      *view.Views
}

// NewHandlers creates a new Handlers instance with the given storage and theme
func NewHandlers(repo *storage.Storage, themeFS fs.FS) (*Handlers, error) {
	return &Handlers{
		repository: repo,
		views:      view.NewViews(vuego.New(vuego.WithFS(themeFS), vuego.WithLessProcessor())),
	}, nil
}

// NewAdminHandlers creates a new admin-only Handlers instance
func NewAdminHandlers(repo *storage.Storage) *Handlers {
	return &Handlers{
		repository: repo,
		views:      nil,
	}
}

// Repository returns the storage repository
func (h *Handlers) Repository() *storage.Storage {
	return h.repository
}

// Views returns the views instance
func (h *Handlers) Views() *view.Views {
	return h.views
}
