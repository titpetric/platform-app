package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/storage"
)

// Handlers provides HTTP handlers for blog API endpoints.
type Handlers struct {
	repository *storage.Storage
}

// NewHandlers returns a new Handlers instance.
func NewHandlers(repo *storage.Storage) *Handlers {
	return &Handlers{
		repository: repo,
	}
}

// Mount registers the blog API routes on the given router.
func (h *Handlers) Mount(r platform.Router) {
	r.Group(func(r platform.Router) {
		// Public API Routes (JSON)
		r.Get("/api/blog/articles", h.ListArticlesJSON)
		r.Get("/api/blog/articles/{slug}", h.GetArticleJSON)
		r.Get("/api/blog/search", h.SearchArticlesJSON)

		// Admin Routes (JSON)
		r.Get("/admin/articles.json", h.ListArticlesAdminJSON)
		r.Get("/admin/articles/{slug}", h.GetArticleAdminJSON)
	})
}

// ListArticlesJSON returns a JSON list of all articles
func (h *Handlers) ListArticlesJSON(w http.ResponseWriter, r *http.Request) {
	articles, err := h.repository.GetArticles(r.Context(), 0, 9999)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch articles: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")

	list := &model.ArticleList{
		Articles: articles,
		Total:    len(articles),
		Page:     1,
		PageSize: len(articles),
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetArticleJSON returns a single article as JSON
func (h *Handlers) GetArticleJSON(w http.ResponseWriter, r *http.Request) {
	slug := platform.URLParam(r, "slug")

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, fmt.Sprintf("article not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SearchArticlesJSON performs full-text search on articles
func (h *Handlers) SearchArticlesJSON(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing 'q' query parameter", http.StatusBadRequest)
		return
	}

	articles, err := h.repository.SearchArticles(r.Context(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf("search failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")

	result := map[string]interface{}{
		"articles": articles,
		"total":    len(articles),
		"query":    query,
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListArticlesAdminJSON returns a paginated JSON list of articles for admin
func (h *Handlers) ListArticlesAdminJSON(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	// Default pagination values
	pageNum := 1
	pageSz := 10

	// Parse and validate page parameter
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	// Parse and validate pageSize parameter
	if pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			pageSz = ps
		}
	}

	// Calculate offset
	offset := (pageNum - 1) * pageSz

	// Get total count
	total, err := h.repository.CountArticles(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to count articles: %v", err), http.StatusInternalServerError)
		return
	}

	// Get paginated articles
	articles, err := h.repository.GetArticles(r.Context(), offset, pageSz)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch articles: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	// Build article list response
	list := &model.ArticleList{
		Articles: articles,
		Total:    total,
		Page:     pageNum,
		PageSize: pageSz,
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetArticleAdminJSON returns a single article as JSON for admin
func (h *Handlers) GetArticleAdminJSON(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug parameter is required", http.StatusBadRequest)
		return
	}

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, fmt.Sprintf("article not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	if err := json.NewEncoder(w).Encode(article); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
