# Repository Guidelines

## Project Structure & Module Organization
- Core: `generate.go`, `types.go`, `inline_refs.go`, `genstate.go`.
- DSL and expressions: `dsl/`, `expr/`.
- Example: `examples/calc/` with generated code in `gen/` and binaries in `cmd/`.
- Tests and goldens: `*_test.go` and `testdata/*.json`.
- Docs: `README.md`, `INLINE_REFS.md`.

## Build, Test, and Development Commands
- `make all`: run gen, tests, lint, build examples, then clean (via `../plugins.mk`).
- `make lint`: `goimports` formatting check and `staticcheck` lint.
- `make test`: run `go test ./...`; verbose: `go test -v ./...`.
- `make gen`: regenerate example outputs (includes `examples/calc/gen/docs.json`).
- Update goldens: `go test ./... -- -update` (commit changes in `testdata/`).

## Coding Style & Naming Conventions
- Follow Goa’s CLAUDE.md layout: group declarations in this order per file — types, consts, vars, public funcs, public methods, private funcs, private methods. No section markers.
- Keep files focused and reasonably small; one main construct per file.
- Prefer `any` over `interface{}` in new code; exported identifiers use CamelCase; packages are short, lower-case.
- Never edit generated code in `examples/calc/gen/`; fix generators/templates instead.
- Every exported type, function, method, and interface must have a GoDoc-quality comment beginning with its name. Private declarations also get a short comment describing their role. Public struct fields need field comments.
- Avoid redundant defensive checks—only guard against nil/empty values at genuine system boundaries (e.g., external inputs). Inside our code paths assume invariants already validated.

## Curly Braces Rules
- Default: use multi-line braces for all code blocks (Go and Goa DSL).
- Exceptions only:
  - Empty DSL closures may be single-line, e.g., `JSONRPC(func() {})`.
  - Trivial methods returning a constant may be single-line, e.g., `func (e *Enum) String() string { return "foo" }`.
- Do not compress control flow. Preferred:
  
  ```go
  if err != nil {
      return err
  }
  ```
  
  Avoid: `if err != nil { return err }`.
- Place `else` on the same line as the closing `}` of the preceding block.

## Testing Guidelines
- Use `testing` plus `testify/assert` and `testify/require`.
- Golden tests compare generated output to files in `testdata/` (e.g., `no-payload-no-return.json`).
- When behavior changes intentionally, regenerate goldens and review diffs carefully.

## Commit & Pull Request Guidelines
- Conventional commits: `feat(docs): ...`, `fix(docs): ...`, `chore: ...`.
- PRs include description, rationale, linked issues, and testing notes; keep changes small and scoped.
- If generation output changes, run `make gen` and commit relevant updates under `examples/calc/gen/` and `testdata/`.
