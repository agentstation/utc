# Contributing

Thanks for improving `github.com/agentstation/utc`.

## Development

The root module targets Go 1.18+ and intentionally has no external dependencies. Optional codec integration tests live in nested modules.

Before opening a pull request, run:

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

## Change Guidelines

- Preserve the UTC invariant at constructors, parsers, scanners, and serialization boundaries.
- Keep the wrapped `time.Time` unexported.
- Use `t.Time()` or `t.UTC()` for interop with libraries that require concrete `time.Time`.
- Keep optional third-party integrations out of the root module.
- Regenerate README API docs with `go generate ./...` when exported docs change.

## Releases

Releases use semver tags. Before `v1.0.0`, exported API additions usually use a minor version bump, while docs-only or internal cleanup usually does not require a release.
