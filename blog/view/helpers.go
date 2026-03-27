package view

import (
	"bytes"
	"strings"
)

// StripFrontMatter removes YAML front matter from content and returns just the body.
func StripFrontMatter(content []byte) []byte {
	marker := []byte(`---`)

	// Check if content starts with ---
	if !bytes.HasPrefix(content, marker) {
		// No front matter, return as-is
		return content
	}

	// Find the closing ---
	parts := bytes.SplitN(content, marker, 3)
	if len(parts) < 3 {
		// No closing ---, return as-is
		return content
	}

	// Return content without front matter (parts[2])
	return parts[2]
}

// ExtractCustomYAML extracts non-standard YAML fields from frontmatter.
// Known fields (title, description, date, layout, draft) are excluded.
func ExtractCustomYAML(content []byte) string {
	marker := []byte(`---`)

	if !bytes.HasPrefix(content, marker) {
		return ""
	}

	parts := bytes.SplitN(content, marker, 3)
	if len(parts) < 3 {
		return ""
	}

	// Known fields to exclude
	knownFields := map[string]bool{
		"title":       true,
		"description": true,
		"date":        true,
		"layout":      true,
		"draft":       true,
	}

	frontmatter := string(parts[1])
	var customLines []string

	for _, line := range strings.Split(frontmatter, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Get the field name (before the colon)
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		fieldName := strings.TrimSpace(line[:colonIdx])
		if !knownFields[fieldName] {
			customLines = append(customLines, line)
		}
	}

	return strings.Join(customLines, "\n")
}
