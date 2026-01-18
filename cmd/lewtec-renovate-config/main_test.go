package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidInput(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"valid-name", true},
		{"valid_name", true},
		{"valid.name", true},
		{"ValidName123", true},
		{".", false},
		{"..", false},
		{"../foo", false},
		{"foo/bar", false},
		{"foo;bar", false},
		{"-flag", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isValidInput(tt.input); got != tt.expected {
				t.Errorf("isValidInput(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAddPresetToConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "renovate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name            string
		initialContent  string
		presetRef       string
		expectChange    bool
		expectedExtends []string
	}{
		{
			name:            "Empty config",
			initialContent:  `{}`,
			presetRef:       "foo",
			expectChange:    true,
			expectedExtends: []string{"foo"},
		},
		{
			name:            "Existing extends string",
			initialContent:  `{"extends": "bar"}`,
			presetRef:       "foo",
			expectChange:    true,
			expectedExtends: []string{"bar", "foo"},
		},
		{
			name:            "Existing extends list",
			initialContent:  `{"extends": ["bar"]}`,
			presetRef:       "foo",
			expectChange:    true,
			expectedExtends: []string{"bar", "foo"},
		},
		{
			name:            "Preset already at end",
			initialContent:  `{"extends": ["bar", "foo"]}`,
			presetRef:       "foo",
			expectChange:    false,
			expectedExtends: []string{"bar", "foo"},
		},
		{
			name:            "Preset in middle",
			initialContent:  `{"extends": ["foo", "bar"]}`,
			presetRef:       "foo",
			expectChange:    true,
			expectedExtends: []string{"bar", "foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tmpDir, "renovate.json")
			err := os.WriteFile(configPath, []byte(tt.initialContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			changed, err := addPresetToConfig(configPath, tt.presetRef)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if changed != tt.expectChange {
				t.Errorf("expected change %v, got %v", tt.expectChange, changed)
			}

			// Verify content
			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatal(err)
			}

			var config map[string]interface{}
			if err := json.Unmarshal(content, &config); err != nil {
				t.Fatal(err)
			}

			extendsRaw := config["extends"]
			var extends []string
			if list, ok := extendsRaw.([]interface{}); ok {
				for _, v := range list {
					extends = append(extends, v.(string))
				}
			}

			if len(extends) != len(tt.expectedExtends) {
				t.Errorf("expected extends length %d, got %d", len(tt.expectedExtends), len(extends))
			}

			for i, v := range tt.expectedExtends {
				if i < len(extends) && extends[i] != v {
					t.Errorf("expected extends[%d] to be %s, got %s", i, v, extends[i])
				}
			}
		})
	}
}
