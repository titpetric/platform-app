package view

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func getTemplatesFS(t *testing.T) fs.FS {
	// Get templates directory relative to this file
	wd, err := os.Getwd()
	require.NoError(t, err)

	templatesPath := filepath.Join(filepath.Dir(wd), "templates")
	return os.DirFS(templatesPath)
}

func TestRendererLogin(t *testing.T) {
	fsys := getTemplatesFS(t)
	renderer := NewRenderer(fsys, map[string]any{})
	ctx := context.Background()

	t.Run("login returns template", func(t *testing.T) {
		data := LoginData{
			Email: "test@example.com",
		}

		tpl := renderer.Login(data)
		require.NotNil(t, tpl)

		var buf bytes.Buffer
		err := tpl.Render(ctx, &buf)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "test@example.com")
	})

	t.Run("login with error message", func(t *testing.T) {
		data := LoginData{
			Email:        "user@example.com",
			ErrorMessage: "Invalid credentials",
		}

		tpl := renderer.Login(data)
		require.NotNil(t, tpl)

		var buf bytes.Buffer
		err := tpl.Render(ctx, &buf)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "user@example.com")
		assert.Contains(t, output, "Invalid credentials")
	})
}

func TestRendererLogout(t *testing.T) {
	fsys := getTemplatesFS(t)
	renderer := NewRenderer(fsys, map[string]any{})

	t.Run("logout returns template", func(t *testing.T) {
		data := LogoutData{
			User: "testuser",
		}

		tpl := renderer.Logout(data)
		require.NotNil(t, tpl)
	})
}

func TestRendererRegister(t *testing.T) {
	fsys := getTemplatesFS(t)
	renderer := NewRenderer(fsys, map[string]any{})
	ctx := context.Background()

	t.Run("register returns template", func(t *testing.T) {
		data := RegisterData{
			Email:     "newuser@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		tpl := renderer.Register(data)
		require.NotNil(t, tpl)

		var buf bytes.Buffer
		err := tpl.Render(ctx, &buf)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "newuser@example.com")
		assert.Contains(t, output, "John")
		assert.Contains(t, output, "Doe")
	})

	t.Run("register with error message", func(t *testing.T) {
		data := RegisterData{
			FirstName:    "Jane",
			LastName:     "Smith",
			Email:        "jane@example.com",
			ErrorMessage: "Email already exists",
		}

		tpl := renderer.Register(data)
		require.NotNil(t, tpl)

		var buf bytes.Buffer
		err := tpl.Render(ctx, &buf)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "jane@example.com")
		assert.Contains(t, output, "Jane")
		assert.Contains(t, output, "Smith")
		assert.Contains(t, output, "Email already exists")
	})
}
