package view

import (
	"fmt"
	"os"

	"github.com/titpetric/vuego"
	yaml "gopkg.in/yaml.v3"
)

// MenuItem represents a single navigation item.
type MenuItem struct {
	Label     string `yaml:"label"`
	URL       string `yaml:"url"`
	Icon      string `yaml:"icon,omitempty"`
	Target    string `yaml:"target,omitempty"`
	LoggedOut bool   `yaml:"loggedOut,omitempty"`
	LoggedIn  bool   `yaml:"loggedIn,omitempty"`
}

// Show returns true if this menu item should be displayed based on the login state.
func (m MenuItem) Show(ctx *vuego.VueContext) bool {
	loggedIn, _ := ctx.Stack().Resolve("loggedIn")
	isLoggedIn, _ := loggedIn.(bool)
	// Items with loggedOut: true are only shown when NOT logged in
	if m.LoggedOut {
		return !isLoggedIn
	}
	// Items with loggedIn: true are only shown when logged in
	if m.LoggedIn {
		return isLoggedIn
	}
	// Items without restrictions are always shown
	return true
}

// MenuConfig holds the main site navigation with separate header and footer.
type MenuConfig struct {
	Header []MenuItem `yaml:"header"`
	Footer []MenuItem `yaml:"footer"`
}

// AdminMenuConfig holds the admin panel navigation (separate from main site).
type AdminMenuConfig struct {
	Admin struct {
		Top     []MenuItem `yaml:"top"`
		Side    []MenuItem `yaml:"side"`
		Account []MenuItem `yaml:"account"`
	} `yaml:"admin"`
}

// LoadMenuConfig loads the main site menu configuration from a file.
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

// LoadAdminMenuConfig loads the admin panel menu configuration from a file.
func LoadAdminMenuConfig(filename string) (*AdminMenuConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", filename, err)
	}
	defer f.Close()

	var config AdminMenuConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	return &config, nil
}

// MenuData holds the navigation data for main site template rendering.
type MenuData struct {
	Header   []MenuItem
	Footer   []MenuItem
	LoggedIn bool
}

// NewMenuData creates MenuData for the main site.
func NewMenuData(config *MenuConfig, loggedIn bool) *MenuData {
	return &MenuData{
		Header:   config.Header,
		Footer:   config.Footer,
		LoggedIn: loggedIn,
	}
}

// Map converts MenuData to a map[string]any for template rendering.
func (d *MenuData) Map() map[string]any {
	return map[string]any{
		"header":   d.Header,
		"footer":   d.Footer,
		"loggedIn": d.LoggedIn,
	}
}

// AdminNavigation holds the admin panel navigation structure.
type AdminNavigation struct {
	Top     []MenuItem
	Side    []MenuItem
	Account []MenuItem
}

// NewAdminNavigation creates admin navigation from config.
func NewAdminNavigation(config *AdminMenuConfig) *AdminNavigation {
	return &AdminNavigation{
		Top:     config.Admin.Top,
		Side:    config.Admin.Side,
		Account: config.Admin.Account,
	}
}

// WithActive returns a copy of the navigation with the active item marked.
func (n *AdminNavigation) WithActive(activeURL string) map[string]any {
	markActive := func(items []MenuItem) []map[string]any {
		result := make([]map[string]any, len(items))
		for i, item := range items {
			result[i] = map[string]any{
				"label":  item.Label,
				"url":    item.URL,
				"icon":   item.Icon,
				"target": item.Target,
				"active": item.URL == activeURL,
			}
		}
		return result
	}

	return map[string]any{
		"top":     markActive(n.Top),
		"side":    markActive(n.Side),
		"account": markActive(n.Account),
	}
}
