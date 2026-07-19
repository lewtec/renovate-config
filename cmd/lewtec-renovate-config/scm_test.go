package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetDefaultBranchFallsBackWhenOriginHEADMissing(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")
	// No origin remote / no origin/HEAD → recoverable fallback to main.
	if got := getDefaultBranch(dir); got != "main" {
		t.Fatalf("getDefaultBranch() = %q, want %q", got, "main")
	}
}

func TestGetDefaultBranchReadsOriginHEAD(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "test")
	// Detached default branch name for the local repo.
	runGit(t, dir, "checkout", "-b", "develop")
	if err := os.WriteFile(filepath.Join(dir, "README"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README")
	runGit(t, dir, "commit", "-m", "init")
	// Simulate a cloned remote-tracking default branch named develop.
	runGit(t, dir, "update-ref", "refs/remotes/origin/develop", "HEAD")
	runGit(t, dir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

	if got := getDefaultBranch(dir); got != "develop" {
		t.Fatalf("getDefaultBranch() = %q, want %q", got, "develop")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
