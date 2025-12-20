# Agent Commands and Conventions

This document defines standard commands and development conventions for the blog module.

## Build Commands

All build and test commands use [Task](https://taskfile.dev). Install it with:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### Common Commands

**Run the service:**

```bash
task          # Start blog service on http://127.0.0.1:8080
```

**Format code:**

```bash
task fmt
```

**Run all checks (lint, test, build):**

```bash
task all
```

**Run tests:**

```bash
task test          # All tests
task test:unit     # Unit tests only
task test:integration  # Integration tests only
task test:coverage # Tests with coverage report
```

**Build:**

```bash
task build
```

**Run benchmarks:**

```bash
task bench         # Standard benchmarks
task bench:cpu     # With CPU profiling
task bench:memory  # With memory profiling
```

**Lint and vet:**

```bash
task lint
task vet
```

**Clean:**

```bash
task clean         # Remove build artifacts
```

**CI pipeline:**

```bash
task ci            # Runs: fmt:check, lint, vet, test, coverage:check
```

## Code Style

- Use `goimports` for import organization
- Run `task fmt` before committing
- Maximum line length: 100 characters (enforced by linter)
- Use meaningful variable names
- Add comments to exported functions

## Template Syntax (Vuego)

- Use `{{ variable }}` for variable interpolation
- Use pipes for filters: `{{ value | filter }}`
- Use `<vuego include="path/to/component.vuego">` for component composition
- Include paths are relative to the filesystem root passed to vuego.NewVue()
- Use `v-if` for conditional rendering
- Use `v-for="item in array"` for loops (not "of")
- Use `v-html="content"` for unescaped HTML
- Use `:attr="value"` for attribute binding

## Testing Standards

### Unit Tests
- Files: `*_test.go` in same package
- Use `setupTestDB()` and `cleanupTestDB()` helpers
- Each test is independent and can run in isolation
- Mock external dependencies

### Integration Tests
- Files: `integration_test.go` in `storage/` package
- Test full workflows across multiple functions
- Use real SQLite database (`:memory:`)
- Test error conditions and edge cases

### Benchmarks
- Files: `*_test.go` with `Benchmark` prefix
- Use `task bench` to run
- Profile with `task bench:cpu` or `task bench:memory`

## Code Organization

```
blog/
├── model/              # Data types and structures
├── markdown/           # Markdown rendering with syntax highlighting
│   └── markdown.go     # Markdown to HTML conversion with Chroma
├── storage/            # Database layer
│   ├── storage.go      # Public interface
│   ├── articles.go     # SQL operations
│   ├── db.go          # Connection helper
│   ├── storage_test.go # Unit tests
│   └── integration_test.go  # Integration tests
├── view/               # Template rendering
├── theme/              # Static theme assets
├── handlers.go        # HTTP handlers
├── blog.go           # Module implementation
└── schema/           # Database schema migrations
```

## Git Workflow

1. Create feature branch: `git checkout -b feature/name`
2. Make changes and commit: `git add . && git commit -m "description"`
3. Run checks: `task ci`
4. Ensure all tests pass before pushing
5. Push branch and create pull request

## Documentation

- Update relevant `.md` files when changing behavior
- Keep README.md in sync with implementation
- Document breaking changes in commit messages
- Use clear, concise language

## Dependencies

### Go Modules

Update dependencies with:

```bash
go get -u ./...
go mod tidy
```

Check for updates:

```bash
task deps
```

### Required Tools

- Go 1.25.4+
- SQLite 3.x
- Task (for build automation)
- golangci-lint (for code quality)

### Core Dependencies

- `github.com/jmoiron/sqlx` - Database access
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/titpetric/platform` - Platform framework
- `github.com/russross/blackfriday/v2` - Markdown rendering
- `github.com/alecthomas/chroma/v2` - Syntax highlighting for code blocks
- `github.com/titpetric/vuego` - Template engine
- `gopkg.in/yaml.v3` - YAML parsing

## Common Patterns

### Storage Operations

All storage methods follow this pattern:

```go
func (s *Storage) GetArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	return GetArticleBySlug(ctx, s.db, slug)
}
```

Public methods on `Storage` delegate to package-level functions that accept `*sqlx.DB`.

### Error Handling

- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Propagate errors up; don't log in utility functions
- Use specific error types for client-facing errors

### Context Usage

All database operations accept `context.Context`:
- Enables request cancellation
- Supports distributed tracing
- Pass `r.Context()` from HTTP handlers
- Use `context.Background()` in tests

## Continuous Integration

The `task ci` command runs:
1. `task fmt:check` - Verify code formatting
2. `task lint` - Run golangci-lint
3. `task vet` - Run go vet
4. `task test` - Run all tests
5. `task test:coverage:check` - Verify coverage ≥ 70%

## Debugging

### Database Issues

Check SQLite directly:

```bash
sqlite3 :memory: "SELECT COUNT(*) FROM articles;"
```

### Test Failures

Run specific test with verbose output:

```bash
go test -v -run TestName ./storage
```

### Performance

Profile a benchmark:

```bash
task bench:cpu
task bench:memory
```

## Schema Management

Database schema lives in `schema/`:
- Forward-only migrations (no rollback support)
- Single migration file: `2025-01-01-000000-articles-initial.up.sql`
- Embedded in code via `//go:embed`
- Applied automatically on module startup

Never edit schema outside migrations. To change schema:
1. Create new migration file (timestamp-based)
2. Test with `task test`
3. Document breaking changes
