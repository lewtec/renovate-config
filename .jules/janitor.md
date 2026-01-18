## 2026-01-14 - Fix outdated README.md
**Issue:** The README.md was severely outdated, describing a Python CLI that no longer exists in the repository. It included incorrect installation, prerequisite, and development instructions.
**Root Cause:** The CLI was migrated from Python to Go, but the documentation was not updated to reflect this significant change, leading to documentation drift.
**Solution:** I updated the README.md to accurately describe the current Go-based CLI. This included changing the installation instructions to use `go install`, updating prerequisites to Go, and revising the development section to use `mise` and Go commands.
**Pattern:** Documentation must be treated as code. When a core piece of functionality is rewritten or migrated, its corresponding documentation must be updated in the same changeset to prevent confusion and maintain usability.

## 2026-01-14 - Simplify extends normalization
**Issue:** The `addPresetToConfig` function contained nested if-else logic to handle the `extends` field, which can be either a string or a list. This increased cognitive load and cyclomatic complexity.
**Root Cause:** The logic for normalizing the `extends` field was inline, mixing data normalization with business logic.
**Solution:** Extracted the normalization logic into a pure helper function `normalizeExtends`. This simplifies the main function and makes the normalization logic independently testable.
**Pattern:** Extract pure data transformation logic into small, testable helper functions to keep the main business logic flow linear and clean.
