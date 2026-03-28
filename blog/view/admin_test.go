package view

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/platform-app/blog/model"
)

func TestAdminDashboardData_Map(t *testing.T) {
	data := NewAdminDashboardData(5, 3, 10, nil, nil, nil)

	assert.Equal(t, "Blog Dashboard", data.Title)
	assert.Equal(t, 5, data.DraftsCount)
	assert.Equal(t, 3, data.ScheduledCount)
	assert.Equal(t, 10, data.PublishedCount)
}

func TestNewAdminListData(t *testing.T) {
	articles := []model.Article{{Slug: "test"}}
	data := NewAdminListData("Drafts", articles, 100, 2, 10)

	assert.Equal(t, "Drafts", data.Title)
	assert.Len(t, data.Articles, 1)
	assert.Equal(t, 100, data.Total)
	assert.Equal(t, 2, data.Page)
	assert.Equal(t, 10, data.PageSize)
}

func TestAdminListData_TotalPages(t *testing.T) {
	tests := []struct {
		total    int
		pageSize int
		expected int
	}{
		{100, 10, 10},
		{101, 10, 11},
		{99, 10, 10},
		{0, 10, 0},
		{10, 0, 0},
	}

	for _, tt := range tests {
		data := &AdminListData{Total: tt.total, PageSize: tt.pageSize}
		assert.Equal(t, tt.expected, data.TotalPages())
	}
}

func TestNewAdminEditData_NewArticle(t *testing.T) {
	data := NewAdminEditData(nil, "", "")

	assert.Equal(t, "New Article", data.Title)
	assert.True(t, data.IsNew)
	assert.Nil(t, data.Article)
}

func TestNewAdminEditData_ExistingArticle(t *testing.T) {
	now := time.Now()
	article := &model.Article{
		Slug:        "test",
		Title:       "Test Article",
		Description: "A test",
		Date:        &now,
		Draft:       1,
	}

	data := NewAdminEditData(article, "# Content", "og_image: /test.jpg")

	assert.Equal(t, "Edit: Test Article", data.Title)
	assert.False(t, data.IsNew)
	assert.NotNil(t, data.Article)
	assert.Equal(t, "# Content", data.Content)
	assert.Equal(t, "og_image: /test.jpg", data.CustomYaml)
}

func TestAdminEditData_Map(t *testing.T) {
	now := time.Now()
	article := &model.Article{
		Slug:        "test",
		Title:       "Test Article",
		Description: "A test",
		Date:        &now,
		Draft:       1,
	}

	data := NewAdminEditData(article, "# Content", "")
	m := data.Map()

	assert.Equal(t, "Edit: Test Article", m["title"])
	assert.Equal(t, "# Content", m["bodyContent"])
	assert.False(t, m["isNew"].(bool))
	assert.Equal(t, "test", m["slug"])
	assert.Equal(t, "Test Article", m["articleTitle"])
	assert.True(t, m["draft"].(bool))
}
