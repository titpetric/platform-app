# Building a platform-app application

This guide describes the repository conventions for adding an application module. It is data-model first: define and generate the database contract before writing storage, HTTP, or frontend code.

## 1. Start with SQL migrations and generate the model

Create the module at the repository root and write its initial migration as `<module>/schema/<module>.up.sql`. Migration files are append-only after they have been applied anywhere shared; later changes belong in another `*.up.sql` file.

Use portable SQLite SQL unless the module explicitly targets another database:

- Use singular, lowercase `snake_case` table names: `schedule`, `booking`, `booking_administrator`.
- Use lowercase `snake_case` columns.
- Use `CHAR(26)` for ULID identifiers and related IDs. Generate values in Go with `platform/pkg/ulid.String()`; databases do not generate application IDs.
- Name user relationships `user_id`, `created_by_user_id`, or another explicit `*_user_id`. Do not add SQL foreign keys: modules can share identity values without coupling migrations or database connections.
- Store lifecycle timestamps as `DATETIME`: normally `created_at`, `updated_at`, and nullable `deleted_at`. Application queries must exclude `deleted_at IS NOT NULL` rows by default.
- Name boolean state positively with `is_`, such as `is_active` and `is_admin`, and declare it `BOOLEAN NOT NULL DEFAULT 0` (or `1` when active by default).
- Use `NOT NULL` and defaults deliberately. Nullable values should represent meaningful absence, not an omitted design decision.
- Add indexes for ownership, active-list, ordering, and lookup predicates used by storage. Partial unique indexes are useful for active records in SQLite.
- Prefer explicit status text plus a `CHECK` constraint when a record has more than two lifecycle states.
- Keep table comments and index names descriptive. Prefix indexes with `idx_` and unique indexes with `uidx_`.

After writing the SQL, run the repository's existing schema automation from the new module directory:

```sh
atkins "$HOME/.atkins/skills/schema.yml" -w ./booking migrate
```

The schema skill creates a disposable `<module>.db`, applies `schema/*.sql`, generates Markdown and `schema/schema.yml`, then runs `mig gen --go.fill-json --output=./model`. Commit the generated `model/types.mig.go`, schema YAML, and useful schema docs. Generated Go files say `DO NOT EDIT`; change SQL and regenerate instead.

Embed the migrations in `<module>/schema/schema.go` with `//go:embed *.up.sql`. The storage constructor should apply that filesystem through `mig/migrate`, so a module initializes its own database consistently in development, tests, and production.

## 2. Create packages around ownership boundaries

A developed module normally has this shape:

```text
<module>/
├── <module>.go            platform module lifecycle and route assembly
├── model/
│   └── types.mig.go       generated database records and query helpers
├── schema/
│   ├── <module>.up.sql    append-only migrations
│   ├── schema.go          embedded migrations
│   ├── schema.yml         generated schema description
│   └── docs/              generated schema documentation
├── storage/
│   ├── db.go              platform database connection
│   ├── migrate.go         embedded migration runner
│   └── ...                persistence and transactions
└── service/
    ├── api/               authenticated user JSON handlers
    ├── admin/             explicitly authorized admin JSON handlers
    └── web/               optional HTML handlers, added only when needed
```

`model` represents persisted data, not request payloads or business operations. Put request/response types beside the HTTP handlers when they are transport-specific.

`storage` owns SQL, transactions, user scoping, capacity/uniqueness invariants, and conversion of zero affected rows to `sql.ErrNoRows`. Keep multi-write state transitions such as booking, cancellation promotion, or rescheduling in one transaction.

`service/api` owns decoding, input validation, HTTP status mapping, and JSON encoding. It calls storage rather than issuing SQL. `service/admin` is a separate package and route namespace even when it currently serves only JSON; a future frontend can consume the same API without mixing templates into persistence.

The root module implements `platform.Module`: connect and migrate in `Start`, assemble handlers after storage is ready, mount routes in `Mount`, and report a stable `Name`. Do not build frontend packages until there is actual frontend behavior to own.

## 3. Authenticate with the user module

Populate browser/session authentication by mounting user middleware around the route group. For JSON APIs whose storage layer requires a user, use optional mode so an absent session reaches the handler and can be returned as a consistent JSON `401` (the current required middleware stops the chain without writing an HTTP error):

```go
r.Group(func(r platform.Router) {
	r.Use(user.NewMiddleware(user.AuthCookie(), user.AuthOptional()))
	// authenticated routes
})
```

For bearer-token APIs use `user.AuthHeader()`. Optional middleware is not authorization by itself: every protected handler or storage operation must reject a missing session user. It is also appropriate for genuinely public routes that render extra state for logged-in users.

The middleware validates the session or JWT, loads the active user, and writes it to the request context. Retrieve it with:

```go
sessionUser, ok := user.GetSessionUser(r.Context())
if !ok {
	return user.ErrLoginRequired
}
userID := sessionUser.ID
```

Storage operations for user-owned records should also read the session user from `context.Context` and include `user_id = ?` in reads and writes. That provides defense in depth and prevents handlers from accepting a caller-supplied owner ID. Tests can bind a user with `user.SetSessionUser(ctx, &model.User{ID: ...})`.

Authentication is not administration. A route named `/admin` must check an application-owned role or membership after user authentication. Keep that check in storage/service authorization and return forbidden for authenticated non-admin users. Plan how the first administrator is provisioned (migration, operational command, or direct controlled insert); do not make the first caller an administrator automatically.

## 4. Implement and test lifecycle behavior

Write storage tests against an isolated SQLite database with embedded migrations. Use `testify/assert` for result checks and `testify/require` for setup and preconditions. At minimum, test:

- migrations create all expected tables and indexes;
- authentication is required and ownership is enforced;
- normal create/list/get/update and soft-delete behavior;
- uniqueness and capacity boundaries;
- transactional state transitions and rollback paths;
- administrator allow/deny behavior;
- HTTP decoding, statuses, and response bodies with `httptest`.

For time-based applications, inject or pass explicit times where practical and store UTC. Test the boundary rules rather than sleeping. For concurrent limits such as schedule capacity, calculate and write inside a transaction; do not trust a prior list response to remain current.

Run the narrow module suite first:

```sh
go test ./booking/...
```

Then format and run broader checks when shared contracts or module assembly changed. The repository structural linter is useful for confirming package shape:

```sh
go run ./internal/cmd/structurelint
```

## 5. Keep generated and handwritten responsibilities clear

Review generated diffs after every schema run. Never add business methods to `types.mig.go`; put them in another handwritten file in `model` only when they are model behavior, or in storage/service when they depend on persistence or a request. Do not expose generated SQL helpers as an HTTP contract.

Document status meanings and transition rules near the owning module. Schema constraints protect valid values, storage transactions protect cross-row invariants, API validation protects callers, and tests preserve the intended lifecycle.
