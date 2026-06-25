package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// setupTestDB creates a temporary SQLite database for testing.
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

// setupHandlers returns a Handlers with repo + tmp GitFS for CRUD tests.
// Views are nil; we test JSON endpoints only.
func setupHandlers(t *testing.T) (*Handlers, *storage.Storage, *storage.GitFS) {
	db := setupTestDB(t)
	repo, err := storage.NewStorage(t.Context(), db)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	gfs, err := storage.NewGitFS(tmpDir)
	require.NoError(t, err)

	h := &Handlers{repository: repo, contentFS: gfs}
	return h, repo, gfs
}

func chiRouter(method, pattern string, handler http.HandlerFunc) *chi.Mux {
	r := chi.NewRouter()
	r.Method(method, pattern, handler)
	return r
}

func TestCreateArticleJSON_Success(t *testing.T) {
	h, repo, gfs := setupHandlers(t)

	req := ArticleRequest{
		Slug:        "hello-world",
		Title:       "Hello, World",
		Description: "Greetings",
		Content:     "# Hello\n\nBody text.",
		Date:        "2024-06-01",
		Draft:       false,
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)

	require.Equal(t, http.StatusCreated, w.Code, "body: %s", w.Body.String())

	var got model.Article
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "hello-world", got.Slug)
	assert.Equal(t, "hello-world.md", got.Filename)

	// DB row exists
	stored, err := repo.GetArticleBySlug(t.Context(), "hello-world")
	require.NoError(t, err)
	assert.Equal(t, "Hello, World", stored.Title)

	// File exists in GitFS
	content, err := gfs.ReadFile("hello-world.md")
	require.NoError(t, err)
	assert.Contains(t, string(got.Filename), "hello-world.md")
	assert.Contains(t, string(content), `title: "Hello, World"`)
}

func TestCreateArticleJSON_InvalidPayload(t *testing.T) {
	h, _, _ := setupHandlers(t)

	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateArticleJSON_ValidationError(t *testing.T) {
	h, _, _ := setupHandlers(t)

	req := ArticleRequest{Slug: "", Title: "x", Content: "y"}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "slug is required")
}

func TestUpdateArticleJSON_Success(t *testing.T) {
	h, repo, gfs := setupHandlers(t)

	// Seed an article on disk + in DB
	createReq := ArticleRequest{
		Slug: "draft-post", Title: "Original", Content: "old", Draft: true,
	}
	body, _ := json.Marshal(createReq)
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	// Update
	updateReq := ArticleRequest{
		Slug: "draft-post", Title: "Updated", Description: "new desc", Content: "new content", Draft: false,
	}
	updateBody, _ := json.Marshal(updateReq)

	router := chiRouter(http.MethodPut, "/api/admin/blog/articles/{slug}", h.UpdateArticleJSON)
	r2 := httptest.NewRequest(http.MethodPut, "/api/admin/blog/articles/draft-post", bytes.NewReader(updateBody))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)

	require.Equal(t, http.StatusOK, w2.Code, "body: %s", w2.Body.String())

	stored, err := repo.GetArticleBySlug(t.Context(), "draft-post")
	require.NoError(t, err)
	assert.Equal(t, "Updated", stored.Title)
	assert.Equal(t, "new desc", stored.Description)
	assert.Equal(t, int64(0), stored.Draft)

	// File reflects the new content
	content, err := gfs.ReadFile("draft-post.md")
	require.NoError(t, err)
	assert.Contains(t, string(content), "new content")
	assert.Contains(t, string(content), `title: "Updated"`)
	assert.NotContains(t, string(content), "draft: true")
}

func TestUpdateArticleJSON_ValidationFails(t *testing.T) {
	h, _, _ := setupHandlers(t)

	// Seed
	createReq := ArticleRequest{Slug: "post-a", Title: "T", Content: "C"}
	body, _ := json.Marshal(createReq)
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	// Update with empty content should fail validation now
	bad := ArticleRequest{Slug: "post-a", Title: "", Content: ""}
	badBody, _ := json.Marshal(bad)
	router := chiRouter(http.MethodPut, "/api/admin/blog/articles/{slug}", h.UpdateArticleJSON)
	r2 := httptest.NewRequest(http.MethodPut, "/api/admin/blog/articles/post-a", bytes.NewReader(badBody))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)

	require.Equal(t, http.StatusBadRequest, w2.Code, "body: %s", w2.Body.String())
}

