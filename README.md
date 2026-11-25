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

Files marked `// Code generated … DO NOT EDIT.` receive no diagnostics or fixes, but remain available as analysis context.

## Analyzers

| Analyzer | What it checks |
|---|---|
| `badname` | Flags local error variables not named `err`, `testing.TB` variables named `t` instead of `tb`, single-letter names outside short scopes, redundant `Str`/`Map` suffixes, and generic package or file names using `util`, `utils`, `helper`, `helpers`, or `misc`. |
| `blankparam` | Flags explicit `_` parameters in functions, methods, literals, interfaces, and function types. Suggests collision-safe names, including `ctx` for `context.Context`. |
| `boolswitch` | Prefers `if` statements and early returns over boolean or expressionless `switch` statements. |
| `declnear` | Flags unexported, method-free structs separated from their first explicit function reference by unrelated functions. Initially skips grouped, shared, and test-file types. |
| `declorder` | Enforces ordering within top-level type blocks and places test entrypoints before helper types and functions in `_test.go` files. |
| `errcheck` | Requires `err != nil` comparisons and flags explicit nil error returns inside non-nil error branches. Does not check for ignored error results. |
| `forcederef` | Flags explicit pointer dereferences used to satisfy function or method signatures within the same module. |
| `forint` | Flags traditional `for` loops with integer counters. |
| `importname` | Requires lowercase import names without underscores and flags redundant aliases. Offers collision-safe fixes; `_` and `.` imports are flagged without fixes, while cgo's `C` import is exempt. |
| `keepflat` | Flags nesting deeper than five blocks within a function. |
| `limitscope` | Flags opportunities to narrow variable scope, such as moving an `err` assignment into the following `if` initializer. |
| `nilreceiver` | Flags direct nil comparisons and nil switch cases on method receivers. |
| `nilresult` | For functions returning a pointer or struct alongside an error, expects a nil pointer or zero-value struct on error. |
| `noelse` | Flags unnecessary `else` branches after returning branches, with scope-safe fixes. |
| `okcheck` | Flags discarding the boolean result of a `(value, bool)` function call or comma-ok type assertion with `_`. |
| `returnstructs` | Flags functions and methods that return interfaces rather than concrete types, with exceptions such as `error` and `any`. |
| `scopesprawl` | Flags files with too many concepts, grouping types/methods, constructors, functional-option families, and exclusively owned private helpers/data/state. Shared private constants do not add concepts. Reports group names; defaults to 5 (`-scopesprawl.limit=N`). Test and generated files are excluded from reporting. |
| `singleassign` | Prefers separate assignments for separate right-hand-side expressions, while allowing swaps and rotations. |
| `testingtb` | Recommends `testing.TB` for directly called helpers that do not need `*testing.T`. Preserves test entrypoints, callbacks, and T-specific uses. Diagnostic-only. |
| `zeroinit` | Flags redundant zero-value initializers for basic types, such as `count := 0` where `var count int` would do. |
