# Platform app

This is an example platform application. It provides two modules:

- `expvar` - a simple monitoring endpoint under `/debug/vars`.
- `user` - a user registration, login and logout MVC.

The functionality of the `expvar` package is provided by the platform.
Telemetry traces will be exported by name.

- [Testing coverage report](./docs/testing-coverage.md)