package main

import "testing"

func TestValidateInput(t *testing.T) {
	tests := []struct {
		input    string
		isValid  bool
	}{
		{"valid-repo", true},
		{"valid_repo", true},
		{"valid.repo", true},
		{"validRepo123", true},
		{"../invalid", false},
		{"invalid/path", false},
		{".", false},
		{"..", false},
		{"-invalid", true}, // Allowed by regex, though maybe bad practice for starting char, but safe from traversal
		{"", false},
		{"invalid\\path", false},
	}

	for _, tt := range tests {
		err := validateInput(tt.input)
		if tt.isValid && err != nil {
			t.Errorf("validateInput(%q) returned error: %v", tt.input, err)
		}
		if !tt.isValid && err == nil {
			t.Errorf("validateInput(%q) expected error, got nil", tt.input)
		}
	}
}
