package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectIndentation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "default when no indentation",
			content: "{\n\"extends\": []\n}\n",
			want:    "  ",
		},
		{
			name:    "two spaces",
			content: "{\n  \"extends\": []\n}\n",
			want:    "  ",
		},
		{
			name:    "four spaces",
			content: "{\n    \"extends\": []\n}\n",
			want:    "    ",
		},
		{
			name:    "tabs",
			content: "{\n\t\"extends\": []\n}\n",
			want:    "\t",
		},
		{
			name:    "empty content",
			content: "",
			want:    "  ",
		},
		{
			name:    "single line no break",
			content: `{"extends":[]}`,
			want:    "  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectIndentation(tt.content)
			if got != tt.want {
				t.Fatalf("detectIndentation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func writeEmptyConfig(t *testing.T, dir, rel string) string {
	t.Helper()
	path := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestEnsureRenovateConfig(t *testing.T) {
	t.Run("creates root renovate.json when missing", func(t *testing.T) {
		dir := t.TempDir()
		got, err := ensureRenovateConfig(dir)
		if err != nil {
			t.Fatalf("ensureRenovateConfig() error = %v", err)
		}
		want := filepath.Join(dir, "renovate.json")
		if got != want {
			t.Fatalf("ensureRenovateConfig() = %q, want %q", got, want)
		}
		content, err := os.ReadFile(got)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "{}\n" {
			t.Fatalf("created file content = %q, want %q", content, "{}\n")
		}
	})

	t.Run("returns existing config without rewriting it", func(t *testing.T) {
		dir := t.TempDir()
		existing := writeEmptyConfig(t, dir, ".github/renovate.json")
		if err := os.WriteFile(existing, []byte("{\n  \"extends\": []\n}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := ensureRenovateConfig(dir)
		if err != nil {
			t.Fatalf("ensureRenovateConfig() error = %v", err)
		}
		if got != existing {
			t.Fatalf("ensureRenovateConfig() = %q, want %q", got, existing)
		}
		if _, err := os.Stat(filepath.Join(dir, "renovate.json")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("did not expect a new root renovate.json, stat err = %v", err)
		}
		content, err := os.ReadFile(existing)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != "{\n  \"extends\": []\n}\n" {
			t.Fatalf("existing file was rewritten: %q", content)
		}
	})
}

func TestFindRenovateConfig(t *testing.T) {
	if len(renovateConfigLocations) == 0 {
		t.Fatal("renovateConfigLocations is empty")
	}

	t.Run("missing returns empty", func(t *testing.T) {
		dir := t.TempDir()
		if got := findRenovateConfig(dir); got != "" {
			t.Fatalf("findRenovateConfig() = %q, want empty", got)
		}
	})

	// Drive cases from renovateConfigLocations so the search list is not copied here.
	for _, loc := range renovateConfigLocations {
		t.Run("finds "+loc, func(t *testing.T) {
			dir := t.TempDir()
			writeEmptyConfig(t, dir, loc)
			want := filepath.Join(dir, loc)
			got := findRenovateConfig(dir)
			if got != want {
				t.Fatalf("findRenovateConfig() = %q, want %q", got, want)
			}
		})
	}

	t.Run("prefers first location when all exist", func(t *testing.T) {
		dir := t.TempDir()
		// Create later paths first so we prove precedence, not discovery race.
		for i := len(renovateConfigLocations) - 1; i >= 0; i-- {
			writeEmptyConfig(t, dir, renovateConfigLocations[i])
		}
		want := filepath.Join(dir, renovateConfigLocations[0])
		got := findRenovateConfig(dir)
		if got != want {
			t.Fatalf("findRenovateConfig() = %q, want %q", got, want)
		}
	})
}
