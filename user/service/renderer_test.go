package service

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRendererLogin(t *testing.T) {
	renderer := NewRenderer(map[string]any{})
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
	renderer := NewRenderer(map[string]any{})

	t.Run("logout returns template", func(t *testing.T) {
		data := LogoutData{
			User: "testuser",
		}

		tpl := renderer.Logout(data)
		require.NotNil(t, tpl)
	})
}

func TestRendererRegister(t *testing.T) {
	renderer := NewRenderer(map[string]any{})
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
