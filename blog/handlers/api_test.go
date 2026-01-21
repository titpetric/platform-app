package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/storage"
)

// TestListArticlesJSON_Empty tests listing all articles with none present
func TestListArticlesJSON_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/api/blog/articles", nil)
	w := httptest.NewRecorder()

	h.ListArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=300", w.Header().Get("Cache-Control"))

	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, len(result.Articles))
	assert.Equal(t, 1, result.Page)
}

// TestListArticlesJSON_MultipleArticles tests listing multiple articles
func TestListArticlesJSON_MultipleArticles(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	now := time.Now()
	for i := 1; i <= 5; i++ {
		article := &model.Article{
			ID:    "api-a" + string(rune(i)),
			Slug:  "api-article-" + string(rune(i)),
			Title: "API Article " + string(rune(i)),
			Date:  &now,
		}
		err = repo.InsertArticle(ctx, article)
		require.NoError(t, err)
	}

	req := httptest.NewRequest("GET", "/api/blog/articles", nil)
	w := httptest.NewRecorder()

	h.ListArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 5, result.Total)
	assert.Equal(t, 5, len(result.Articles))
	assert.Equal(t, 1, result.Page)
}

// TestListArticlesJSON_HeadersPresent validates required headers
func TestListArticlesJSON_HeadersPresent(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/api/blog/articles", nil)
	w := httptest.NewRecorder()

	h.ListArticlesJSON(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=300", w.Header().Get("Cache-Control"))
}

// TestGetArticleJSON_Found tests retrieving a single article by slug
func TestGetArticleJSON_Found(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	now := time.Now()
	article := &model.Article{
		ID:          "api-single-1",
		Slug:        "api-single-article",
		Title:       "API Single Article",
		Description: "This is a test article",
		OgImage:     "/images/og.png",
		Date:        &now,
	}

	err = repo.InsertArticle(ctx, article)
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Get("/api/blog/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/api/blog/articles/api-single-article", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=3600", w.Header().Get("Cache-Control"))

	var result model.Article
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "api-single-article", result.Slug)
	assert.Equal(t, "API Single Article", result.Title)
	assert.Equal(t, "This is a test article", result.Description)
}

// TestGetArticleJSON_NotFound tests 404 for non-existent article
func TestGetArticleJSON_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	r := chi.NewRouter()
	r.Get("/api/blog/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/api/blog/articles/missing-article", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "article not found")
}

// TestGetArticleJSON_CacheHeaders validates caching strategy
func TestGetArticleJSON_CacheHeaders(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	// Insert an article so we get a 200 response
	now := time.Now()
	article := &model.Article{
		ID:    "cache-test",
		Slug:  "cache-test",
		Title: "Cache Test",
		Date:  &now,
	}
	err = repo.InsertArticle(ctx, article)
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Get("/api/blog/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/api/blog/articles/cache-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "public, max-age=3600", w.Header().Get("Cache-Control"))
}

// TestSearchArticlesJSON_MissingQuery tests missing search query
func TestSearchArticlesJSON_MissingQuery(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/api/blog/search", nil)
	w := httptest.NewRecorder()

	h.SearchArticlesJSON(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing 'q' query parameter")
}

// TestSearchArticlesJSON_EmptyResults tests search with no matches
func TestSearchArticlesJSON_EmptyResults(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	now := time.Now()
	article := &model.Article{
		ID:    "search-1",
		Slug:  "golang-basics",
		Title: "Go Basics",
		Date:  &now,
	}
	err = repo.InsertArticle(ctx, article)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/blog/search?q=python", nil)
	w := httptest.NewRecorder()

	h.SearchArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, 0, int(result["total"].(float64)))
	assert.Equal(t, "python", result["query"])
}

// TestSearchArticlesJSON_SingleMatch tests search with one match
func TestSearchArticlesJSON_SingleMatch(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	now := time.Now()
	articles := []model.Article{
		{
			ID:          "search-2",
			Slug:        "go-tutorial",
			Title:       "Go Tutorial",
			Description: "Learn Go",
			Date:        &now,
		},
		{
			ID:          "search-3",
			Slug:        "rust-guide",
			Title:       "Rust Guide",
			Description: "Learn Rust",
			Date:        &now,
		},
	}

	for i := range articles {
		err = repo.InsertArticle(ctx, &articles[i])
		require.NoError(t, err)
	}

	req := httptest.NewRequest("GET", "/api/blog/search?q=Go", nil)
	w := httptest.NewRecorder()

	h.SearchArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 1, int(result["total"].(float64)))
	assert.Equal(t, "Go", result["query"])
}

// TestSearchArticlesJSON_MultipleMatches tests search with multiple matches
func TestSearchArticlesJSON_MultipleMatches(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	now := time.Now()
	articles := []model.Article{
		{ID: "s1", Slug: "go-1", Title: "Go Guide", Date: &now},
		{ID: "s2", Slug: "go-2", Title: "Go Tips", Date: &now},
		{ID: "s3", Slug: "rust-1", Title: "Rust Basics", Date: &now},
	}

	for i := range articles {
		require.NoError(t, repo.InsertArticle(ctx, &articles[i]))
	}

	req := httptest.NewRequest("GET", "/api/blog/search?q=Go", nil)
	w := httptest.NewRecorder()

	h.SearchArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 2, int(result["total"].(float64)))
}

// TestSearchArticlesJSON_CacheHeaders validates cache control
func TestSearchArticlesJSON_CacheHeaders(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/api/blog/search?q=test", nil)
	w := httptest.NewRecorder()

	h.SearchArticlesJSON(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "public, max-age=300", w.Header().Get("Cache-Control"))
}

// TestSearchArticlesJSON_CaseInsensitive tests case-insensitive search
func TestSearchArticlesJSON_CaseInsensitive(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)
	ctx := t.Context()

	now := time.Now()
	article := &model.Article{
		ID:          "search-case",
		Slug:        "programming-101",
		Title:       "Programming Basics",
		Description: "Introduction to Programming",
		Date:        &now,
	}
	err = repo.InsertArticle(ctx, article)
	require.NoError(t, err)

	// Search with different case
	req := httptest.NewRequest("GET", "/api/blog/search?q=programming", nil)
	w := httptest.NewRecorder()
	h.SearchArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(result["total"].(float64)))
}
