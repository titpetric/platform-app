package blog

import (
	yaml "gopkg.in/yaml.v3"
)

// Navigation represents CMS navigation structure
type Navigation struct {
	Name  string
	Links []NavigationLink
}

// NavigationLink represents a single navigation link
type NavigationLink struct {
	Label string
	URL   string
	Icon  string // Optional icon identifier
}

// Form represents a complete form definition
type Form struct {
	Name   string            `json:"name" yaml:"name"`
	Title  string            `json:"title" yaml:"title"`
	Fields []Field           `json:"fields" yaml:"fields"`
	Meta   map[string]string `json:"meta,omitempty" yaml:"meta,omitempty"`
}

// Field represents a form field definition
type Field struct {
	Name        string            `json:"name" yaml:"name"`
	Title       string            `json:"title" yaml:"title"`
	Desc        string            `json:"desc" yaml:"desc"` // Description (from formidable)
	Type        string            `json:"type" yaml:"type"`
	Value       interface{}       `json:"value,omitempty" yaml:"value,omitempty"`
	Placeholder string            `json:"placeholder,omitempty" yaml:"placeholder,omitempty"`
	Required    bool              `json:"required" yaml:"required"`
	Options     []Option          `json:"options,omitempty" yaml:"options,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`
	Size        string            `json:"size,omitempty" yaml:"size,omitempty"` // full, half
	Height      int               `json:"height,omitempty" yaml:"height,omitempty"`
	Rows        int               `json:"rows,omitempty" yaml:"rows,omitempty"`
	Cols        int               `json:"cols,omitempty" yaml:"cols,omitempty"`
	Link        string            `json:"link,omitempty" yaml:"link,omitempty"` // For search_popup
	Sort        string            `json:"sort,omitempty" yaml:"sort,omitempty"`
	Accept      string            `json:"accept,omitempty" yaml:"accept,omitempty"` // For file uploads
	Multiple    bool              `json:"multiple" yaml:"multiple"`
	Disabled    bool              `json:"disabled" yaml:"disabled"`
	Readonly    bool              `json:"readonly" yaml:"readonly"`
	Error       string            `json:"error,omitempty" yaml:"error,omitempty"`
	Help        string            `json:"help,omitempty" yaml:"help,omitempty"`
}

// Option represents a select/radio option
type Option struct {
	Value string `json:"value" yaml:"value"`
	Label string `json:"label" yaml:"label"`
}

// ParseFieldValue parses enum-style values (key=label,key=label) into Options
func ParseFieldValue(value interface{}) []Option {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		return parseEnumString(v)
	case []interface{}:
		return parseEnumArray(v)
	}
	return nil
}

// parseEnumString handles "key=label,key=label" format
func parseEnumString(s string) []Option {
	var opts []Option
	// This would parse CSV-style enum values
	// Implementation depends on exact format
	return opts
}

// parseEnumArray handles YAML array of options
func parseEnumArray(arr []interface{}) []Option {
	var opts []Option
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			opt := Option{}
			if v, ok := m["key"]; ok {
				opt.Value = v.(string)
			}
			if v, ok := m["value"]; ok {
				opt.Label = v.(string)
			}
			if opt.Value != "" && opt.Label != "" {
				opts = append(opts, opt)
			}
		}
	}
	return opts
}

// UnmarshalYAML customizes YAML unmarshaling for Form
func (f *Form) UnmarshalYAML(node *yaml.Node) error {
	// Map form data from YAML structure
	type formAlias Form
	return node.Decode((*formAlias)(f))
}
