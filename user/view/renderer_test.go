package view

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed all:testdata
var testdata embed.FS

func TestRenderLayout(t *testing.T) {
	fsys, err := fs.Sub(testdata, "testdata")
	assert.NoError(t, err)

	renderer := NewRenderer(fsys, map[string]any{})
	ctx := context.Background()

	t.Run("blog.vuego", func(t *testing.T) {
		var buf bytes.Buffer
		err := renderer.Load("blog.vuego", map[string]any{
			"content": "Test Content",
		}).Render(ctx, &buf)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, ">Layout: base<")
		assert.Contains(t, output, ">Test Content<")
	})

	t.Run("index.vuego", func(t *testing.T) {
		var buf bytes.Buffer
		err := renderer.Load("index.vuego", map[string]any{
			"content": "Test Content",
		}).Render(ctx, &buf)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, ">Layout: base<")
		assert.Contains(t, output, ">Layout: post<")
		assert.Contains(t, output, ">Test Content<")
	})
}
