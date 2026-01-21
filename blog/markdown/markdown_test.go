package markdown

import (
	"strings"
	"testing"
)

func TestRenderSimpleMarkdown(t *testing.T) {
	renderer := NewRenderer()
	markdown := []byte("# Hello\n\nThis is a **test**.")
	result := renderer.Render(markdown)

	if !strings.Contains(string(result), "<h1>") {
		t.Error("expected h1 tag in result")
	}
	if !strings.Contains(string(result), "<strong>test</strong>") {
		t.Error("expected strong tag in result")
	}
}

func TestRenderCodeBlockWithSyntaxHighlighting(t *testing.T) {
	renderer := NewRenderer()

	markdown := []byte("# Code Example\n\n```go\npackage main\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```\n")

	result := renderer.Render(markdown)
	resultStr := string(result)

	// Verify the code block is present
	if !strings.Contains(resultStr, "main") {
		t.Error("expected code content in result")
	}

	// Verify highlighting is applied (pre tag with chroma class)
	if !strings.Contains(resultStr, "class=\"chroma") {
		t.Error("expected chroma class in pre tag")
	}

	// Verify it's a code block
	if !strings.Contains(resultStr, "<code") {
		t.Error("expected code tag in result")
	}
}

func TestRenderPlainCodeBlock(t *testing.T) {
	renderer := NewRenderer()

	// Indented code block (not fenced)
	markdown := []byte("Some code:\n\n    plain code line 1\n    plain code line 2\n")

	result := renderer.Render(markdown)
	resultStr := string(result)

	if !strings.Contains(resultStr, "<pre") {
		t.Error("expected pre tag for code block")
	}
	if !strings.Contains(resultStr, "plain code line 1") {
		t.Error("expected code content in result")
	}
}

func TestRenderEmptyString(t *testing.T) {
	renderer := NewRenderer()
	result := renderer.Render([]byte(""))

	if result == nil {
		t.Error("expected non-nil result for empty input")
	}
}

func TestHighlightCodeWithLanguage(t *testing.T) {
	code := []byte("package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"test\") }")

	result := highlightCode(code, "go")

	if !strings.Contains(string(result), "go") {
		t.Error("expected go language indicator")
	}
	if !strings.Contains(string(result), "<code") {
		t.Error("expected code tag")
	}
}

func TestHighlightCodeWithoutLanguage(t *testing.T) {
	code := []byte("plain text code")
	result := highlightCode(code, "")

	if !strings.Contains(string(result), "plain text code") {
		t.Error("expected code content")
	}
}

func TestHighlightEmptyCode(t *testing.T) {
	result := highlightCode([]byte(""), "go")

	if !strings.Contains(string(result), "<pre>") {
		t.Error("expected pre tag")
	}
}
