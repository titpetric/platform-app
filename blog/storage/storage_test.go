package storage

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/blog/model"
	"github.com/titpetric/platform-app/blog/schema"
)

// setupTestDB creates a temporary SQLite database for testing with automatic cleanup
func setupTestDB(t *testing.T) *sqlx.DB {
	// Use temporary file for each test to ensure isolation
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// SQLite only supports one concurrent connection
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	// Initialize schema
	if err := Migrate(context.Background(), db, schema.Migrations); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Logf("warning: failed to close test database: %v", err)
		}
	})

	return db
}

// cleanupTestDB is a no-op since setupTestDB handles cleanup via t.Cleanup()
func cleanupTestDB(t *testing.T, db *sqlx.DB) {
	// Cleanup is handled by t.Cleanup() in setupTestDB
}

// TestInitSchema tests schema initialization
func TestInitSchema(t *testing.T) {
	db := setupTestDB(t)

	_, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Verify tables exist (NewStorage runs migrations automatically)
	var tableCount int
	require.NoError(t, db.GetContext(ctx, &tableCount,
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='article'"))

	require.Equal(t, 1, tableCount, "expected article table to exist")
}

// TestInsertArticle tests inserting an article
func TestInsertArticle(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	article := &model.Article{
		ID:          "test-1",
		Slug:        "test-article",
		Title:       "Test Article",
		Description: "This is a test article",
		OgImage:     "/images/test.png",
		Layout:      "post",
		Source:      "",
		URL:         "/blog/test-article/",
	}

	err = storage.InsertArticle(ctx, article)
	if err != nil {
		t.Fatalf("InsertArticle() failed: %v", err)
	}

	// Verify article was inserted
	retrieved, err := storage.GetArticleBySlug(ctx, "test-article")
	if err != nil {
		t.Fatalf("GetArticleBySlug() failed: %v", err)
	}

	if retrieved.Title != "Test Article" {
		t.Errorf("expected title 'Test Article', got '%s'", retrieved.Title)
	}

	if retrieved.Slug != "test-article" {
		t.Errorf("expected slug 'test-article', got '%s'", retrieved.Slug)
	}
}

// TestGetArticleBySlug tests retrieving an article by slug
func TestGetArticleBySlug(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert test article
	article := &model.Article{
		ID:          "test-2",
		Slug:        "my-article",
		Title:       "My Article",
		Description: "Description",
		OgImage:     "/og.png",
		Layout:      "post",
		URL:         "/blog/my-article/",
	}

	storage.InsertArticle(ctx, article)

	// Retrieve and verify
	retrieved, err := storage.GetArticleBySlug(ctx, "my-article")
	if err != nil {
		t.Fatalf("GetArticleBySlug() failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected article, got nil")
	}

	if retrieved.ID != "test-2" {
		t.Errorf("expected ID 'test-2', got '%s'", retrieved.ID)
	}
}

// TestGetArticleBySlugNotFound tests retrieving non-existent article
func TestGetArticleBySlugNotFound(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	_, err = storage.GetArticleBySlug(ctx, "non-existent")
	require.Error(t, err, "expected error for non-existent article")
}

// TestGetArticles tests retrieving all articles
func TestGetArticles(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert test articles
	articles := []model.Article{
		{
			ID:    "test-3",
			Slug:  "first",
			Title: "First Article",
		},
		{
			ID:    "test-4",
			Slug:  "second",
			Title: "Second Article",
		},
		{
			ID:    "test-5",
			Slug:  "third",
			Title: "Third Article",
		},
	}

	for _, article := range articles {
		storage.InsertArticle(ctx, &article)
	}

	// Retrieve all
	retrieved, err := storage.GetArticles(ctx, 0, 10)
	if err != nil {
		t.Fatalf("GetAllArticles() failed: %v", err)
	}

	if len(retrieved) != 3 {
		t.Errorf("expected 3 articles, got %d", len(retrieved))
	}

	// Verify sorted by date DESC (newest first)
	if retrieved[0].Slug != "third" {
		t.Errorf("expected first article 'third', got '%s'", retrieved[0].Slug)
	}

	if retrieved[2].Slug != "first" {
		t.Errorf("expected last article 'first', got '%s'", retrieved[2].Slug)
	}
}

// TestSearchArticles tests searching articles
func TestSearchArticles(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert test articles
	articles := []model.Article{
		{
			ID:          "test-6",
			Slug:        "go-programming",
			Title:       "Learning Go Programming",
			Description: "A comprehensive guide to Go",
		},
		{
			ID:          "test-7",
			Slug:        "rust-basics",
			Title:       "Rust Basics",
			Description: "Introduction to Rust programming",
		},
		{
			ID:          "test-8",
			Slug:        "go-frameworks",
			Title:       "Go Web Frameworks",
			Description: "Popular web frameworks in Go",
		},
	}

	for i := range articles {
		storage.InsertArticle(ctx, &articles[i])
	}

	// Search for "Go"
	results, err := storage.SearchArticles(ctx, "Go")
	if err != nil {
		t.Fatalf("SearchArticles() failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results for 'Go', got %d", len(results))
	}

	// Search for "Rust"
	results, err = storage.SearchArticles(ctx, "Rust")
	if err != nil {
		t.Fatalf("SearchArticles() failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result for 'Rust', got %d", len(results))
	}

	if results[0].Slug != "rust-basics" {
		t.Errorf("expected 'rust-basics', got '%s'", results[0].Slug)
	}
}

// TestCountArticles tests counting articles
func TestCountArticles(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Should be zero initially
	count, err := storage.CountArticles(ctx)
	if err != nil {
		t.Fatalf("CountArticles() failed: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 articles initially, got %d", count)
	}

	// Insert some articles
	for i := 0; i < 5; i++ {
		article := &model.Article{
			ID:    "test-" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article",
		}
		storage.InsertArticle(ctx, article)
	}

	// Count should be 5
	count, err = storage.CountArticles(ctx)
	if err != nil {
		t.Fatalf("CountArticles() failed: %v", err)
	}

	if count != 5 {
		t.Errorf("expected 5 articles, got %d", count)
	}
}

// TestInsertArticleUpdate tests updating an article via INSERT OR REPLACE
func TestInsertArticleUpdate(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert initial article
	article := &model.Article{
		ID:    "test-update",
		Slug:  "update-test",
		Title: "Original Title",
	}

	storage.InsertArticle(ctx, article)

	// Update the article
	article.Title = "Updated Title"

	err = storage.InsertArticle(ctx, article)
	if err != nil {
		t.Fatalf("InsertArticle() update failed: %v", err)
	}

	// Retrieve and verify
	retrieved, err := storage.GetArticleBySlug(ctx, "update-test")
	if err != nil {
		t.Fatalf("GetArticleBySlug() failed: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("expected updated title 'Updated Title', got '%s'", retrieved.Title)
	}
}

// TestSchemaConstraints tests that INSERT OR REPLACE replaces on duplicate slug
func TestSchemaConstraints(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert first article
	article1 := &model.Article{
		ID:    "test-constraint-1",
		Slug:  "unique-slug",
		Title: "Article 1",
	}

	err = storage.InsertArticle(ctx, article1)
	if err != nil {
		t.Fatalf("first insert failed: %v", err)
	}

	// Insert article with duplicate slug (should replace)
	article2 := &model.Article{
		ID:    "test-constraint-2",
		Slug:  "unique-slug",
		Title: "Article 2",
	}

	err = storage.InsertArticle(ctx, article2)
	if err != nil {
		t.Fatalf("second insert should succeed: %v", err)
	}

	// Verify article was replaced
	retrieved, err := storage.GetArticleBySlug(ctx, "unique-slug")
	if err != nil {
		t.Fatalf("failed to retrieve article: %v", err)
	}

	if retrieved.ID != "test-constraint-2" {
		t.Errorf("expected ID test-constraint-2, got %s", retrieved.ID)
	}
}

// TestSearchArticlesEmpty tests searching when no articles match
func TestSearchArticlesEmpty(t *testing.T) {
	db := setupTestDB(t)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Search with no articles
	results, err := storage.SearchArticles(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("SearchArticles() failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

// BenchmarkInsertArticle benchmarks article insertion
func BenchmarkInsertArticle(b *testing.B) {
	db, _ := sqlx.Open("sqlite", ":memory:")
	defer db.Close()

	storage, err := NewStorage(b.Context(), db)
	if err != nil {
		b.Fatalf("failed to create storage: %v", err)
	}
	ctx := b.Context()

	article := &model.Article{
		ID:          "bench",
		Slug:        "benchmark",
		Title:       "Benchmark Article",
		Description: "This is for benchmarking",
		URL:         "/blog/benchmark/",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		article.ID = "bench-" + string(rune(i))
		storage.InsertArticle(ctx, article)
	}
}

// BenchmarkGetArticleBySlug benchmarks article retrieval by slug
func BenchmarkGetArticleBySlug(b *testing.B) {
	db, _ := sqlx.Open("sqlite", ":memory:")
	defer db.Close()

	storage, err := NewStorage(b.Context(), db)
	if err != nil {
		b.Fatalf("failed to create storage: %v", err)
	}
	ctx := b.Context()

	// Insert test articles
	for i := 0; i < 100; i++ {
		article := &model.Article{
			ID:    "bench-get-" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article",
		}
		storage.InsertArticle(ctx, article)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.GetArticleBySlug(ctx, "article-50")
	}
}

// BenchmarkGetArticles benchmarks retrieving 20% of articles
func BenchmarkGetAllArticles(b *testing.B) {
	db, _ := sqlx.Open("sqlite", ":memory:")
	defer db.Close()

	storage, err := NewStorage(b.Context(), db)
	if err != nil {
		b.Fatalf("failed to create storage: %v", err)
	}
	ctx := b.Context()

	// Insert test articles
	for i := 0; i < 100; i++ {
		article := &model.Article{
			ID:    "bench-all-" + string(rune(i)),
			Slug:  "article-" + string(rune(i)),
			Title: "Article",
		}
		storage.InsertArticle(ctx, article)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.GetArticles(ctx, 0, 20)
	}
}

// BenchmarkSearchArticles benchmarks searching articles
func BenchmarkSearchArticles(b *testing.B) {
	db, _ := sqlx.Open("sqlite", ":memory:")
	defer db.Close()

	storage, err := NewStorage(b.Context(), db)
	if err != nil {
		b.Fatalf("failed to create storage: %v", err)
	}
	ctx := b.Context()

	// Insert test articles
	for i := 0; i < 100; i++ {
		article := &model.Article{
			ID:          "bench-search-" + string(rune(i)),
			Slug:        "article-" + string(rune(i)),
			Title:       "Test Article Number " + string(rune(i)),
			Description: "This is a test article about testing",
		}
		storage.InsertArticle(ctx, article)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.SearchArticles(ctx, "test")
	}
}
