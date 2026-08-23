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
	t.Run("missing returns empty", func(t *testing.T) {
		dir := t.TempDir()
		if got := findRenovateConfig(dir); got != "" {
			t.Fatalf("findRenovateConfig() = %q, want empty", got)
		}
	})

	tests := []struct {
		name  string
		files []string // relative paths to create; first entry is the expected match
	}{
		{
			// Create a later location first so we prove order, not discovery race.
			name:  "prefers first location in precedence order",
			files: []string{"renovate.json", ".renovaterc"},
		},
		{
			name:  "finds nested github path when root missing",
			files: []string{".github/renovate.json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			// Create non-preferred paths first so discovery order is the subject under test.
			for i := len(tt.files) - 1; i >= 0; i-- {
				writeEmptyConfig(t, dir, tt.files[i])
			}
			want := filepath.Join(dir, tt.files[0])
			got := findRenovateConfig(dir)
			if got != want {
				t.Fatalf("findRenovateConfig() = %q, want %q", got, want)
			}
		})
	}
}
