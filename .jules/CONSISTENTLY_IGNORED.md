# Consistently Ignored Changes

This file lists patterns of changes that have been consistently rejected by human reviewers. All agents MUST consult this file before proposing a new change. If a planned change matches any pattern described below, it MUST be abandoned.

---

## IGNORE: Input Validation for Repo/Owner

**- Pattern:** Adding strict regex validation (e.g. `^[a-zA-Z0-9_.-]+$`) or path traversal checks for `owner` and `repo` CLI arguments.
**- Justification:** Multiple security PRs attempting to patch path traversal in these arguments have been rejected (#17, #22, #23, #25). This suggests the validation is either too restrictive or the risk is accepted/handled elsewhere.
**- Files Affected:** `cmd/lewtec-renovate-config/main.go`

---

## IGNORE: Extracting Logic from addPresetToConfig

**- Pattern:** Refactoring `addPresetToConfig` to use helper functions like `updateConfigMap` or `normalizeExtends` to separate business logic from I/O.
**- Justification:** PRs (#18, #21, #24) proposing these specific refactorings were rejected. The current inline implementation in `main.go` is preferred.
**- Files Affected:** `cmd/lewtec-renovate-config/main.go`

---

## IGNORE: Adding Autorelease/CI Workflows

**- Pattern:** Adding `.github/workflows/autorelease.yml`, `.golangci.yml`, `dprint.json` or related `mise` tasks for linting/formatting/releasing.
**- Justification:** PRs (#19, #27, #32) setting up this specific CI/CD and linter infrastructure have been consistently closed.
**- Files Affected:** `.golangci.yml`, `dprint.json`, `.github/workflows/*`

---

## IGNORE: Verbose JSDoc for Simple Identifiers

**- Pattern:** Using JSDoc-style block comments (`/** ... */`) for variables, constants, or simple functions where standard Go comments (`//`) suffice, or adding redundant documentation (e.g., repeating the variable name).
**- Justification:** PRs #28 and #33 adding extensive JSDoc comments to `main.go` were rejected. The project prefers an "Essentialist" documentation style, reserving detailed block comments for complex logic or public APIs, and avoiding noise for internal CLI variables.
**- Files Affected:** `cmd/lewtec-renovate-config/*.go`
