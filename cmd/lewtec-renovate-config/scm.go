package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// ghRepo represents a repository in the GitHub CLI JSON output.
type ghRepo struct {
	Name   string `json:"name"`
	IsFork bool   `json:"isFork"`
}

// runCommand executes a shell command in a specified directory.
//
// This wrapper handles:
// - Setting the working directory.
// - Logging the command and its output at Debug level (to avoid clutter).
// - Capturing stdout separately from stderr to avoid contamination of structured output.
func runCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	slog.Debug("Running command", "command", name, "args", args)
	output, err := cmd.Output()
	outStr := string(output)
	if outStr != "" {
		slog.Debug("Command output", "output", outStr)
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outStr, fmt.Errorf("%w (stderr: %s)", err, string(exitErr.Stderr))
		}
		return outStr, err
	}
	return outStr, nil
}

// checkGhCli verifies that the GitHub CLI (`gh`) is installed and authenticated.
//
// Returns true when available and authenticated. On failure, err is a sentinel
// (ErrGhCLINotFound or ErrGhCLINotAuth) suitable for ReportError / slog.
func checkGhCli() (ok bool, err error) {
	if _, err := runCommand("", "gh", "--version"); err != nil {
		return false, ErrGhCLINotFound
	}
	if _, err := runCommand("", "gh", "auth", "status"); err != nil {
		return false, ErrGhCLINotAuth
	}
	return true, nil
}

// cloneRepository clones a repository to a destination path.
//
// Attempts to use `gh repo clone` first (for better auth handling),
// and falls back to standard `git clone` (HTTPS) if `gh` fails.
//
// Returns true if cloning succeeded, false otherwise.
func cloneRepository(owner, repo, destination string) bool {
	_, err := runCommand("", "gh", "repo", "clone", fmt.Sprintf("%s/%s", owner, repo), destination)
	if err != nil {
		slog.Warn("Error cloning with gh", "error", err)
		slog.Info("Fallback to git clone")
		_, err = runCommand("", "git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", owner, repo), destination)
		if err != nil {
			ReportError(err, "Error cloning with git")
			return false
		}
	}
	return true
}

// getDefaultBranch identifies the default branch of the repository (e.g., main, master).
//
// Queries the git remote `origin/HEAD` to find the default branch.
// Falls back to "main" if the query fails or returns nothing. Missing
// origin/HEAD is common after some clones, so that path is recoverable and
// must not go through ReportError (reserved for unexpected failures).
func getDefaultBranch(repoPath string) string {
	output, err := runCommand(repoPath, "git", "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		slog.Warn("Could not resolve origin/HEAD; defaulting to main", "error", err)
		return "main"
	}
	branch := strings.TrimSpace(output)
	branch = strings.TrimPrefix(branch, "refs/remotes/origin/")
	if branch == "" {
		return "main"
	}
	return branch
}

// createGitHubPR creates a Pull Request using the GitHub CLI.
//
// If PR creation fails (e.g., branch already exists, permission issues),
// it logs the error but does not crash, allowing the user to create it manually.
func createGitHubPR(repoPath, owner, repo, baseBranch, headBranch, title, body string) bool {
	output, err := runCommand(repoPath, "gh", "pr", "create",
		"--repo", fmt.Sprintf("%s/%s", owner, repo),
		"--title", title,
		"--body", body,
		"--base", baseBranch,
		"--head", headBranch,
	)
	if err != nil {
		ReportError(err, "Error creating PR", "output", output)
		slog.Info("Branch pushed. Create PR manually", "branch", headBranch)
		return false
	}
	slog.Info("Pull request created successfully", "output", output)
	return true
}

// getRepositories fetches a list of repository names for a given owner (user or org).
//
// Uses `gh repo list` to retrieve repositories.
// - Filters out archived repositories automatically (via `gh` flag).
// - Optionally includes forks based on the `includeForks` parameter.
// - Returns a list of repository names (without owner prefix).
func getRepositories(owner string, includeForks bool) ([]string, error) {
	slog.Info("Fetching repositories", "owner", owner)
	output, err := runCommand("", "gh", "repo", "list", owner, "--json", "name,isFork", "--limit", "1000", "--no-archived")
	if err != nil {
		return nil, err
	}

	var repos []ghRepo
	if err := json.Unmarshal([]byte(output), &repos); err != nil {
		return nil, err
	}

	var names []string
	for _, r := range repos {
		if !includeForks && r.IsFork {
			continue
		}
		names = append(names, r.Name)
	}
	return names, nil
}
