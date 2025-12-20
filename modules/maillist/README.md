# Maillist

This is an example platform application.

## About

While the module isn't implemented, this demonstrates structure, and
handles some typical concerns for modules like migrating SQL schema.

- `schema/` and `schema.go` handle migrations,
- `main.go` registers a new module and runs migrations,
- `model/types.mig.go` is generated from database schema (`task migrate`)

The database connection for the maillist is set up separately to the
default connection. Requests against the platform user module use a
different connection, as defined in Taskfile.yml.

```yaml
    env:
      PLATFORM_DB_DEFAULT: "mysql://platform:platform@tcp(127.0.0.1:3306)/platform"
      PLATFORM_DB_MAILLIST: "mysql://platform:platform@tcp(127.0.0.1:3306)/platform_maillist"
```

This allows the maillist project to work on maillist schema, while
integrating with the platform for database access and authentication.

## Integrations

- [titpetric/platform-app](https://github.com/titpetric/platform-app) - A `user` module bringing in a user system
- [go-bridget/mig](https://github.com/go-bridget/mig) - SQL migration functionality
