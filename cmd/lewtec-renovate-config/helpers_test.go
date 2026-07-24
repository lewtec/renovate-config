package main

import (
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

func TestFindRenovateConfig(t *testing.T) {
	t.Run("missing returns empty", func(t *testing.T) {
		dir := t.TempDir()
		if got := findRenovateConfig(dir); got != "" {
			t.Fatalf("findRenovateConfig() = %q, want empty", got)
		}
	})

	t.Run("prefers first location in precedence order", func(t *testing.T) {
		dir := t.TempDir()
		// Create a later location first so we prove order, not discovery race.
		later := filepath.Join(dir, ".renovaterc")
		if err := os.WriteFile(later, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		preferred := filepath.Join(dir, "renovate.json")
		if err := os.WriteFile(preferred, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		got := findRenovateConfig(dir)
		if got != preferred {
			t.Fatalf("findRenovateConfig() = %q, want %q", got, preferred)
		}
	})

	t.Run("finds nested github path when root missing", func(t *testing.T) {
		dir := t.TempDir()
		nestedDir := filepath.Join(dir, ".github")
		if err := os.MkdirAll(nestedDir, 0o755); err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(nestedDir, "renovate.json")
		if err := os.WriteFile(want, []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		got := findRenovateConfig(dir)
		if got != want {
			t.Fatalf("findRenovateConfig() = %q, want %q", got, want)
		}
	})
}
