package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	chi "github.com/go-chi/chi/v5"
	"github.com/titpetric/platform"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/storage"
)

// slugPattern matches lowercase alphanumeric slugs with hyphens.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// maxSlugLength is the maximum allowed slug length.
const maxSlugLength = 100

// isValidSlug validates that a slug is well-formed.
func isValidSlug(slug string) bool {
	if slug == "" || len(slug) > maxSlugLength {
		return false
	}
	return slugPattern.MatchString(slug)
}

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

		// Admin Routes (JSON) - moved to /api/admin/blog/*
		r.Get("/api/admin/blog/articles", h.ListArticlesAdminJSON)
		r.Get("/api/admin/blog/articles/{slug}", h.GetArticleAdminJSON)
	})
}

// ListArticlesJSON returns a JSON list of all articles.
func (h *Handlers) ListArticlesJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listArticlesJSON(w, r))
}

func (h *Handlers) listArticlesJSON(w http.ResponseWriter, r *http.Request) error {
	// Public listing returns only published articles (no drafts, no scheduled)
	articles, err := h.repository.GetPublishedArticles(r.Context(), 0, 9999)
	if err != nil {
		return ErrInternal("failed to fetch articles", err)
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
		return ErrInternal("failed to encode response", err)
	}
	return nil
}

// GetArticleJSON returns a single article as JSON.
func (h *Handlers) GetArticleJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getArticleJSON(w, r))
}

func (h *Handlers) getArticleJSON(w http.ResponseWriter, r *http.Request) error {
	slug := platform.URLParam(r, "slug")
	if !isValidSlug(slug) {
		return ErrBadRequest("invalid slug", nil)
	}

	// Public detail endpoint must not expose drafts or scheduled articles
	article, err := h.repository.GetPublishedArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	if err := json.NewEncoder(w).Encode(article); err != nil {
		return ErrInternal("failed to encode response", err)
	}
	return nil
}

// SearchArticlesJSON performs full-text search on articles.
func (h *Handlers) SearchArticlesJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.searchArticlesJSON(w, r))
}

// maxSearchQueryLength is the maximum allowed length for a search query.
const maxSearchQueryLength = 200

func (h *Handlers) searchArticlesJSON(w http.ResponseWriter, r *http.Request) error {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		return ErrBadRequest("missing 'q' query parameter", nil)
	}
	if len(query) > maxSearchQueryLength {
		return ErrBadRequest("query exceeds maximum length", nil)
	}

	// Public search must not return drafts or scheduled articles
	articles, err := h.repository.SearchPublishedArticles(r.Context(), query)
	if err != nil {
		return ErrInternal("search failed", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")

	result := map[string]interface{}{
		"articles": articles,
		"total":    len(articles),
		"query":    query,
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return ErrInternal("failed to encode response", err)
	}
	return nil
}

// ListArticlesAdminJSON returns a paginated JSON list of articles for admin.
func (h *Handlers) ListArticlesAdminJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.listArticlesAdminJSON(w, r))
}

func (h *Handlers) listArticlesAdminJSON(w http.ResponseWriter, r *http.Request) error {
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
		return ErrInternal("failed to count articles", err)
	}

	// Get paginated articles
	articles, err := h.repository.GetArticles(r.Context(), offset, pageSz)
	if err != nil {
		return ErrInternal("failed to fetch articles", err)
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
		return ErrInternal("failed to encode response", err)
	}
	return nil
}

// GetArticleAdminJSON returns a single article as JSON for admin.
func (h *Handlers) GetArticleAdminJSON(w http.ResponseWriter, r *http.Request) {
	h.errorHandler(w, r, h.getArticleAdminJSON(w, r))
}

func (h *Handlers) getArticleAdminJSON(w http.ResponseWriter, r *http.Request) error {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		return ErrBadRequest("slug parameter is required", nil)
	}
	if !isValidSlug(slug) {
		return ErrBadRequest("invalid slug", nil)
	}

	article, err := h.repository.GetArticleBySlug(r.Context(), slug)
	if err != nil {
		return ErrNotFound("article not found", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	if err := json.NewEncoder(w).Encode(article); err != nil {
		return ErrInternal("failed to encode response", err)
	}
	return nil
}
