# Migrations

Migrations are templated for all applications. For any given app, the
process how to do migration repeats. Database schema is supported with
`mig gen`, generating the `model/` folder from schema.

Additionally, `mig docs` generates `schema/docs` markdown.

The database schema is treated as source of truth, and `mig lint`
verifies that the tables and columns are documented at source with
`COMMENT` for supported databases (MySQL, Postgres).

So, each package:

- `schema/` - source of truth,
- `schema/docs` - generated with `mig docs`,
- `model/*.mig.go` - generated with `mig gen`

The authorship areas for new modules are:

- `storage/` - implemented repositories to store/retrieve values,
- `storage/db.go` - a named connection for the application,
- `storage/storage.go` - implementation of retrieval with your supported database

While apps can be portable to some measure for simple CRUD operations,
it's not the default. The default path is to choose a database, design
the schema and cover it to standards, then write driver specific SQL for
the more complex queries. The mig model provides query builder utilities
with function chaining for common queries.

## Decisions

1. Why not `storage/model.go`?

In simple applications I want to encourage a separate model package.
This allows to build new functionality based on the data model, but not
couple to other implementation in the package.

This is why `model` and `storage` packages are kept separately.

With this choice, we favour internal composition. Complexity isn't
unusual with service oriented development, and a good measure of
grouping can occur within an application scope, for example:

- multiple `cmd/tool`'s
- multiple commands inside a tool (`git pull`, `git clone`,...)
- CQRS - separate dependency trees for read/write operations

A `model` package enables these subdivisions to packages beyond just
`storage`. With more complex applications, a `dto` or some similar
package should likely be maintained to handle type conversions and
validations, particularly around user provided input.
