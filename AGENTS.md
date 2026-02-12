# AGENTS.md

Guidance for coding agents working in `github.com/tomski747/pvm`.

## Project overview

- **Language:** Go (module targets Go 1.21)
- **Type:** CLI application (Pulumi Version Manager)
- **Entrypoint:** `cmd/pvm/main.go`
- **Command framework:** `spf13/cobra`

## Repository map

- `cmd/pvm/main.go` - starts the CLI and calls `commands.Execute()`
- `internal/commands/` - Cobra commands (`install`, `use`, `list`, `current`, `remove`, `help`, `version`)
- `internal/utils/` - version resolution, GitHub release fetching, archive download/extraction, color output
- `internal/config/` - path/layout constants and runtime/test configuration
- `.github/workflows/ci.yml` - CI checks (lint, tests, coverage)
- `Makefile` - local developer commands

## Local development commands

Use these commands from repository root:

- `go mod download` - fetch dependencies
- `make build` - build `bin/pvm`
- `make test` or `go test -v ./...` - run tests
- `make lint` - run `golangci-lint` if installed
- `make coverage` - produce coverage artifacts in `bin/`

## Behavioral guardrails

- PVM data lives under `~/.pvm` by default (`internal/config/config.go`).
- Installed versions are stored in `~/.pvm/versions/<version>`.
- Active Pulumi binaries are symlinked into `~/.pvm/bin`.
- Version inputs can be exact or prefixes; resolution is handled by `utils.ResolveVersion`.
- Release fetches should use cache-aware helpers (`FetchGitHubReleases`, `GetAvailableVersions`) unless a direct latest lookup is explicitly required.
- Keep output style consistent with `internal/utils/color.go` helpers when touching command UX.

## Command change checklist

When adding or changing CLI behavior:

1. Update the relevant file in `internal/commands/`.
2. Ensure the command is registered on `rootCmd` (usually in `internal/commands/root.go` or command `init`).
3. Keep help text and examples consistent (`internal/commands/help.go`, command descriptions).
4. Update `README.md` usage examples when user-facing behavior changes.
5. Add or update unit tests for the new behavior.

## Testing guidance

- Prefer unit tests over integration/network-dependent tests.
- Use `config.SetTestConfig(...)` + `config.ResetConfig()` to isolate filesystem paths in tests.
- Use `httptest` servers for GitHub API behavior tests; avoid real API calls.
- Keep tests deterministic and independent of host machine state.

## Code style

- Run `gofmt` on changed Go files before committing.
- Keep functions focused and errors wrapped with context (`fmt.Errorf("...: %w", err)`).
- Preserve existing package boundaries (`commands` for CLI wiring, `utils` for core operations, `config` for environment/path concerns).
