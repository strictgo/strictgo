# strictgo

An intentionally opinionated Go linter for keeping Go code readable and consistent in the age of AI-driven development. It enforces strict, idiomatic Go conventions that would be impractical to demand from human developers but are inexpensive for coding agents to follow.

Its analyzers focus on simpler control flow, smaller variable scopes, clearer API boundaries, predictable error handling, and idiomatic naming and initialization.

## Install

```sh
go install github.com/strictgo/strictgo/cmd/strictgo@latest
```

Ensure your Go bin directory (`$(go env GOPATH)/bin` by default) is on `PATH`.

## Usage

Run from your Go module:

```sh
strictgo ./...
```

Apply available fixes automatically:

```sh
strictgo -fix ./...
```

## Analyzers

| Analyzer | What it checks |
|---|---|
| `badname` | Flags error variables not named `err`, single-letter names outside short scopes, redundant `Str`/`Map` suffixes, and generic package or file names using `util`, `utils`, `helper`, `helpers`, or `misc`. |
| `blankparam` | Flags explicit `_` parameters in functions, methods, literals, interfaces, and function types. Suggests collision-safe names, including `ctx` for `context.Context`. |
| `boolswitch` | Prefers `if` statements and early returns over boolean or expressionless `switch` statements. |
| `declorder` | Enforces a consistent order for top-level declarations, including types, constructors, methods, and functions. |
| `errcheck` | Requires comparisons in the form `err != nil`. Despite the name, it does not check for ignored errors. |
| `forcederef` | Flags explicit pointer dereferences used to satisfy function or method signatures within the same module. |
| `forint` | Flags traditional `for` loops with integer counters. |
| `importname` | Requires lowercase import names without underscores and flags redundant aliases. Offers collision-safe fixes; `_` and `.` imports are flagged without fixes, while cgo's `C` import is exempt. |
| `keepflat` | Flags nesting deeper than five blocks within a function. |
| `limitscope` | Flags opportunities to narrow variable scope, such as moving an `err` assignment into the following `if` initializer. |
| `nilresult` | For functions returning a pointer or struct alongside an error, expects a nil pointer or zero-value struct on error. |
| `noelse` | Flags unnecessary `else` branches after returning branches, with scope-safe fixes. |
| `okcheck` | Flags discarding the boolean result of a `(value, bool)` function call or comma-ok type assertion with `_`. |
| `returnstructs` | Flags functions and methods that return interfaces rather than concrete types, with exceptions such as `error` and `any`. |
| `singleassign` | Prefers separate assignments for separate right-hand-side expressions, while allowing swaps and rotations. |
| `zeroinit` | Flags redundant zero-value initializers for basic types, such as `count := 0` where `var count int` would do. |
