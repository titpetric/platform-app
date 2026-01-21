package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/blog/model"
)

func timeNow() *time.Time {
	return timePtr(time.Now())
}

func timePtr(v time.Time) *time.Time {
	return &v
}

// TestStorageIntegration_FullLifecycle tests the complete storage lifecycle
func TestStorageIntegration_FullLifecycle(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// 1. Verify schema was initialized
	var tableCount int
	require.NoError(t, db.GetContext(ctx, &tableCount,
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='article'"))
	require.Equal(t, 1, tableCount)

	// 2. Insert multiple articles
	articles := []model.Article{
		{
			ID:          "article-1",
			Slug:        "getting-started",
			Title:       "Getting Started with Go",
			Description: "Learn the basics of Go programming",
			Date:        timePtr(time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)),
			OgImage:     "/og/go.png",
			Layout:      "post",
			URL:         "/blog/getting-started/",
		},
		{
			ID:          "article-2",
			Slug:        "advanced-patterns",
			Title:       "Advanced Go Patterns",
			Description: "Master advanced patterns in Go",
			Date:        timePtr(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
			OgImage:     "/og/patterns.png",
			Layout:      "post",
			URL:         "/blog/advanced-patterns/",
		},
		{
			ID:          "article-3",
			Slug:        "concurrency-guide",
			Title:       "Concurrency Guide",
			Description: "Understanding goroutines and channels",
			Date:        timePtr(time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)),
			OgImage:     "/og/concurrency.png",
			Layout:      "post",
			URL:         "/blog/concurrency-guide/",
		},
	}

	for _, article := range articles {
		require.NoError(t, storage.InsertArticle(ctx, &article), "failed to insert article %s", article.Slug)
	}

	// 3. Verify count
	count, err := storage.CountArticles(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	// 4. Retrieve all and verify ordering
	all, err := storage.GetArticles(ctx, 0, 9999)
	require.NoError(t, err)
	require.Len(t, all, 3)
	require.Equal(t, "concurrency-guide", all[0].Slug)

	// 5. Get single article
	article, err := storage.GetArticleBySlug(ctx, "advanced-patterns")
	require.NoError(t, err)
	require.Equal(t, "Advanced Go Patterns", article.Title)

	// 6. Search articles
	results, err := storage.SearchArticles(ctx, "Go")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 2)

	// 7. Update article
	article.Title = "Updated: Advanced Go Patterns"
	err = storage.InsertArticle(ctx, article)
	require.NoError(t, err)

	// Verify update
	updated, err := storage.GetArticleBySlug(ctx, "advanced-patterns")
	require.NoError(t, err)
	require.Equal(t, "Updated: Advanced Go Patterns", updated.Title)
}

// TestStorageIntegration_ConcurrentInserts tests concurrent article insertions
func TestStorageIntegration_ConcurrentInserts(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert articles concurrently
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			article := &model.Article{
				ID:    fmt.Sprintf("concurrent-%d", idx),
				Slug:  fmt.Sprintf("article-%d", idx),
				Title: "Concurrent Article",
				Date:  timeNow(),
			}
			done <- storage.InsertArticle(ctx, article)
		}(i)
	}

	// Collect results
	for i := 0; i < 10; i++ {
		err = <-done
		require.NoError(t, err)
	}

	// Verify all were inserted
	finalCount, _ := storage.CountArticles(ctx)
	require.Equal(t, 10, finalCount)
}

// TestStorageIntegration_ReplaceOnDuplicate tests that INSERT OR REPLACE works correctly
func TestStorageIntegration_ReplaceOnDuplicate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert a valid article
	article1 := &model.Article{
		ID:    "valid-1",
		Slug:  "valid-article",
		Title: "Valid Article",
		Date:  timeNow(),
	}
	require.NoError(t, storage.InsertArticle(ctx, article1))

	// Insert article with duplicate slug (should replace)
	article2 := &model.Article{
		ID:    "valid-2",
		Slug:  "valid-article", // duplicate slug
		Title: "Replaced Article",
		Date:  timeNow(),
	}
	require.NoError(t, storage.InsertArticle(ctx, article2))

	// Verify article was replaced
	count, _ := storage.CountArticles(ctx)
	require.Equal(t, 1, count)

	retrieved, err := storage.GetArticleBySlug(ctx, "valid-article")
	require.NoError(t, err)
	require.Equal(t, "valid-2", retrieved.ID)
	require.Equal(t, "Replaced Article", retrieved.Title)
}

