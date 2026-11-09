# strictgo

An intentionally opinionated Go linter for keeping Go code readable and consistent in the age of AI-driven development. It enforces strict, idiomatic Go conventions that would be impractical to demand from human developers but are inexpensive for coding agents to follow.

Its analyzers focus on simpler control flow, smaller variable scopes, clearer API boundaries, predictable error handling, and idiomatic naming and initialization.

## Install

```sh
go install github.com/konradreiche/strictgo/cmd/strictgo@latest
```

Ensure your Go bin directory (`$(go env GOPATH)/bin` by default) is on `PATH`.

## Usage

Run from your Go module:

```sh
strictgo ./...
```
