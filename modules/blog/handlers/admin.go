package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform-app/modules/blog/model"
)

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

// ListArticlesAdminHTML renders the admin articles list as HTML using CMS layout
func (h *Handlers) ListArticlesAdminHTML(w http.ResponseWriter, r *http.Request) {
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

	// Prepare data for template rendering
	data := map[string]interface{}{
		"title":       "Articles - Blog Admin",
		"total":       total,
		"page":        pageNum,
		"pageSize":    pageSz,
		"breadcrumbs": []map[string]string{},
		"articles":    articles,
	}

	// Render using layout renderer with blog_admin_articles.vuego
	err = h.layoutRenderer.Render(r.Context(), w, "blog_admin_articles.vuego", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}
