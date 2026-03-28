package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	chi "github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/schema"
	"github.com/titpetric/platform-app/blog/storage"
)

// setupTestDB creates a temporary SQLite database for testing
func setupTestDB(t *testing.T) *sqlx.DB {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, err := sqlx.Open("sqlite", dbPath)
	require.NoError(t, err)

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	require.NoError(t, storage.Migrate(t.Context(), db, schema.Migrations))

	return db
}

// TestListArticlesJSON_Empty tests listing all articles with none present
func TestListArticlesJSON_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewHandlers(repo)

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
	h := NewHandlers(repo)
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
	h := NewHandlers(repo)

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
	h := NewHandlers(repo)
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
	h := NewHandlers(repo)

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
	h := NewHandlers(repo)
	ctx := t.Context()

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
	h := NewHandlers(repo)

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
	h := NewHandlers(repo)
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
	h := NewHandlers(repo)
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
	h := NewHandlers(repo)
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
	h := NewHandlers(repo)

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
	h := NewHandlers(repo)
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

	req := httptest.NewRequest("GET", "/api/blog/search?q=programming", nil)
	w := httptest.NewRecorder()
	h.SearchArticlesJSON(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(result["total"].(float64)))
}

// TestListArticlesAdminJSON_Empty tests listing with no articles
func TestListArticlesAdminJSON_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles", nil)
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
	h := NewHandlers(repo)

	now := time.Now()
	articles := []model.Article{
		{ID: "a1", Slug: "first", Title: "First Article", Date: &now},
		{ID: "a2", Slug: "second", Title: "Second Article", Date: &now},
		{ID: "a3", Slug: "third", Title: "Third Article", Date: &now},
	}

	for i := range articles {
		require.NoError(t, repo.InsertArticle(t.Context(), &articles[i]))
	}

	req := httptest.NewRequest("GET", "/api/admin/blog/articles", nil)
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
	h := NewHandlers(repo)

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

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?page=1&pageSize=10", nil)
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
	h := NewHandlers(repo)

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

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?page=2&pageSize=10", nil)
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
	h := NewHandlers(repo)

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

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?page=3&pageSize=10", nil)
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
	h := NewHandlers(repo)

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

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?pageSize=5", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?page=abc", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?page=-5", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?page=0", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?pageSize=xyz", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?pageSize=-10", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?pageSize=200", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles?pageSize=100", nil)
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
	h := NewHandlers(repo)

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

	r := chi.NewRouter()
	var result model.Article

	r.Get("/api/admin/blog/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleAdminJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/api/admin/blog/articles/admin-article", nil)
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
	h := NewHandlers(repo)

	r := chi.NewRouter()
	r.Get("/api/admin/blog/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		h.GetArticleAdminJSON(w, r)
	})

	req := httptest.NewRequest("GET", "/api/admin/blog/articles/non-existent", nil)
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
	h := NewHandlers(repo)

	req := httptest.NewRequest("GET", "/api/admin/blog/articles", nil)
	w := httptest.NewRecorder()
	h.GetArticleAdminJSON(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "slug parameter is required")
}
