# Assert embeddings

This package exists to drop the dependency on stretchr/testify/assert.
It mutes the conversation on which to use (assert or require), and
provides a limited API that does about the same.

It's a copy of `internal/assert` in github.com/titpetric/platform, which
is itself based on [fortio.org/assert](https://github.com/fortio/assert)
and contains the following modifications:

- Use `testing.TB` instead of explicit `*testing.T`
- Use `any` or `...any` for printf arguments

On top of the platform copy this one adds:

- `GreaterOrEqual` in `assert_extra.go`, next to `Greater`

Nothing here calls `FailNow`. A failed assertion reports and the test
continues, so a call site that dereferences a value the assertion was
guarding will panic rather than stop.

It's meant to be used from inside the current module and not
intended as an imported dependency.
