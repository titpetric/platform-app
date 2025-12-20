# Markdown Package

The `markdown` package handles conversion of markdown content to HTML with automatic syntax highlighting for code blocks.

## Features

- **Markdown Rendering**: Converts markdown to HTML using `github.com/russross/blackfriday/v2`
- **Syntax Highlighting**: Automatically highlights code blocks using `github.com/alecthomas/chroma/v2`
- **Language Detection**: Detects code block language from fenced code block info strings
- **Fallback Detection**: Automatically detects code language if not specified
- **Monokai Theme**: Uses the popular Monokai color scheme for syntax highlighting

## Usage

```go
import "github.com/titpetric/platform-example/blog/markdown"

// Create a renderer
renderer := markdown.NewRenderer()

// Render markdown with syntax highlighting
html := renderer.Render([]byte("# Hello\n\n```go\nfmt.Println(\"hi\")\n```"))

// html now contains highlighted code block
```

## Supported Languages

Any language supported by Chroma lexers is automatically highlighted. Common languages include:

- Go
- Python
- JavaScript/TypeScript
- Java
- C/C++
- Rust
- SQL
- HTML/CSS
- Bash/Shell
- YAML
- JSON
- And many more...

## Code Block Syntax

### Fenced Code Blocks (Recommended)

Specify language for better highlighting:

```markdown
```go
package main

func main() {
    fmt.Println("Hello, World!")
}
```

```

### Indented Code Blocks

Code indented by 4 spaces is also converted to code blocks:

```markdown
    code line 1
    code line 2
```

## Styling

The CSS for syntax highlighting is provided in `theme/styles/syntax-highlighting.css`. Include this in your HTML templates:

```html
<link rel="stylesheet" href="/styles/syntax-highlighting.css">
```

## Implementation Details

The renderer works in two passes:

1. **Markdown Rendering**: Blackfriday converts markdown to HTML, escaping code content
2. **Syntax Highlighting**: Chroma unescapes and highlights the code, then re-escapes for safe HTML output

This approach ensures:
- All HTML is properly escaped for security
- Code blocks are correctly identified by regex pattern matching
- Language-specific syntax highlighting is applied

## Testing

Unit tests are available in `markdown_test.go`:

```bash
go test -v ./markdown
```

Tests cover:
- Simple markdown rendering
- Fenced code blocks with language specification
- Indented code blocks
- Code highlighting with and without language
- Edge cases (empty content, etc.)