func TestUpdateArticleJSON_NotFound(t *testing.T) {
	h, _, _ := setupHandlers(t)

	updateReq := ArticleRequest{Slug: "missing", Title: "T", Content: "C"}
	body, _ := json.Marshal(updateReq)

	router := chiRouter(http.MethodPut, "/api/admin/blog/articles/{slug}", h.UpdateArticleJSON)
	r := httptest.NewRequest(http.MethodPut, "/api/admin/blog/articles/missing", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteArticleJSON_Success(t *testing.T) {
	h, repo, gfs := setupHandlers(t)

	createReq := ArticleRequest{Slug: "to-delete", Title: "T", Content: "C"}
	body, _ := json.Marshal(createReq)
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	router := chiRouter(http.MethodDelete, "/api/admin/blog/articles/{slug}", h.DeleteArticleJSON)
	r2 := httptest.NewRequest(http.MethodDelete, "/api/admin/blog/articles/to-delete", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)

	require.Equal(t, http.StatusNoContent, w2.Code, "body: %s", w2.Body.String())

	_, err := repo.GetArticleBySlug(t.Context(), "to-delete")
	require.Error(t, err, "expected article to be gone from DB")

	_, err = gfs.Stat("to-delete.md")
	assert.Error(t, err, "expected file to be removed")
}

func TestDeleteArticleJSON_NotFound(t *testing.T) {
	h, _, _ := setupHandlers(t)

	router := chiRouter(http.MethodDelete, "/api/admin/blog/articles/{slug}", h.DeleteArticleJSON)
	r := httptest.NewRequest(http.MethodDelete, "/api/admin/blog/articles/missing-thing", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteArticleJSON_InvalidSlug(t *testing.T) {
	h, _, _ := setupHandlers(t)

	router := chiRouter(http.MethodDelete, "/api/admin/blog/articles/{slug}", h.DeleteArticleJSON)
	r := httptest.NewRequest(http.MethodDelete, "/api/admin/blog/articles/Invalid_Slug", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublishArticleJSON_UpdatesFileAndDB(t *testing.T) {
	h, repo, gfs := setupHandlers(t)

	// Create a draft article
	createReq := ArticleRequest{Slug: "to-publish", Title: "T", Content: "C", Draft: true}
	body, _ := json.Marshal(createReq)
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	// Verify draft marker is in file
	before, err := gfs.ReadFile("to-publish.md")
	require.NoError(t, err)
	require.Contains(t, string(before), "draft: true")

	// Publish
	router := chiRouter(http.MethodPost, "/api/admin/blog/articles/{slug}/publish", h.PublishArticleJSON)
	r2 := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles/to-publish/publish", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)

	require.Equal(t, http.StatusOK, w2.Code, "body: %s", w2.Body.String())

	// DB updated
	stored, err := repo.GetArticleBySlug(t.Context(), "to-publish")
	require.NoError(t, err)
	assert.Equal(t, int64(0), stored.Draft)

	// File frontmatter no longer contains "draft: true"
	after, err := gfs.ReadFile("to-publish.md")
	require.NoError(t, err)
	assert.NotContains(t, string(after), "draft: true")
}

func TestCheckSlugJSON_AvailableAndInvalid(t *testing.T) {
	h, _, _ := setupHandlers(t)

	router := chiRouter(http.MethodGet, "/api/admin/blog/articles/{slug}/check", h.CheckSlugJSON)

	// Available slug
	r := httptest.NewRequest(http.MethodGet, "/api/admin/blog/articles/free-slug/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"available":true`)

	// Invalid slug
	r2 := httptest.NewRequest(http.MethodGet, "/api/admin/blog/articles/Bad_Slug/check", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, r2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

func TestSaveSettingsJSON_Validation(t *testing.T) {
	h, _, _ := setupHandlers(t)

	tests := []struct {
		name    string
		body    map[string]any
		wantErr string
	}{
		{
			name:    "posts_per_page negative",
			body:    map[string]any{"posts_per_page": -1},
			wantErr: "posts_per_page",
		},
		{
			name:    "posts_per_page too large",
			body:    map[string]any{"posts_per_page": 5000},
			wantErr: "posts_per_page",
		},
		{
			name:    "invalid url",
			body:    map[string]any{"meta_url": "not-a-url"},
			wantErr: "meta_url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/settings", bytes.NewReader(body))
			w := httptest.NewRecorder()
			h.SaveSettingsJSON(w, r)
			require.Equal(t, http.StatusBadRequest, w.Code, "body: %s", w.Body.String())
			assert.Contains(t, w.Body.String(), tt.wantErr)
		})
	}
}

func TestSaveSettingsJSON_Success(t *testing.T) {
	h, repo, _ := setupHandlers(t)

	body, _ := json.Marshal(map[string]any{
		"meta_lang":      "en",
		"meta_url":       "https://example.com",
		"posts_per_page": 25,
	})
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/settings", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.SaveSettingsJSON(w, r)

	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())

	settings, err := repo.GetGlobalSettings(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "en", settings.MetaLang)
	assert.Equal(t, "https://example.com", settings.MetaURL)
	assert.Equal(t, int64(25), settings.PostsPerPage)
}

func TestRemoveDraftFromFrontmatter(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
		hasDoc bool
	}{
		{
			name:  "removes draft true line",
			input: "---\ntitle: \"X\"\ndraft: true\n---\n\nbody",
			want:  "---\ntitle: \"X\"\n---\n\nbody",
		},
		{
			name:  "removes draft anywhere",
			input: "---\ndraft: true\ntitle: \"X\"\n---\nbody",
			want:  "---\ntitle: \"X\"\n---\nbody",
		},
		{
			name:  "no frontmatter is unchanged",
			input: "no frontmatter here",
			want:  "no frontmatter here",
		},
		{
			name:  "no draft line is unchanged",
			input: "---\ntitle: \"X\"\n---\nbody",
			want:  "---\ntitle: \"X\"\n---\nbody",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(removeDraftFromFrontmatter([]byte(tt.input)))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEscapeYAML_NewlineAndQuotes(t *testing.T) {
	// Newlines must not leak into a single-line YAML scalar
	got := escapeYAML("hello\nworld")
	assert.Equal(t, `hello\nworld`, got)

	got = escapeYAML("with \"quotes\"")
	assert.Equal(t, `with \"quotes\"`, got)

	got = escapeYAML("carriage\rreturn")
	assert.Equal(t, `carriage\rreturn`, got)
}

func TestIsValidSlugAdmin(t *testing.T) {
	tests := []struct {
		slug string
		ok   bool
	}{
		{"hello", true},
		{"hello-world", true},
		{"hello-world-123", true},
		{"", false},
		{"Hello", false},
		{"hello_world", false},
		{"hello world", false},
		{"-hello", false},
		{"hello-", false},
		{"hello--world", false},
		{"/etc/passwd", false},
		{"../escape", false},
	}
	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			assert.Equal(t, tt.ok, isValidSlugAdmin(tt.slug))
		})
	}
}

func TestCreateArticleJSON_DuplicateSlugConflict(t *testing.T) {
	h, _, _ := setupHandlers(t)

	req := ArticleRequest{Slug: "same-slug", Title: "First", Content: "C"}
	body, _ := json.Marshal(req)

	// First create succeeds
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateArticleJSON(w, r)
	require.Equal(t, http.StatusCreated, w.Code, "body: %s", w.Body.String())

	// Second create with same slug must return 409 and NOT silently overwrite
	dup := ArticleRequest{Slug: "same-slug", Title: "Second", Content: "X"}
	body2, _ := json.Marshal(dup)
	r2 := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body2))
	w2 := httptest.NewRecorder()
	h.CreateArticleJSON(w2, r2)
	assert.Equal(t, http.StatusConflict, w2.Code, "body: %s", w2.Body.String())
}

func TestEditArticleHTML_InvalidSlug(t *testing.T) {
	h, _, _ := setupHandlers(t)

	router := chiRouter(http.MethodGet, "/admin/blog/articles/{slug}", h.EditArticleHTML)
	r := httptest.NewRequest(http.MethodGet, "/admin/blog/articles/Bad_Slug", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPublishArticleJSON_NotFound(t *testing.T) {
	h, _, _ := setupHandlers(t)

	router := chiRouter(http.MethodPost, "/api/admin/blog/articles/{slug}/publish", h.PublishArticleJSON)
	r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles/missing/publish", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestCreateAndListInteroperability verifies that creating articles works at
// volume and that they show up in subsequent listing calls.
func TestCreateAndListInteroperability(t *testing.T) {
	h, repo, _ := setupHandlers(t)

	for i := 0; i < 5; i++ {
		req := ArticleRequest{
			Slug:    fmt.Sprintf("post-%d", i),
			Title:   fmt.Sprintf("Post %d", i),
			Content: "body",
			Draft:   i%2 == 0,
			Date:    time.Now().Format("2006-01-02"),
		}
		body, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "/api/admin/blog/articles", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.CreateArticleJSON(w, r)
		require.Equal(t, http.StatusCreated, w.Code, "iter %d body: %s", i, w.Body.String())
	}

	drafts, err := repo.GetDraftArticles(t.Context(), 0, 100)
	require.NoError(t, err)
	assert.Equal(t, 3, len(drafts), "expected 3 drafts (i=0,2,4)")
}
