package view

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMenuConfig(t *testing.T) {
	tmpDir := t.TempDir()
	menuFile := filepath.Join(tmpDir, "menu.yml")

	content := `
menu:
  - type: group
    label: Public
menuLoggedIn:
  - type: group
    label: Admin
`
	err := os.WriteFile(menuFile, []byte(content), 0o644)
	require.NoError(t, err)

	config, err := LoadMenuConfig(menuFile)
	require.NoError(t, err)

	assert.Len(t, config.Menu, 1)
	assert.Len(t, config.MenuLoggedIn, 1)
}

func TestLoadMenuConfig_FileNotFound(t *testing.T) {
	_, err := LoadMenuConfig("/nonexistent/menu.yml")
	require.Error(t, err)
}

func TestMenuConfig_GetMenu(t *testing.T) {
	config := &MenuConfig{
		Menu:         []any{"public"},
		MenuLoggedIn: []any{"admin"},
	}

	t.Run("not logged in", func(t *testing.T) {
		menu := config.GetMenu(false)
		assert.Equal(t, []any{"public"}, menu)
	})

	t.Run("logged in", func(t *testing.T) {
		menu := config.GetMenu(true)
		assert.Equal(t, []any{"admin"}, menu)
	})

	t.Run("logged in with empty menuLoggedIn", func(t *testing.T) {
		configNoAdmin := &MenuConfig{
			Menu:         []any{"public"},
			MenuLoggedIn: []any{},
		}
		menu := configNoAdmin.GetMenu(true)
		assert.Equal(t, []any{"public"}, menu)
	})
}

func TestNewMenuData(t *testing.T) {
	config := &MenuConfig{
		Menu:         []any{"public"},
		MenuLoggedIn: []any{"admin"},
	}

	data := NewMenuData(config, true)
	assert.Equal(t, []any{"admin"}, data.Menu)
	assert.True(t, data.LoggedIn)
}

func TestMenuData_Map(t *testing.T) {
	data := &MenuData{
		Menu:     []any{"test"},
		LoggedIn: true,
	}

	m := data.Map()
	assert.Equal(t, []any{"test"}, m["menu"])
	assert.True(t, m["loggedIn"].(bool))
}
