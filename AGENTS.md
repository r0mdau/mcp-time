# MCP Time Agents Guide

## Project Structure
The codebase follows standard Go project layout:
- `cmd/mcp-time/main.go` - Application entry point (flags, server setup)
- `internal/types/` - Shared type definitions (TimeResult, ConvertTimeInput, etc.)
- `internal/handlers/` - MCP tool handlers (GetCurrentTime, ConvertTime, RegisterTools)
- `internal/timezone/` - Low-level timezone operations (GetNowInLocation, IsDST, FormatISOSeconds)
- `internal/timeutil/` - Time utilities (BuildTimeResult, ParseTimeInput, FormatTimeDifference)

## Dev environment tips
- Use `make build` to format (`go fmt`), vet, and compile the server into `build/mcp-time` in one step.
- Command-line flags: `--local-timezone` overrides the detected IANA zone, `--port` sets the HTTP listener (default `8080`). Run `./build/mcp-time --help` for details.
- Code organization: timezone operations in `internal/timezone/`, time utilities in `internal/timeutil/`, MCP handlers in `internal/handlers/`. Each package has its own focused test file.
- Coverage artifacts are written to `coverage.out`. Open `go tool cover -html=coverage.out` for an interactive report when researching regressions.
- The HTTP server uses the MCP Go SDK. See `RegisterTools` in `internal/handlers/handlers.go` for tool schemas and update both schemas and helper logic together.

## Testing instructions
- Primary suite: `make test` (runs `go test -cover ./...`). Current coverage: handlers 96.8%, timeutil 100%, timezone 95.1%. Keep overall coverage above 80%.
- For targeted runs, use `go test -v ./...` or scope to a package:
  - `go test -v ./internal/handlers` - Test MCP tool handlers
  - `go test -v ./internal/timezone` - Test timezone operations
  - `go test -v ./internal/timeutil` - Test time utilities
- Benchmarks live in `internal/timezone/timezone_test.go` and `internal/timeutil/timeutil_test.go`; run with `go test -bench=. -benchmem ./internal/...` when tuning performance.
- After modifying functions, regenerate coverage with `go test -coverprofile=coverage.out ./...` and inspect via `go tool cover -func=coverage.out` or `go tool cover -html=coverage.out`.
- Each package has table-driven tests; add new test cases alongside changes to keep coverage and edge-case checks intact.

## PR instructions
- Title format: `<summary>`.
- Before pushing, run `make build` (formats + vets + compiles) and `make test` (full coverage suite).
- Include tests for new behavior. If coverage dips below 80%, add or update tests before requesting review.
