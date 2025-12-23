package cms

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// LoadFormFromYAML loads a form definition from YAML data
func LoadFormFromYAML(name string, data []byte) (*Form, error) {
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Get the module data (first key in the map)
	var moduleData map[string]interface{}
	for _, v := range raw {
		if m, ok := v.(map[string]interface{}); ok {
			moduleData = m
			break
		}
	}

	if moduleData == nil {
		return nil, fmt.Errorf("no module data found in YAML")
	}

	form := &Form{
		Name:   name,
		Title:  name,
		Fields: []Field{},
	}

	// Convert each field
	for fieldName, fieldData := range moduleData {
		if fieldMap, ok := fieldData.(map[string]interface{}); ok {
			field := fieldFromMap(fieldName, fieldMap)
			form.Fields = append(form.Fields, field)
		}
	}

	return form, nil
}

// fieldFromMap converts a YAML field definition to a Field struct
func fieldFromMap(name string, data map[string]interface{}) Field {
	field := Field{
		Name:     name,
		Title:    name,
		Required: false,
		Disabled: false,
		Multiple: false,
		Readonly: false,
	}

	// Map YAML fields to Field struct
	if desc, ok := data["desc"].(string); ok {
		field.Desc = desc
		field.Title = cleanTitle(desc)
	}

	if fieldType, ok := data["type"].(string); ok {
		field.Type = fieldType
	}

	if value, ok := data["value"]; ok {
		field.Value = value
		// Parse enum-style values
		if str, ok := value.(string); ok && field.Type == "enum" {
			field.Options = parseEnumValue(str)
		}
	}

	if placeholder, ok := data["placeholder"].(string); ok {
		field.Placeholder = placeholder
	}

	if size, ok := data["size"].(string); ok {
		field.Size = size
	}

	if height, ok := data["height"].(int); ok {
		field.Height = height
	}

	if rows, ok := data["rows"].(int); ok {
		field.Rows = rows
	}

	if cols, ok := data["cols"].(int); ok {
		field.Cols = cols
	}

	if link, ok := data["link"].(string); ok {
		field.Link = link
	}

	if sort, ok := data["sort"].(string); ok {
		field.Sort = sort
	}

	if accept, ok := data["accept"].(string); ok {
		field.Accept = accept
	}

	// Determine if field is required based on type
	field.Required = isRequiredField(field.Type)

	return field
}

// parseEnumValue parses "0=Label,1=Label" format into Options
func parseEnumValue(value string) []Option {
	var options []Option
	pairs := strings.Split(value, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), "=")
		if len(parts) == 2 {
			options = append(options, Option{
				Value: strings.TrimSpace(parts[0]),
				Label: strings.TrimSpace(parts[1]),
			})
		}
	}
	return options
}

// cleanTitle converts field description to a clean title
func cleanTitle(desc string) string {
	// Remove HTML tags
	desc = strings.ReplaceAll(desc, "<br>", " ")
	desc = strings.ReplaceAll(desc, "<br/>", " ")
	desc = strings.ReplaceAll(desc, "<br />", " ")
	// Trim spaces
	desc = strings.TrimSpace(desc)
	return desc
}

// isRequiredField determines if a field type is typically required
func isRequiredField(fieldType string) bool {
	requiredTypes := map[string]bool{
		"key":                 false,
		"parent_key":          false,
		"explicit_parent_key": false,
		"hidden":              false,
		"stamp":               false,
	}
	return !requiredTypes[fieldType]
}