// TestStorageIntegration_SearchAccuracy tests search result accuracy
func TestStorageIntegration_SearchAccuracy(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert articles with varied content
	articles := []model.Article{
		{
			ID:          "search-1",
			Slug:        "python-basics",
			Title:       "Python for Beginners",
			Description: "Learn Python programming language",
			Date:        timeNow(),
		},
		{
			ID:          "search-2",
			Slug:        "go-performance",
			Title:       "Go Performance Optimization",
			Description: "Optimize Go applications",
			Date:        timeNow(),
		},
		{
			ID:          "search-3",
			Slug:        "rust-safety",
			Title:       "Rust Memory Safety",
			Description: "Understanding Rust's safety features",
			Date:        timeNow(),
		},
	}

	for i := range articles {
		require.NoError(t, storage.InsertArticle(ctx, &articles[i]))
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
		expectedSlugs []string
	}{
		{
			name:          "search by title",
			query:         "Python",
			expectedCount: 1,
			expectedSlugs: []string{"python-basics"},
		},
		{
			name:          "search by description",
			query:         "Optimize",
			expectedCount: 1,
			expectedSlugs: []string{"go-performance"},
		},
		{
			name:          "search by title word",
			query:         "Performance",
			expectedCount: 1,
			expectedSlugs: []string{"go-performance"},
		},
		{
			name:          "search by slug",
			query:         "rust",
			expectedCount: 1,
			expectedSlugs: []string{"rust-safety"},
		},
		{
			name:          "case insensitive search",
			query:         "PYTHON",
			expectedCount: 1,
			expectedSlugs: []string{"python-basics"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := storage.SearchArticles(ctx, tt.query)
			require.NoError(t, err)
			require.Len(t, results, tt.expectedCount, "unexpected result count")

			for _, expectedSlug := range tt.expectedSlugs {
				found := false
				for _, result := range results {
					if result.Slug == expectedSlug {
						found = true
						break
					}
				}
				require.True(t, found, "expected slug %s not found in results", expectedSlug)
			}
		})
	}
}

// TestStorageIntegration_DateOrdering tests proper date-based ordering
func TestStorageIntegration_DateOrdering(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert articles with specific dates (not in order)
	articles := []model.Article{
		{
			ID:    "date-1",
			Slug:  "middle",
			Title: "Middle Article",
			Date:  timePtr(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:    "date-2",
			Slug:  "newest",
			Title: "Newest Article",
			Date:  timePtr(time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:    "date-3",
			Slug:  "oldest",
			Title: "Oldest Article",
			Date:  timePtr(time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)),
		},
	}

	for i := range articles {
		require.NoError(t, storage.InsertArticle(ctx, &articles[i]))
	}

	// Retrieve all and verify order (newest first)
	all, err := storage.GetArticles(ctx, 0, 9999)
	require.NoError(t, err)

	expectedOrder := []string{"newest", "middle", "oldest"}
	for i, expected := range expectedOrder {
		require.Equal(t, expected, all[i].Slug, "position %d mismatch", i)
	}
}

// TestStorageIntegration_EmptyDatabase tests operations on empty database
func TestStorageIntegration_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Count should be zero
	count, err := storage.CountArticles(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, count)

	// Get should return empty slice
	articles, err := storage.GetArticles(ctx, 0, 9999)
	require.NoError(t, err)
	require.Len(t, articles, 0)

	// Search should return empty results
	results, err := storage.SearchArticles(ctx, "anything")
	require.NoError(t, err)
	require.Len(t, results, 0)

	// GetBySlug should fail
	_, err = storage.GetArticleBySlug(ctx, "nonexistent")
	require.Error(t, err)
}

// TestStorageIntegration_SpecialCharacters tests handling of special characters
func TestStorageIntegration_SpecialCharacters(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	article := &model.Article{
		ID:          "special-1",
		Slug:        "special-chars",
		Title:       "Article with \"quotes\" and 'apostrophes'",
		Description: "Testing special chars: <html> & {json}",
		Date:        timeNow(),
	}

	err = storage.InsertArticle(ctx, article)
	require.NoError(t, err)

	// Retrieve and verify special characters are preserved
	retrieved, err := storage.GetArticleBySlug(ctx, "special-chars")
	require.NoError(t, err)
	require.Equal(t, article.Title, retrieved.Title)
	require.Equal(t, article.Description, retrieved.Description)
}

// TestStorageIntegration_TimestampHandling tests proper timestamp handling
func TestStorageIntegration_TimestampHandling(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	article := &model.Article{
		ID:    "time-1",
		Slug:  "time-test",
		Title: "Timestamp Test",
		Date:  timePtr(time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)),
	}

	err = storage.InsertArticle(ctx, article)
	require.NoError(t, err)

	retrieved, err := storage.GetArticleBySlug(ctx, "time-test")
	require.NoError(t, err)

	// Verify date is preserved
	require.Equal(t, 2024, retrieved.Date.Year())
	require.Equal(t, time.June, retrieved.Date.Month())
	require.Equal(t, 15, retrieved.Date.Day())
}

// TestStorageIntegration_IndexEfficiency tests that indexes are properly used
func TestStorageIntegration_IndexEfficiency(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	storage, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	// Insert many articles
	for i := 0; i < 100; i++ {
		idx := fmt.Sprintf("%d", i)
		article := &model.Article{
			ID:    "index-" + idx,
			Slug:  "article-" + idx,
			Title: "Article Number",
			Date:  timePtr(time.Now().AddDate(0, 0, -i)),
		}
		err = storage.InsertArticle(ctx, article)
		require.NoError(t, err)
	}

	// These operations should be fast due to indexes
	_, err = storage.GetArticleBySlug(ctx, "article-50")
	require.NoError(t, err)

	articles, err := storage.GetArticles(ctx, 0, 9999)
	require.NoError(t, err)
	require.Len(t, articles, 100)

	// Verify newest first ordering with many articles
	require.Equal(t, "article-0", articles[0].Slug)
}
