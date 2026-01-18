# Sentinel Journal

## 2026-01-18 - Path Traversal in Repository Arguments
**Vulnerability:** The CLI accepts `owner` and `repo` arguments without validation and uses them in `filepath.Join` to construct the local clone path. This allows path traversal attacks (e.g., using `../` in the repo name to write to arbitrary directories).
**Learning:** Even in CLI tools where the user attacks themselves, input validation is crucial for defense-in-depth and to prevent accidental misuse or potentially malicious scripts invoking the tool.
**Prevention:** Always validate user inputs against an allowlist (e.g., alphanumeric regex) before using them in file system operations or command execution.
