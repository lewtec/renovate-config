## 2026-01-14 - Fix outdated README.md
**Issue:** The README.md was severely outdated, describing a Python CLI that no longer exists in the repository. It included incorrect installation, prerequisite, and development instructions.
**Root Cause:** The CLI was migrated from Python to Go, but the documentation was not updated to reflect this significant change, leading to documentation drift.
**Solution:** I updated the README.md to accurately describe the current Go-based CLI. This included changing the installation instructions to use `go install`, updating prerequisites to Go, and revising the development section to use `mise` and Go commands.
**Pattern:** Documentation must be treated as code. When a core piece of functionality is rewritten or migrated, its corresponding documentation must be updated in the same changeset to prevent confusion and maintain usability.

## 2026-01-16 - Refactor addPresetToConfig logic
**Issue:** The `addPresetToConfig` function in `main.go` was mixing file I/O, JSON marshalling, and business logic (modifying the `extends` slice), making it complex and harder to test.
**Root Cause:** Organic growth of the function handling multiple concerns.
**Solution:** Extracted the configuration update logic into a pure function `updateConfigMap`. This separates the "what" (business logic) from the "how" (persistence).
**Pattern:** Isolate impure I/O from pure business logic. This improves testability and readability, adhering to the Single Responsibility Principle.
