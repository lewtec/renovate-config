## 2026-01-14 - Fix outdated README.md
**Issue:** The README.md was severely outdated, describing a Python CLI that no longer exists in the repository. It included incorrect installation, prerequisite, and development instructions.
**Root Cause:** The CLI was migrated from Python to Go, but the documentation was not updated to reflect this significant change, leading to documentation drift.
**Solution:** I updated the README.md to accurately describe the current Go-based CLI. This included changing the installation instructions to use `go install`, updating prerequisites to Go, and revising the development section to use `mise` and Go commands.
**Pattern:** Documentation must be treated as code. When a core piece of functionality is rewritten or migrated, its corresponding documentation must be updated in the same changeset to prevent confusion and maintain usability.

## 2026-01-18 - Extract logic from addPresetToConfig
**Issue:** The `addPresetToConfig` function was mixing I/O operations (file reading/writing, JSON marshaling) with business logic (modifying the `extends` array). This made it harder to test the logic in isolation and violated the Single Responsibility Principle.
**Root Cause:** Rapid development of the CLI tool led to "god functions" that handled everything from start to finish.
**Solution:** Extracted the core configuration manipulation logic into a pure helper function `updateConfigMap`. This function takes the configuration map and the preset reference, modifies the map if needed, and returns a boolean indicating if changes were made. `addPresetToConfig` now coordinates the I/O and delegates the logic to `updateConfigMap`.
**Pattern:** Separate "Business Logic" from "I/O Logic". Pure functions are easier to test and reason about than functions that perform side effects (like file system access).
