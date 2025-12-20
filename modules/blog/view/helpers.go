package view

import "bytes"

// StripFrontMatter removes YAML front matter from content and returns just the body
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
