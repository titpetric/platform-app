package schema

import (
	yaml "gopkg.in/yaml.v3"
)

// Schema represents the entire CMS schema from schema.yaml
type Schema struct {
	Tables []*Table `yaml:"tables"`
}

// Table represents a database table definition
type Table struct {
	Name    string    `yaml:"name"`
	Comment string    `yaml:"comment"`
	Columns []*Column `yaml:"columns"`
	Indexes []*Index  `yaml:"indexes"`
}

// Column represents a table column definition
type Column struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"` // INTEGER, TEXT, TIMESTAMP, etc.
	Key        string            `yaml:"key"`  // PRI, MUL, UNI
	Comment    string            `yaml:"comment"`
	DataType   string            `yaml:"datatype"` // integer, text, enum, timestamp
	Values     []string          `yaml:"values"`
	EnumValues []string          `yaml:"enum_values"`
	DefaultVal string            `yaml:"default"`
	Nullable   bool              `yaml:"nullable"`
	Attributes map[string]string `yaml:"attributes,omitempty"`
}

// Index represents a table index definition
type Index struct {
	Name    string   `yaml:"name"`
	Columns []string `yaml:"columns"`
	Primary bool     `yaml:"primary"`
	Unique  bool     `yaml:"unique"`
}

// IsPrimaryKey returns true if this column is a primary key
func (c *Column) IsPrimaryKey() bool {
	return c.Key == "PRI"
}

// IsIndexed returns true if this column is indexed
func (c *Column) IsIndexed() bool {
	return c.Key == "MUL" || c.Key == "UNI" || c.Key == "PRI"
}

// IsTimestamp returns true if this column is a timestamp type
func (c *Column) IsTimestamp() bool {
	return c.DataType == "timestamp"
}

// IsEnum returns true if this column is an enum type
func (c *Column) IsEnum() bool {
	return c.DataType == "enum"
}

// GetLabel returns a user-friendly label from comment or name
func (c *Column) GetLabel() string {
	if c.Comment != "" {
		return c.Comment
	}
	return humanize(c.Name)
}

// GetFieldType returns the HTML input type for this column
func (c *Column) GetFieldType() string {
	switch c.DataType {
	case "integer":
		return "number"
	case "text":
		return "text"
	case "enum":
		return "select"
	case "timestamp":
		return "datetime-local"
	case "boolean":
		return "checkbox"
	default:
		return "text"
	}
}

// LoadSchema parses YAML bytes into a Schema
func LoadSchema(data []byte) (*Schema, error) {
	var schema Schema
	err := yaml.Unmarshal(data, &schema)
	if err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetTableNames returns a list of names of tables in the database.
func (s *Schema) GetTableNames() []string {
	result := make([]string, 0, len(s.Tables))
	for i := range s.Tables {
		result = append(result, s.Tables[i].Name)
	}
	return result
}

// Walk will call the callback for each table.
func (s *Schema) Walk(cb func(t *Table)) {
	for _, v := range s.Tables {
		cb(v)
	}
}

// GetTable returns a table by name or nil.
func (s *Schema) GetTable(name string) *Table {
	for _, v := range s.Tables {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// humanize converts snake_case to Title Case
func humanize(s string) string {
	// Simple humanization: replace _ with space and capitalize words
	words := ""
	prevUnderscore := true
	for _, ch := range s {
		if ch == '_' {
			words += " "
			prevUnderscore = true
		} else if prevUnderscore {
			words += string(toUpper(ch))
			prevUnderscore = false
		} else {
			words += string(ch)
		}
	}
	return words
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}
