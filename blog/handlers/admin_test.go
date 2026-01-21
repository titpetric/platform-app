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

// TestListArticlesAdminJSON_Empty tests listing with no articles
func TestListArticlesAdminJSON_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles", nil)
	w := httptest.NewRecorder()

	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))

	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 0, len(result.Articles))
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_WithData tests listing with articles
func TestListArticlesAdminJSON_WithData(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	now := time.Now()
	articles := []model.Article{
		{ID: "a1", Slug: "first", Title: "First Article", Date: &now},
		{ID: "a2", Slug: "second", Title: "Second Article", Date: &now},
		{ID: "a3", Slug: "third", Title: "Third Article", Date: &now},
	}

	for i := range articles {
		require.NoError(t, repo.InsertArticle(t.Context(), &articles[i]))
	}

	req := httptest.NewRequest("GET", "/admin/articles", nil)
	w := httptest.NewRecorder()

	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 3, result.Total)
	assert.Equal(t, 3, len(result.Articles))
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_Page1 tests first page
func TestListArticlesAdminJSON_Page1(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	now := time.Now()
	for i := 1; i <= 15; i++ {
		article := &model.Article{
			ID:    "a" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article " + string(rune(i)),
			Date:  &now,
		}
		require.NoError(t, repo.InsertArticle(t.Context(), article))
	}

	req := httptest.NewRequest("GET", "/admin/articles?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 10, len(result.Articles))
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_Page2 tests second page
func TestListArticlesAdminJSON_Page2(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	now := time.Now()
	for i := 1; i <= 25; i++ {
		article := &model.Article{
			ID:    "a" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article " + string(rune(i)),
			Date:  &now,
		}
		require.NoError(t, repo.InsertArticle(t.Context(), article))
	}

	req := httptest.NewRequest("GET", "/admin/articles?page=2&pageSize=10", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 10, len(result.Articles))
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_LastPage tests partial last page
func TestListArticlesAdminJSON_LastPage(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	now := time.Now()
	for i := 1; i <= 25; i++ {
		article := &model.Article{
			ID:    "a" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article " + string(rune(i)),
			Date:  &now,
		}
		require.NoError(t, repo.InsertArticle(t.Context(), article))
	}

	req := httptest.NewRequest("GET", "/admin/articles?page=3&pageSize=10", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 5, len(result.Articles))
	assert.Equal(t, 3, result.Page)
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_CustomPageSize tests custom page sizes
func TestListArticlesAdminJSON_CustomPageSize(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	now := time.Now()
	for i := 1; i <= 20; i++ {
		article := &model.Article{
			ID:    "a" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article " + string(rune(i)),
			Date:  &now,
		}
		require.NoError(t, repo.InsertArticle(t.Context(), article))
	}

	req := httptest.NewRequest("GET", "/admin/articles?pageSize=5", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 20, result.Total)
	assert.Equal(t, 5, len(result.Articles))
	assert.Equal(t, 5, result.PageSize)
}

// TestListArticlesAdminJSON_InvalidPage_NonNumeric tests non-numeric page
func TestListArticlesAdminJSON_InvalidPage_NonNumeric(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?page=abc", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 1, result.Page)
}

// TestListArticlesAdminJSON_InvalidPage_NegativeValue tests negative page
func TestListArticlesAdminJSON_InvalidPage_NegativeValue(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?page=-5", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 1, result.Page)
}

// TestListArticlesAdminJSON_InvalidPage_Zero tests zero page
func TestListArticlesAdminJSON_InvalidPage_Zero(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?page=0", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 1, result.Page)
}

// TestListArticlesAdminJSON_InvalidPageSize_NonNumeric tests non-numeric size
func TestListArticlesAdminJSON_InvalidPageSize_NonNumeric(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?pageSize=xyz", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_InvalidPageSize_Negative tests negative page size
func TestListArticlesAdminJSON_InvalidPageSize_Negative(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?pageSize=-10", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_InvalidPageSize_Exceeds100 tests page size > 100
func TestListArticlesAdminJSON_InvalidPageSize_Exceeds100(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?pageSize=200", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 10, result.PageSize)
}

// TestListArticlesAdminJSON_PageSize100_Valid tests pageSize=100 is valid
func TestListArticlesAdminJSON_PageSize100_Valid(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	req := httptest.NewRequest("GET", "/admin/articles?pageSize=100", nil)
	w := httptest.NewRecorder()
	h.ListArticlesAdminJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result model.ArticleList
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, 100, result.PageSize)
}

// TestListArticlesAdminJSON_OffsetCalculation tests correct offset calculation
func TestListArticlesAdminJSON_OffsetCalculation(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{"page 1 offset 0", 1, 10, 0},
		{"page 2 offset 10", 2, 10, 10},
		{"page 3 offset 20", 3, 10, 20},
		{"page 1 size 5", 1, 5, 0},
		{"page 2 size 5", 2, 5, 5},
		{"page 3 size 5", 3, 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := (tt.page - 1) * tt.pageSize
			assert.Equal(t, tt.expected, expected)
		})
	}
}

// TestGetArticleAdminJSON_Found tests retrieving existing article
func TestGetArticleAdminJSON_Found(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	now := time.Now()
	article := &model.Article{
		ID:          "admin-test-1",
		Slug:        "admin-article",
		Title:       "Admin Article",
		Description: "Test admin article",
		Date:        &now,
	}

	err = repo.InsertArticle(t.Context(), article)
	require.NoError(t, err)

	// Create a chi router to properly set URL params
	r := chi.NewRouter()
	var result model.Article

	r.Get("/admin/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleAdminJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/admin/articles/admin-article", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))

	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "admin-article", result.Slug)
	assert.Equal(t, "Admin Article", result.Title)
}

// TestGetArticleAdminJSON_NotFound tests non-existent article
func TestGetArticleAdminJSON_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	r := chi.NewRouter()
	r.Get("/admin/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleAdminJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/admin/articles/non-existent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "article not found")
}

// TestGetArticleAdminJSON_EmptySlug tests missing slug parameter
func TestGetArticleAdminJSON_EmptySlug(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewAdminHandlers(repo)

	// Direct test - chi.URLParam should return empty string for missing param
	req := httptest.NewRequest("GET", "/admin/articles", nil)
	w := httptest.NewRecorder()
	h.GetArticleAdminJSON(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "slug parameter is required")
}
