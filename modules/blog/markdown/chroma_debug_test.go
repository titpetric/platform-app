package markdown

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func TestChromaDirectOutput(t *testing.T) {
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

	// Get lexer
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Tokenize
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		t.Fatalf("tokenize failed: %v", err)
	}

	// Format with HTML
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.WithClasses(true),
		html.PreventSurroundingPre(true),
		html.TabWidth(4),
	)

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		t.Fatalf("format failed: %v", err)
	}

	output := buf.String()
	fmt.Printf("\n=== CHROMA RAW OUTPUT ===\n%q\n", output)
	fmt.Printf("\n=== CHROMA RENDERED ===\n%s\n", output)
	fmt.Printf("\n=== NEWLINE CHECK ===\n")
	fmt.Printf("Contains newlines: %v\n", bytes.Contains(buf.Bytes(), []byte("\n")))
	fmt.Printf("Newline count: %d\n", bytes.Count(buf.Bytes(), []byte("\n")))

	// Wrap like highlightCode does
	languageClass := ` language-go`
	result := fmt.Sprintf(
		"<pre class=\"chroma%s\"><code>%s</code></pre>",
		languageClass,
		output,
	)

	fmt.Printf("\n=== WRAPPED OUTPUT ===\n%q\n", result)
	fmt.Printf("\n=== WRAPPED RENDERED ===\n%s\n", result)
}
