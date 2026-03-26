package view

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// MenuConfig holds both public and logged-in menu configurations.
type MenuConfig struct {
	Menu         []any `yaml:"menu"`
	MenuLoggedIn []any `yaml:"menuLoggedIn"`
}

// LoadMenuConfig loads the menu configuration from a file.
func LoadMenuConfig(filename string) (*MenuConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", filename, err)
	}
	defer f.Close()

	var config MenuConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	return &config, nil
}

// GetMenu returns the appropriate menu based on login state.
func (c *MenuConfig) GetMenu(loggedIn bool) []any {
	if loggedIn && len(c.MenuLoggedIn) > 0 {
		return c.MenuLoggedIn
	}
	return c.Menu
}

// MenuData holds the menu for template rendering.
type MenuData struct {
	Menu     []any
	LoggedIn bool
}

// NewMenuData creates MenuData with the appropriate menu for the user state.
func NewMenuData(config *MenuConfig, loggedIn bool) *MenuData {
	return &MenuData{
		Menu:     config.GetMenu(loggedIn),
		LoggedIn: loggedIn,
	}
}

// Map converts MenuData to a map[string]any.
func (d *MenuData) Map() map[string]any {
	return map[string]any{
		"menu":     d.Menu,
		"loggedIn": d.LoggedIn,
	}
}
