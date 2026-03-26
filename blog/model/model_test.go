package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/platform-app/blog/model"
)

func TestArticle_IsDraft(t *testing.T) {
	t.Run("draft article", func(t *testing.T) {
		article := &model.Article{Draft: 1}
		assert.True(t, article.IsDraft())
	})

	t.Run("non-draft article", func(t *testing.T) {
		article := &model.Article{Draft: 0}
		assert.False(t, article.IsDraft())
	})
}

func TestArticle_IsScheduled(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	t.Run("scheduled article (future date)", func(t *testing.T) {
		article := &model.Article{Draft: 0, Date: &futureDate}
		assert.True(t, article.IsScheduled())
	})

	t.Run("draft article with future date is not scheduled", func(t *testing.T) {
		article := &model.Article{Draft: 1, Date: &futureDate}
		assert.False(t, article.IsScheduled())
	})

	t.Run("past date is not scheduled", func(t *testing.T) {
		article := &model.Article{Draft: 0, Date: &pastDate}
		assert.False(t, article.IsScheduled())
	})

	t.Run("nil date is not scheduled", func(t *testing.T) {
		article := &model.Article{Draft: 0, Date: nil}
		assert.False(t, article.IsScheduled())
	})
}

func TestArticle_IsPublished(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	t.Run("published article (past date)", func(t *testing.T) {
		article := &model.Article{Draft: 0, Date: &pastDate}
		assert.True(t, article.IsPublished())
	})

	t.Run("published article (nil date)", func(t *testing.T) {
		article := &model.Article{Draft: 0, Date: nil}
		assert.True(t, article.IsPublished())
	})

	t.Run("draft article is not published", func(t *testing.T) {
		article := &model.Article{Draft: 1, Date: &pastDate}
		assert.False(t, article.IsPublished())
	})

	t.Run("scheduled article is not published", func(t *testing.T) {
		article := &model.Article{Draft: 0, Date: &futureDate}
		assert.False(t, article.IsPublished())
	})
}

func TestArticle_Status(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)
	pastDate := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name     string
		article  *model.Article
		expected model.ArticleStatus
	}{
		{
			name:     "draft status",
			article:  &model.Article{Draft: 1},
			expected: model.StatusDraft,
		},
		{
			name:     "scheduled status",
			article:  &model.Article{Draft: 0, Date: &futureDate},
			expected: model.StatusScheduled,
		},
		{
			name:     "published status",
			article:  &model.Article{Draft: 0, Date: &pastDate},
			expected: model.StatusPublished,
		},
		{
			name:     "published status (nil date)",
			article:  &model.Article{Draft: 0, Date: nil},
			expected: model.StatusPublished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.article.Status())
		})
	}
}
