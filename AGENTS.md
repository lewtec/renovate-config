# Project Conventions

## Documentation

This project deviates from standard Go conventions by enforcing **JSDoc/TSDoc style** documentation comments.

- **Format:** Use `/** ... */` block comments for all exported functions, types, and variables.
- **Content:**
  - **Why > What:** Explain *why* the code exists and *how* it fits into the larger system.
  - **Non-Obvious:** Focus on edge cases, side effects, and permissions.
  - **No Tautologies:** Do not write "Returns user" for `getUser`. Explain *how* (e.g., "Retrieves cached user, fallback to DB").
- **Style:**
  - Block comments only (`/** ... */`).
  - Avoid line comments (`//`) unless inside function bodies for specific implementation details.

## Coding Standards

- **Logging:**
  - Use `log/slog`.
  - External command outputs must be logged at `Debug` level.
- **Structure:**
  - **Single Responsibility:** Isolate logic (e.g., SCM interactions) into dedicated packages or files.
  - **Pure Functions:** Extract business logic into pure functions to facilitate testing.
- **Environment:**
  - **Mise:** Use `mise` for all task execution and tool management.
  - **Dependencies:** Do not rely on globally installed tools; use `mise` to manage versions.

## Git & Workflow

- **Commit Messages:** Follow conventional commits (e.g., `feat:`, `fix:`, `docs:`, `chore:`).
