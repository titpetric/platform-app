package schema

import (
	"embed"

	yaml "gopkg.in/yaml.v3"
)

// Migrations contains sql migrations contained in this folder.
//
//go:embed *.up.sql
var Migrations embed.FS

// SchemaYAML contains the schema definition.
//
//go:embed schema.yml
var SchemaYAML []byte

// Table represents a database table schema.
type Table struct {
	Name    string   `yaml:"name"`
	Comment string   `yaml:"comment"`
	Columns []Column `yaml:"columns"`
}

// Column represents a database column schema.
type Column struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Key      string `yaml:"key"`
	Comment  string `yaml:"comment"`
	DataType string `yaml:"datatype"`
}

// LoadSchema parses the embedded schema.yml and returns the tables.
func LoadSchema() ([]Table, error) {
	var tables []Table
	if err := yaml.Unmarshal(SchemaYAML, &tables); err != nil {
		return nil, err
	}
	return tables, nil
}

// GetTable returns a specific table by name.
func GetTable(name string) (*Table, error) {
	tables, err := LoadSchema()
	if err != nil {
		return nil, err
	}
	for i := range tables {
		if tables[i].Name == name {
			return &tables[i], nil
		}
	}
	return nil, nil
}
