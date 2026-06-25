package storage

import (
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform-app/blog/model"
)

func seedVisibility(t *testing.T, repo *Storage) {
	t.Helper()
	ctx := t.Context()
	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)

	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "v-pub", Slug: "published", Title: "Pub", Date: &past, Draft: 0,
	}))
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "v-draft", Slug: "draft", Title: "Draft", Date: &past, Draft: 1,
	}))
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "v-sched", Slug: "scheduled", Title: "Sched", Date: &future, Draft: 0,
	}))
}

func TestGetPublishedArticleBySlug_ExcludesDraftAndScheduled(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	seedVisibility(t, repo)

	got, err := repo.GetPublishedArticleBySlug(t.Context(), "published")
	require.NoError(t, err)
	assert.Equal(t, "published", got.Slug)

	_, err = repo.GetPublishedArticleBySlug(t.Context(), "draft")
	assert.Error(t, err, "draft should not be returned")

	_, err = repo.GetPublishedArticleBySlug(t.Context(), "scheduled")
	assert.Error(t, err, "scheduled should not be returned")
}

func TestSearchPublishedArticles_OnlyPublished(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	seedVisibility(t, repo)

	// All three articles match "v-" prefix conceptually via title or slug,
	// but only "published" should come back from a published-only search.
	results, err := repo.SearchPublishedArticles(t.Context(), "pub")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "published", results[0].Slug)
}

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"plain", "plain"},
		{"50%", `50\%`},
		{"with_under", `with\_under`},
		{`back\slash`, `back\\slash`},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, escapeLike(tt.in))
		})
	}
}

func TestSearchArticles_LiteralWildcardMatch(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewStorage(t.Context(), db)
	require.NoError(t, err)
	ctx := t.Context()

	now := time.Now()
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "lit-1", Slug: "fifty-percent", Title: "50% off", Date: &now,
	}))
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "lit-2", Slug: "anything", Title: "no percent here", Date: &now,
	}))

	// "%" should now be a literal character; only the one with "50%" should match.
	res, err := repo.SearchArticles(ctx, "50%")
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "fifty-percent", res[0].Slug)

	// Underscore is also literal
	require.NoError(t, repo.InsertArticle(ctx, &model.Article{
		ID: "lit-3", Slug: "with-underscore", Title: "hello_world", Date: &now,
	}))
	res, err = repo.SearchArticles(ctx, "hello_world")
	require.NoError(t, err)
	require.Len(t, res, 1)
}
