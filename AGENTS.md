# Agent Instructions

Guidelines for AI agents working on this codebase.

## SQL Queries

When writing SQL queries with multiple columns, use explicit SELECT with formatted columns aligned to struct fields:

```go
query := `SELECT 
    id, recipient, subject, 
    '' as body, 
    created_at, 
    retry_count, 
    NULL as last_error, 
    NULL as last_retry, 
    sent_at, 
    NULL as error, 
    'sent' as status
FROM email_sent ORDER BY sent_at DESC LIMIT ?`
```

Do NOT use:
- `SELECT *` when columns don't match struct field order
- Single-line long column lists
- Non-explicit column ordering

This ensures:
- Code is maintainable and columns are visible at a glance
- Struct field scanning works correctly with named `db:` tags
- NULL values are explicit and intentional
- Debugging is easier when columns don't match expectations

## Email Module Testing

When running email module tests, ensure:
- Storage queries match the actual table names in `schema/model.up.sql`
- The migration runs via `Migrate()` before tests access tables
- Tests import `_ "github.com/titpetric/platform/pkg/drivers"` to register sqlite driver
- If migrations show "model.up.sql OK" but tables don't exist, check that table names are consistent between schema and storage queries
