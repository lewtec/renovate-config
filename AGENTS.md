# Project Conventions

## Error Handling

- **No Ignored Errors:** Errors must never be ignored. `_ = func()` or empty `if err != nil {}` blocks are strictly forbidden.
- **Centralized Error Reporting:**
  - All unexpected errors must be reported via a centralized function (e.g., `reportError` in `errors.go`).
  - Do not use `fmt.Println` or `log.Printf` or `slog.Error` directly for unexpected errors. Use the centralized handler.
  - The centralized handler should be capable of logging context and potentially sending errors to an external service (like Sentry) in the future without changing call sites.

## Environment

- **Mise:** Use `mise` for all task executions (`mise run ...`).
