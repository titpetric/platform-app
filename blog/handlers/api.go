package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform-app/blog/model"
)

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
	slug := chi.URLParam(r, "slug")

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
