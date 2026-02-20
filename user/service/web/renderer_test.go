package web_test

import (
	"bytes"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/vuego"
	"github.com/titpetric/vuego-cli/basecoat"

	"github.com/titpetric/platform-app/user/service/web"
	"github.com/titpetric/platform-app/user/view"
)

func newViewFS() fs.FS {
	return vuego.NewOverlayFS(basecoat.FS, view.FS)
}

func TestRendererLogin(t *testing.T) {
	renderer := web.NewRenderer(newViewFS(), map[string]any{})

	t.Run("login returns template", func(t *testing.T) {
		ctx := t.Context()
		data := web.LoginData{
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
		ctx := t.Context()
		data := web.LoginData{
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
	renderer := web.NewRenderer(newViewFS(), map[string]any{})

	t.Run("logout returns template", func(t *testing.T) {
		data := web.LogoutData{
			User: "testuser",
		}

		tpl := renderer.Logout(data)
		require.NotNil(t, tpl)
	})
}

func TestRendererRegister(t *testing.T) {
	renderer := web.NewRenderer(newViewFS(), map[string]any{})

	t.Run("register returns template", func(t *testing.T) {
		ctx := t.Context()
		data := web.RegisterData{
			Email:    "newuser@example.com",
			FullName: "John Doe",
		}

		tpl := renderer.Register(data)
		require.NotNil(t, tpl)

		var buf bytes.Buffer
		err := tpl.Render(ctx, &buf)
		require.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "newuser@example.com")
		assert.Contains(t, output, "John Doe")
	})

	t.Run("register with error message", func(t *testing.T) {
		ctx := t.Context()
		data := web.RegisterData{
			FullName:     "Jane Smith",
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
