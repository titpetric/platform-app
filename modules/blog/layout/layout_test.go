package layout

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
		err := renderer.Render(ctx, &buf, "blog.vuego", map[string]any{
			"content": "Test Content",
		})
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, ">Layout: base<")
		assert.Contains(t, output, ">Test Content<")
	})

	t.Run("index.vuego", func(t *testing.T) {
		var buf bytes.Buffer
		err := renderer.Render(ctx, &buf, "index.vuego", map[string]any{
			"content": "Test Content",
		})
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, ">Layout: base<")
		assert.Contains(t, output, ">Layout: post<")
		assert.Contains(t, output, ">Test Content<")
	})

	t.Run("multiple renders with same renderer", func(t *testing.T) {
		// First render
		var buf1 bytes.Buffer
		err := renderer.Render(ctx, &buf1, "blog.vuego", map[string]any{
			"content": "Content 1",
		})
		assert.NoError(t, err)
		output1 := buf1.String()
		assert.Contains(t, output1, "Content 1")

		// Second render with same renderer should not lose data
		var buf2 bytes.Buffer
		err = renderer.Render(ctx, &buf2, "blog.vuego", map[string]any{
			"content": "Content 2",
		})
		assert.NoError(t, err)
		output2 := buf2.String()
		assert.Contains(t, output2, "Content 2")
		assert.NotContains(t, output2, "Content 1")
	})
}
