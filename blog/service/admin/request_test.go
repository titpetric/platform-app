package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticleRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     ArticleRequest
		wantErr string
	}{
		{
			name:    "empty slug",
			req:     ArticleRequest{Title: "Test", Content: "Content"},
			wantErr: "slug is required",
		},
		{
			name:    "invalid slug format",
			req:     ArticleRequest{Slug: "Invalid Slug!", Title: "Test", Content: "Content"},
			wantErr: "slug must be lowercase alphanumeric with hyphens",
		},
		{
			name:    "slug too long",
			req:     ArticleRequest{Slug: "a-very-long-slug-that-exceeds-the-maximum-allowed-length-of-one-hundred-characters-which-is-way-too-long", Title: "Test", Content: "Content"},
			wantErr: "slug exceeds maximum length",
		},
		{
			name:    "empty title",
			req:     ArticleRequest{Slug: "test-slug", Content: "Content"},
			wantErr: "title is required",
		},
		{
			name:    "empty content",
			req:     ArticleRequest{Slug: "test-slug", Title: "Test"},
			wantErr: "content is required",
		},
		{
			name:    "invalid layout",
			req:     ArticleRequest{Slug: "test-slug", Title: "Test", Content: "Content", Layout: "invalid"},
			wantErr: "invalid layout",
		},
		{
			name:    "invalid date format",
			req:     ArticleRequest{Slug: "test-slug", Title: "Test", Content: "Content", Date: "15-01-2024"},
			wantErr: "date must be in YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS format",
		},
		{
			name:    "valid request",
			req:     ArticleRequest{Slug: "test-slug", Title: "Test", Content: "Content"},
			wantErr: "",
		},
		{
			name:    "valid request with layout",
			req:     ArticleRequest{Slug: "test-slug", Title: "Test", Content: "Content", Layout: "page"},
			wantErr: "",
		},
		{
			name:    "valid request with date",
			req:     ArticleRequest{Slug: "test-slug", Title: "Test", Content: "Content", Date: "2024-01-15"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestArticleRequest_ToArticle(t *testing.T) {
	req := ArticleRequest{
		Slug:        "test-article",
		Title:       "Test Article",
		Description: "A test article",
		Content:     "# Hello World",
		Date:        "2024-01-15",
		Draft:       true,
	}

	article := req.ToArticle()

	assert.Equal(t, "test-article", article.Slug)
	assert.Equal(t, "Test Article", article.Title)
	assert.Equal(t, "A test article", article.Description)
	assert.Equal(t, int64(1), article.Draft)
	assert.Equal(t, "post", article.Layout)
	assert.Equal(t, "/blog/test-article/", article.URL)
	assert.NotNil(t, article.Date)
}

func TestArticleRequest_BuildMarkdownContent(t *testing.T) {
	req := ArticleRequest{
		Title:       "Test Article",
		Description: "A test article",
		Date:        "2024-01-15",
		Content:     "# Hello World\n\nThis is content.",
		Draft:       true,
	}

	content := req.BuildMarkdownContent()

	assert.Contains(t, content, "---")
	assert.Contains(t, content, "title: \"Test Article\"")
	assert.Contains(t, content, "description: \"A test article\"")
	assert.Contains(t, content, "date: \"2024-01-15\"")
	assert.Contains(t, content, "draft: true")
	assert.Contains(t, content, "# Hello World")
}

func TestEscapeYAML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`simple`, `simple`},
		{`with "quotes"`, `with \"quotes\"`},
		{`with \backslash`, `with \\backslash`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeYAML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeContent(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldStrip string // substring that should be stripped
	}{
		{
			name:        "plain text preserved",
			input:       "Hello world",
			shouldStrip: "",
		},
		{
			name:        "script tag removed",
			input:       "Hello <script>alert('xss')</script> world",
			shouldStrip: "<script>",
		},
		{
			name:        "onclick removed",
			input:       `<div onclick="alert('xss')">Click me</div>`,
			shouldStrip: "onclick",
		},
		{
			name:        "onerror removed",
			input:       `<img src="x" onerror="alert('xss')">`,
			shouldStrip: "onerror",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeContent(tt.input)
			if tt.shouldStrip != "" {
				assert.NotContains(t, result, tt.shouldStrip, "dangerous content should be stripped")
			}
			// Plain text should be preserved
			if tt.input == "Hello world" {
				assert.Equal(t, tt.input, result)
			}
		})
	}
}
