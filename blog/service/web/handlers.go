package web

import (
	"io/fs"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog/storage"
	"github.com/titpetric/platform-app/blog/view"
)

// Handlers provides HTTP handlers for blog web endpoints.
type Handlers struct {
	repository *storage.Storage
	views      *view.Views
	themeFS    fs.FS
}

// NewHandlers returns a new Handlers instance.
func NewHandlers(repo *storage.Storage, themeFS fs.FS) *Handlers {
	return &Handlers{
		repository: repo,
		views:      view.NewViews(themeFS),
		themeFS:    themeFS,
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

// Mount registers the blog web routes on the given router.
func (h *Handlers) Mount(r platform.Router) {
	// Register static assets
	h.registerAssets(r)

	r.Group(func(r platform.Router) {
		// Public HTML Routes
		r.Get("/", h.IndexHTML)
		r.Get("/blog", h.ListArticlesHTML)
		r.Get("/blog/", h.ListArticlesHTML)
		r.Get("/blog/{slug}", h.GetArticleHTML)
		r.Get("/blog/{slug}/", h.GetArticleHTML)

		// Feed Routes
		r.Get("/feed.xml", h.GetAtomFeed)

		// Admin HTML Routes
		r.Get("/admin/articles", h.ListArticlesAdminHTML)
	})
}
