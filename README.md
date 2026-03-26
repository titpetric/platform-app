# platform-app - An example platform application tree

This is an example [titpetric/platform](https://github.com/titpetric/platform) application tree.

It is built with:

- [Vuego template engine](https://github.com/titpetric/vuego)
- [MIG - schema migration tooling](https://github.com/go-bridget/mig)
- [Atkins runner](https://github.com/titpetric/atkins)

It implements full applications to deploy via docker. The code here has
been used to provide:

- https://vuego.incubator.to
- https://pulse.incubator.to
- https://atkins.incubator.to

The example is a modular system with a structural linter to use as a guide.

Each module is expected to have it's own test suite with `atkins`.

## Structural linter

There is a structural linter you can run in the project:

```
go run ./internal/cmd/structurelint
```

It will print the following output, giving a rudimentary maturity score.

| folder   | README.md | docs | compose.yml | docker | cmd/*          | config | model | schema | storage | service | service/*               | view | *                               |
| -------- | --------- | ---- | ----------- | ------ | -------------- | ------ | ----- | ------ | ------- | ------- | ----------------------- | ---- | ------------------------------- |
| blog     | [x]       | [x]  | [ ]         | [ ]    | blog, generate | [ ]    | [x]   | [x]    | [x]     | [ ]     |                         | [x]  | data, handlers, markdown, theme |
| cms      | [ ]       | [ ]  | [ ]         | [ ]    |                | [ ]    | [ ]   | [x]    | [ ]     | [ ]     |                         | [ ]  | templates, testdata             |
| daily    | [ ]       | [x]  | [ ]         | [ ]    |                | [ ]    | [x]   | [x]    | [x]     | [ ]     |                         | [x]  | templates                       |
| email    | [x]       | [x]  | [ ]         | [ ]    |                | [ ]    | [x]   | [x]    | [x]     | [ ]     |                         | [ ]  | smtp                            |
| expvar   | [ ]       | [ ]  | [ ]         | [ ]    |                | [ ]    | [ ]   | [ ]    | [ ]     | [ ]     |                         | [ ]  |                                 |
| internal | [ ]       | [ ]  | [ ]         | [ ]    | structurelint  | [ ]    | [ ]   | [ ]    | [ ]     | [ ]     |                         | [ ]  |                                 |
| maillist | [x]       | [x]  | [ ]         | [ ]    |                | [ ]    | [x]   | [x]    | [x]     | [ ]     |                         | [ ]  |                                 |
| pulse    | [x]       | [ ]  | [x]         | [x]    | pulse          | [x]    | [x]   | [x]    | [x]     | [x]     | keycounter              | [x]  | client                          |
| user     | [ ]       | [x]  | [ ]         | [ ]    | generate       | [ ]    | [x]   | [x]    | [x]     | [x]     | api, auth, passkey, web | [x]  | opa                             |

## Modules

A module is created in the root folder.

For modules, the conventions followed are some subset of:

- `storage` - provide database access APIs, repository
- `schema` - database migration, single/multi DB
- `docs` - any documentation, schema docs, technical
- `model` - generated app data model
- `view` - a view package for MVC components
- `templates` - a templates package for Vuego templates
- `service` - the module implementation, app
- `data` - usually some source of local input
- `cmd` - modules provide their cmd's
- `main.go` in root of the module or expanded `service/` (package encouraged)

## Folder structure

This document serves as a structural guide for the repository:

- `./docker` - docker build environments for apps
- `./docs` - some generated docs, this structure guide

No other root packages except the modules are expected.

Experimental:

- `./themes` - experimental code for themes via go packages
- `./cms` - experimental code for a cms, non functional

The deprecated and experimental options are meant to be removed over
time. No packages here are intended for use outside of the platform-app
repository, and the tendency is to only keep the module packages in the
root of the repo.

## Docs

- [Testing coverage report](./docs/testing-coverage.md)
- [Code structure diagrams](./docs/structure.md)
