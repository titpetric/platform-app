# Platform app

This is an example [titpetric/platform](https://github.com/titpetric/platform) application tree.

## Modules

A module is created in the root folder. Currently we have several modules:

- `expvar` - a simple monitoring endpoint under `/debug/vars`.
- `user` - a user registration, login and logout MVC.
- `assets` - static files to be served
- `theme` - the layout shared by modules for output
- `daily` - a simple sqlite focused TODO-app
- `blog` - a blog based on ryan mulligan dev

Other modules may be added, some may be removed as they serve mostly as
an example, and not a fully developed app/service.

For modules, the conventions followed are some subset of:

- `storage (5)` - provide database access APIs, repository
- `schema (5)` - database migration, single/multi DB
- `docs (5)` - any documentation, schema docs, technical
- `model (5)` - generated app data model
- `view (3)` - a view package for MVC components
- `templates (3)` - a templates package for Vuego templates
- `service (2)` - the module implementation, app
- `data (2)` - usually some source of local input
- `cmd (2)` - modules provide their cmd's
- `main.go` in root of the module or expanded `service/` (package encouraged)

## Folder structure

This document serves as a structural guide for the repository:

- `./docker` - docker build environments for apps
- `./docs` - some generated docs, this structure guide
- `./cmd` - the initial app entrypoint (deprecated)
- `./autoload` - registers the default platform drivers (deprecated)

No other root packages except the modules are expected.

Experimental:

- `./themes` - experimental code for themes via go packages
- `./cms` - experimental code for a cms, non functional

The deprecated and experimental options are meant to be removed over
time. No packages here are intended for use outside of the platform-app
repository, and the tendency is to only keep the module packages in the
root of the repo.

## Tools used

The tools used are:

- [Vuego template engine](https://github.com/titpetric/vuego)
- [MIG - schema migration tooling](https://github.com/go-bridget/mig)
- [Atkins runner](https://github.com/titpetric/atkins)

Taskfiles are being phased out for a combination of Atkins and Leftpad (for pre-commit hooks).

## Docs

- [Testing coverage report](./docs/testing-coverage.md)
- [Code structure diagrams](./docs/structure.md)