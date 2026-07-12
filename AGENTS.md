# Agent Guide

Use this guide when editing `github.com/agentstation/utc` with coding agents or AI assistants.

## Package Contract

- `utc.Time` is a wrapper around an unexported `time.Time`; it is not a type alias and not a full drop-in replacement for `time.Time`.
- Values entering through constructors, parsers, JSON/text/YAML unmarshaling, or SQL scanning must be normalized to UTC.
- Values leaving through `Time()`, `UTC()`, JSON/text/YAML marshaling, `String()`, and `Value()` must be UTC-normalized.
- Use `utc.New(time.Time)` or `utc.From(utc.UTC)` to create values from existing time-like values.
- Use `t.Time()` or `t.UTC()` when another library requires a concrete `time.Time`.

## Dependency Policy

- Keep the root module free of external dependencies.
- YAML codec integration belongs in `integration/yaml`, not in the root module.
- Keep assertion-only interface checks in `_test.go` files unless the imported package is part of a production method signature.
- `database/sql/driver` remains a production import because `Value() (driver.Value, error)` is the standard SQL value interface.

## Development Commands

Run these before committing behavior or docs changes:

```sh
go generate ./...
go test ./...
go test -race ./...
go test -tags=debug ./...
go vet ./...
golangci-lint run ./...
make test-yaml
(cd integration/yaml && go vet ./... && golangci-lint run ./... && go test -race ./...)
```

## Release Notes

- Keep README prose consistent with generated API docs.
- Use semver tags such as `v0.2.0`; before `v1.0.0`, exported API additions normally justify a minor bump.
- Do not describe the package as enforcing UTC by replacing every `time.Time` API. It enforces UTC at package boundaries while exposing `time.Time` for interoperability.
