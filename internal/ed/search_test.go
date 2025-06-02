package ed

import (
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		initialPos     int
		searchText     string
		expectedPos    int
		expectError    bool
		errorMessage   string
	}{
		{
			name:        "find text from beginning",
			content:     "hello, world!",
			initialPos:  0,
			searchText:  ",",
			expectedPos: 5,
			expectError: false,
		},
		{
			name:        "find text from middle",
			content:     "hello, world!",
			initialPos:  3,
			searchText:  "world",
			expectedPos: 7,
			expectError: false,
		},
		{
			name:        "text not found",
			content:     "hello, world!",
			initialPos:  0,
			searchText:  "xyz",
			expectedPos: 0, // position should not change
			expectError: true,
			errorMessage: "text not found: xyz",
		},
		{
			name:        "empty search text",
			content:     "hello, world!",
			initialPos:  0,
			searchText:  "",
			expectedPos: 0,
			expectError: true,
			errorMessage: "search text cannot be empty",
		},
		{
			name:        "search past current position",
			content:     "hello, world!",
			initialPos:  8,
			searchText:  "hello",
			expectedPos: 8, // should not find text before current position
			expectError: true,
			errorMessage: "text not found: hello",
		},
		{
			name:        "find text at current position",
			content:     "hello, world!",
			initialPos:  7,
			searchText:  "world",
			expectedPos: 7,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New(tt.content)
			editor.setPosition(tt.initialPos)
			
			err := editor.Do(Search(tt.searchText))
			
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if edErr, ok := err.(Error); ok {
					if edErr.Message != tt.errorMessage {
						t.Errorf("expected error message %q, got %q", tt.errorMessage, edErr.Message)
					}
				} else {
					t.Errorf("expected ed.Error, got %T", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			
			if editor.Position() != tt.expectedPos {
				t.Errorf("expected position %d, got %d", tt.expectedPos, editor.Position())
			}
		})
	}
}
