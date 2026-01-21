package markdown

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	chroma "github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	blackfriday "github.com/russross/blackfriday/v2"
)

// Renderer renders markdown content to HTML with syntax highlighting
type Renderer struct {
	// codeBlockPattern matches HTML code blocks for highlighting
	codeBlockPattern *regexp.Regexp
}

// NewRenderer creates a new markdown renderer with syntax highlighting support
func NewRenderer() *Renderer {
	return &Renderer{
		codeBlockPattern: regexp.MustCompile(`<pre><code(?:\s+class="language-([^"]+)")?>([.\s\S]*?)</code></pre>`),
	}
}

// Render converts markdown content to HTML with syntax highlighting for code blocks
// Security: Blackfriday escapes HTML in code blocks, we unescape to get raw code,
// then Chroma's formatter escapes it again when outputting HTML. This is safe.
func (r *Renderer) Render(content []byte) []byte {
	// First, render markdown to HTML using blackfriday
	htmlContent := blackfriday.Run(content)

	// Then apply syntax highlighting to code blocks
	highlighted := r.highlightCodeBlocks(htmlContent)

	return highlighted
}

// highlightCodeBlocks applies syntax highlighting to all code blocks in the HTML
func (r *Renderer) highlightCodeBlocks(html []byte) []byte {
	htmlStr := string(html)

	// Replace all code blocks with highlighted versions
	result := r.codeBlockPattern.ReplaceAllStringFunc(htmlStr, func(match string) string {
		return r.processCodeBlock(match)
	})

	return []byte(result)
}

// processCodeBlock processes a single code block and applies syntax highlighting
func (r *Renderer) processCodeBlock(match string) string {
	// Extract language and code from the match
	m := r.codeBlockPattern.FindStringSubmatch(match)
	if len(m) < 3 {
		return match
	}

	language := m[1]
	code := m[2]

	// Unescape HTML entities in code (blackfriday escapes it, we need original for highlighting)
	code = unescapeHTML(code)

	// Apply syntax highlighting - Chroma handles escaping internally
	highlighted := highlightCode([]byte(code), language)

	return string(highlighted)
}

// highlightCode applies syntax highlighting to code using Chroma
func highlightCode(code []byte, language string) []byte {
	if len(code) == 0 {
		return []byte("<pre><code></code></pre>")
	}

	codeStr := string(code)

	// Select lexer for the language
	var lexer chroma.Lexer
	if language != "" {
		lexer = lexers.Get(language)
	}
	if lexer == nil {
		lexer = lexers.Analyse(codeStr)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Tokenize
	iterator, err := lexer.Tokenise(nil, codeStr)
	if err != nil {
		// On error, return plain code block
		return wrapCodePlain(code, language)
	}

	// Format with HTML and Monokai style
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.WithClasses(false),
		html.PreventSurroundingPre(true),
		html.TabWidth(4),
	)

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return wrapCodePlain(code, language)
	}

	// Don't strip newlines - preserve all whitespace as-is
	formattedCode := buf.String()

	// Wrap highlighted code in pre/code tags with chroma and language classes
	languageClass := ""
	if language != "" {
		languageClass = fmt.Sprintf(` language-%s`, language)
	}

	result := fmt.Sprintf(
		"<pre class=\"chroma%s\"><code>%s</code></pre>",
		languageClass,
		formattedCode,
	)

	return []byte(result)
}

// wrapCodePlain wraps code without highlighting
func wrapCodePlain(code []byte, language string) []byte {
	languageClass := ""
	if language != "" {
		languageClass = fmt.Sprintf(` class="language-%s"`, language)
	}
	return []byte(fmt.Sprintf(
		"<pre><code%s>%s</code></pre>",
		languageClass,
		escapeHTML(string(code)),
	))
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	).Replace(s)
}

// unescapeHTML unescapes HTML special characters
func unescapeHTML(s string) string {
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
	).Replace(s)
}
