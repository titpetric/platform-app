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
header:
  - label: Articles
    url: /blog
    icon: file-text
footer:
  - label: Articles
    url: /blog
  - label: Login
    url: /login
    loggedOut: true
`
	err := os.WriteFile(menuFile, []byte(content), 0o644)
	require.NoError(t, err)

	config, err := LoadMenuConfig(menuFile)
	require.NoError(t, err)

	assert.Len(t, config.Header, 1)
	assert.Len(t, config.Footer, 2)
	assert.Equal(t, "Articles", config.Header[0].Label)
	assert.Equal(t, "/blog", config.Header[0].URL)
	assert.True(t, config.Footer[1].LoggedOut)
}

func TestLoadMenuConfig_FileNotFound(t *testing.T) {
	_, err := LoadMenuConfig("/nonexistent/menu.yml")
	require.Error(t, err)
}

func TestLoadAdminMenuConfig(t *testing.T) {
	tmpDir := t.TempDir()
	menuFile := filepath.Join(tmpDir, "admin_menu.yml")

	content := `
admin:
  top:
    - label: Dashboard
      url: /admin
      icon: bi-speedometer2
  side:
    - label: Articles
      url: /admin/blog/articles
      icon: bi-file-text
  account:
    - label: Logout
      url: /logout
      icon: bi-box-arrow-right
`
	err := os.WriteFile(menuFile, []byte(content), 0o644)
	require.NoError(t, err)

	config, err := LoadAdminMenuConfig(menuFile)
	require.NoError(t, err)

	assert.Len(t, config.Admin.Top, 1)
	assert.Len(t, config.Admin.Side, 1)
	assert.Len(t, config.Admin.Account, 1)
	assert.Equal(t, "Dashboard", config.Admin.Top[0].Label)
}

func TestNewMenuData(t *testing.T) {
	config := &MenuConfig{
		Header: []MenuItem{{Label: "Articles", URL: "/blog"}},
		Footer: []MenuItem{{Label: "Login", URL: "/login", LoggedOut: true}},
	}

	data := NewMenuData(config, true)
	assert.Equal(t, config.Header, data.Header)
	assert.Equal(t, config.Footer, data.Footer)
	assert.True(t, data.LoggedIn)
}

func TestMenuData_Map(t *testing.T) {
	data := &MenuData{
		Header:   []MenuItem{{Label: "Test", URL: "/test"}},
		Footer:   []MenuItem{{Label: "Footer", URL: "/footer"}},
		LoggedIn: true,
	}

	m := data.Map()
	assert.Equal(t, data.Header, m["header"])
	assert.Equal(t, data.Footer, m["footer"])
	assert.True(t, m["loggedIn"].(bool))
}

func TestAdminNavigation_WithActive(t *testing.T) {
	nav := &AdminNavigation{
		Top: []MenuItem{
			{Label: "Dashboard", URL: "/admin", Icon: "bi-speedometer2"},
			{Label: "View Blog", URL: "/blog", Icon: "bi-box-arrow-up-right"},
		},
		Side: []MenuItem{
			{Label: "Articles", URL: "/admin/blog/articles", Icon: "bi-file-text"},
		},
		Account: []MenuItem{
			{Label: "Logout", URL: "/logout", Icon: "bi-box-arrow-right"},
		},
	}

	result := nav.WithActive("/admin")

	top := result["top"].([]map[string]any)
	assert.True(t, top[0]["active"].(bool))
	assert.False(t, top[1]["active"].(bool))

	side := result["side"].([]map[string]any)
	assert.False(t, side[0]["active"].(bool))
}
