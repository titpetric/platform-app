package handlers

import (
	"io/fs"

	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog/storage"
)

// RegisterRoutes registers all blog routes (public and admin)
func RegisterRoutes(r platform.Router, repository *storage.Storage, themeFS fs.FS) error {
	h, err := NewHandlers(repository, themeFS)
	if err != nil {
		return err
	}

	// Register static assets
	RegisterAssets(r, themeFS)

	r.Group(func(r platform.Router) {
		// Public API Routes (JSON)
		r.Get("/api/blog/articles", h.ListArticlesJSON)
		r.Get("/api/blog/articles/{slug}", h.GetArticleJSON)
		r.Get("/api/blog/search", h.SearchArticlesJSON)

		// Public HTML Routes
		r.Get("/", h.IndexHTML)
		r.Get("/blog", h.ListArticlesHTML)
		r.Get("/blog/", h.ListArticlesHTML)
		r.Get("/blog/{slug}", h.GetArticleHTML)
		r.Get("/blog/{slug}/", h.GetArticleHTML)

		// Feed Routes
		r.Get("/feed.xml", h.GetAtomFeed)

		// Admin Routes
		r.Get("/admin/articles", h.ListArticlesAdminHTML)
		r.Get("/admin/articles.json", h.ListArticlesAdminJSON)
		r.Get("/admin/articles/{slug}", h.GetArticleAdminJSON)
	})

	return nil
}
