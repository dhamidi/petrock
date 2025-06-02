package main

import (
	"testing"

	"github.com/dhamidi/petrock/internal/generator"
)

func TestParseFieldDefinition(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    generator.CommandField
		expectError bool
	}{
		{
			name:  "simple string field",
			input: "postID:string",
			expected: generator.CommandField{
				Name: "postID",
				Type: "string",
			},
			expectError: false,
		},
		{
			name:  "time field",
			input: "publishAt:time.Time",
			expected: generator.CommandField{
				Name: "publishAt",
				Type: "time.Time",
			},
			expectError: false,
		},
		{
			name:        "missing colon",
			input:       "postIDstring",
			expectError: true,
		},
		{
			name:        "empty name",
			input:       ":string",
			expectError: true,
		},
		{
			name:        "empty type",
			input:       "postID:",
			expectError: true,
		},
		{
			name:        "invalid identifier",
			input:       "123invalid:string",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFieldDefinition(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result.Name != tt.expected.Name {
				t.Errorf("expected name %q, got %q", tt.expected.Name, result.Name)
			}
			
			if result.Type != tt.expected.Type {
				t.Errorf("expected type %q, got %q", tt.expected.Type, result.Type)
			}
		})
	}
}

func TestIsValidGoIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple identifier", "postID", true},
		{"underscore start", "_private", true},
		{"with digits", "field123", true},
		{"starts with digit", "123field", false},
		{"empty", "", false},
		{"with space", "post ID", false},
		{"with dash", "post-id", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidGoIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, result)
			}
		})
	}
}
