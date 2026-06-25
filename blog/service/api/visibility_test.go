package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	chi "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/storage"
)

// seedVisibilityArticles inserts one published, one draft, and one scheduled article.
func seedVisibilityArticles(t *testing.T, repo *storage.Storage) {
	t.Helper()
	ctx := t.Context()
	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "vis-pub", Slug: "published-post", Title: "Public", Date: &past, Draft: 0,
	}))
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "vis-draft", Slug: "draft-post", Title: "Secret Draft", Date: &past, Draft: 1,
	}))
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "vis-sched", Slug: "scheduled-post", Title: "Future Post", Date: &future, Draft: 0,
	}))
}

func TestListArticlesJSON_ExcludesDraftsAndScheduled(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	seedVisibilityArticles(t, repo)

	h := NewHandlers(repo)
	r := httptest.NewRequest(http.MethodGet, "/api/blog/articles", nil)
	w := httptest.NewRecorder()
	h.ListArticlesJSON(w, r)

	require.Equal(t, http.StatusOK, w.Code)

	var list model.ArticleList
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	require.Equal(t, 1, list.Total, "only published article should be listed")
	require.Equal(t, "published-post", list.Articles[0].Slug)
}

func TestGetArticleJSON_DraftReturns404(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	seedVisibilityArticles(t, repo)

	h := NewHandlers(repo)
	router := chi.NewRouter()
	router.Get("/api/blog/articles/{slug}", h.GetArticleJSON)

	for _, slug := range []string{"draft-post", "scheduled-post"} {
		r := httptest.NewRequest(http.MethodGet, "/api/blog/articles/"+slug, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		require.Equal(t, http.StatusNotFound, w.Code, "slug %s should be 404", slug)
	}

	// Published article is reachable
	r := httptest.NewRequest(http.MethodGet, "/api/blog/articles/published-post", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetArticleJSON_InvalidSlug(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)

	h := NewHandlers(repo)
	router := chi.NewRouter()
	router.Get("/api/blog/articles/{slug}", h.GetArticleJSON)

	r := httptest.NewRequest(http.MethodGet, "/api/blog/articles/Bad_Slug", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSearchArticlesJSON_ExcludesDrafts(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)
	seedVisibilityArticles(t, repo)

	h := NewHandlers(repo)
	r := httptest.NewRequest(http.MethodGet, "/api/blog/search?q=post", nil)
	w := httptest.NewRecorder()
	h.SearchArticlesJSON(w, r)

	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, float64(1), result["total"], "search should only return published")
}

func TestSearchArticlesJSON_QueryLengthValidation(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)

	h := NewHandlers(repo)
	long := make([]byte, maxSearchQueryLength+1)
	for i := range long {
		long[i] = 'a'
	}
	r := httptest.NewRequest(http.MethodGet, "/api/blog/search?q="+string(long), nil)
	w := httptest.NewRecorder()
	h.SearchArticlesJSON(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "exceeds maximum length")
}

func TestSearchArticlesJSON_LikeWildcardsEscaped(t *testing.T) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)

	// Insert two articles where one would match a wildcard literal and one would not.
	past := time.Now().Add(-time.Hour)
	require.NoError(t, repo.InsertArticle(t.Context(), &model.Article{
		ID: "lk-1", Slug: "post-50-pct", Title: "Discount 50%", Date: &past, Draft: 0,
	}))
	require.NoError(t, repo.InsertArticle(t.Context(), &model.Article{
		ID: "lk-2", Slug: "post-plain", Title: "Plain Post", Date: &past, Draft: 0,
	}))

	h := NewHandlers(repo)

	// Searching literally for "50%" must only match the article that contains "50%"
	r := httptest.NewRequest(http.MethodGet, "/api/blog/search?q=50%25", nil) // %25 == "%"
	w := httptest.NewRecorder()
	h.SearchArticlesJSON(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	var result map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	require.Equal(t, float64(1), result["total"], "wildcard '%%' must be escaped, only one article should match")
}

func TestIsValidSlug(t *testing.T) {
	tests := []struct {
		slug string
		ok   bool
	}{
		{"hello", true},
		{"hello-world", true},
		{"", false},
		{"Hello", false},
		{"hello_x", false},
		{"hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			assert.Equal(t, tt.ok, isValidSlug(tt.slug))
		})
	}
}
