package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestGhRepoListArgsExcludesForksServerSide(t *testing.T) {
	args := ghRepoListArgs("acme", false, 10000)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--source") {
		t.Fatalf("expected --source when includeForks=false, got %v", args)
	}
	if !strings.Contains(joined, "--limit 10000") {
		t.Fatalf("expected --limit 10000, got %v", args)
	}
	if !strings.Contains(joined, "--no-archived") {
		t.Fatalf("expected --no-archived, got %v", args)
	}
}

func TestGhRepoListArgsIncludesForksOmitsSource(t *testing.T) {
	args := ghRepoListArgs("acme", true, 50)
	for _, a := range args {
		if a == "--source" {
			t.Fatalf("did not expect --source when includeForks=true, got %v", args)
		}
	}
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--limit 50") {
		t.Fatalf("expected --limit 50, got %v", args)
	}
}

func TestFilterRepoNames(t *testing.T) {
	repos := []ghRepo{
		{Name: "app", IsFork: false},
		{Name: "forked", IsFork: true},
		{Name: "lib", IsFork: false},
	}

	got := filterRepoNames(repos, false)
	if len(got) != 2 || got[0] != "app" || got[1] != "lib" {
		t.Fatalf("includeForks=false: got %v, want [app lib]", got)
	}

	gotAll := filterRepoNames(repos, true)
	if len(gotAll) != 3 {
		t.Fatalf("includeForks=true: got %v, want 3 names", gotAll)
	}
}
