package main

import "testing"

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid alphanumeric", "username", false},
		{"Valid with hyphen", "my-repo", false},
		{"Valid with underscore", "my_repo", false},
		{"Valid with dot", "my.repo", false},
		{"Invalid with slash", "owner/repo", true},
		{"Invalid traversal", "../etc", true},
		{"Invalid special chars", "repo;", true},
		{"Invalid space", "my repo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInput(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
