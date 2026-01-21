package markdown

import (
	"strings"
	"testing"
)

func TestMarkdownCodeBlockRendering(t *testing.T) {
	renderer := NewRenderer()

	markdown := []byte("Here's some code:\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```\n\nEnd of example.")

	result := renderer.Render(markdown)
	resultStr := string(result)

	// Verify basic structure
	if !strings.Contains(resultStr, "<pre class=\"chroma") {
		t.Error("Missing <pre class=\"chroma\"> tag")
	}
	if !strings.Contains(resultStr, "language-go") {
		t.Error("Missing language-go class")
	}
	if !strings.Contains(resultStr, "package") {
		t.Error("Missing 'package' keyword")
	}
	if !strings.Contains(resultStr, "main") {
		t.Error("Missing 'main' identifier")
	}

	// Count and verify newlines in code block
	codeBlockStart := strings.Index(resultStr, "<code>")
	codeBlockEnd := strings.Index(resultStr, "</code>")
	if codeBlockStart == -1 || codeBlockEnd == -1 {
		t.Fatal("Could not find code block boundaries")
	}

	codeContent := resultStr[codeBlockStart+6 : codeBlockEnd]
	newlineCount := strings.Count(codeContent, "\n")

	// Original code has: package main\n\nimport\n\nfunc main()\n{\n\tfmt.Println\n}\n = 7 newlines
	if newlineCount < 5 {
		t.Errorf("Expected at least 5 newlines in code block, got %d", newlineCount)
	}

	// Count and verify tabs in code block
	tabCount := strings.Count(codeContent, "\t")
	if tabCount < 1 {
		t.Errorf("Expected at least 1 tab in code block, got %d", tabCount)
	}
}

func TestMarkdownMultipleCodeBlocks(t *testing.T) {
	renderer := NewRenderer()

	markdown := []byte(`First block:

` + "```" + `python
def hello():
	print("Hello")
	return True
` + "```" + `

Second block:

` + "```" + `javascript
function test() {
	console.log("test");
}
` + "```")

	result := renderer.Render(markdown)
	resultStr := string(result)

	// Count pre tags (should be 2)
	preCount := strings.Count(resultStr, "<pre class=\"chroma")
	if preCount != 2 {
		t.Errorf("Expected 2 code blocks, found %d", preCount)
	}

	// Verify both languages are present
	if !strings.Contains(resultStr, "language-python") {
		t.Error("Missing language-python class")
	}
	if !strings.Contains(resultStr, "language-javascript") {
		t.Error("Missing language-javascript class")
	}

	// Count total tabs across both blocks (should be 3)
	totalTabs := strings.Count(resultStr, "\t")
	if totalTabs < 3 {
		t.Errorf("Expected at least 3 tabs total, got %d", totalTabs)
	}
}

func TestMarkdownCodeBlockWhitespace(t *testing.T) {
	renderer := NewRenderer()

	// Test code with specific whitespace patterns
	markdown := []byte("```go\nfunc example() {\n\tif true {\n\t\tx := 1\n\t\ty := 2\n\t}\n}\n```")

	result := renderer.Render(markdown)
	resultStr := string(result)

	// Extract code block content
	codeBlockStart := strings.Index(resultStr, "<code>")
	codeBlockEnd := strings.Index(resultStr, "</code>")
	codeContent := resultStr[codeBlockStart+6 : codeBlockEnd]

	// Count tabs - should have nested indentation (6 tabs total)
	tabCount := strings.Count(codeContent, "\t")
	if tabCount < 4 {
		t.Errorf("Expected at least 4 tabs for nested indentation, got %d", tabCount)
	}

	// Count newlines - should preserve structure
	newlineCount := strings.Count(codeContent, "\n")
	if newlineCount < 5 {
		t.Errorf("Expected at least 5 newlines to preserve code structure, got %d", newlineCount)
	}
}

func TestCodeBlockEmptyLines(t *testing.T) {
	renderer := NewRenderer()

	markdown := []byte("```go\nfirst line\n\nsecond line after blank\n\nthird line\n```")

	result := renderer.Render(markdown)
	resultStr := string(result)

	// Extract code block content
	codeBlockStart := strings.Index(resultStr, "<code>")
	codeBlockEnd := strings.Index(resultStr, "</code>")
	codeContent := resultStr[codeBlockStart+6 : codeBlockEnd]

	// Should have 4 newlines for the blank lines and line breaks
	newlineCount := strings.Count(codeContent, "\n")
	if newlineCount < 4 {
		t.Errorf("Expected at least 4 newlines to preserve blank lines, got %d", newlineCount)
	}
}
