## 2026-01-31 - Optimize regex compilation in detectIndentation
**Issue:** The `detectIndentation` function recompiled the regex `\n(\s+)` on every call, which is inefficient and clutters the function logic.
**Root Cause:** The regex was defined locally inside the function scope instead of being a package-level constant.
**Solution:** Extracted the regex to a package-level variable `indentationRegex` using `regexp.MustCompile`. This ensures the regex is compiled only once at startup.
**Pattern:** Regex compilation (`regexp.MustCompile`) is expensive. Always define static regexes as package-level variables to avoid recompilation costs, especially in hot paths or frequently called functions.

## 2026-01-14 - Fix outdated README.md
**Issue:** The README.md was severely outdated, describing a Python CLI that no longer exists in the repository. It included incorrect installation, prerequisite, and development instructions.
**Root Cause:** The CLI was migrated from Python to Go, but the documentation was not updated to reflect this significant change, leading to documentation drift.
**Solution:** I updated the README.md to accurately describe the current Go-based CLI. This included changing the installation instructions to use `go install`, updating prerequisites to Go, and revising the development section to use `mise` and Go commands.
**Pattern:** Documentation must be treated as code. When a core piece of functionality is rewritten or migrated, its corresponding documentation must be updated in the same changeset to prevent confusion and maintain usability.
