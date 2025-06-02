package ed

import "testing"

func TestReplaceRegion(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		markPos         int
		cursorPos       int
		replacement     string
		expectedContent string
		expectedPos     int
		expectError     bool
		errorMessage    string
	}{
		{
			name:            "replace forward selection",
			content:         "hello, world!",
			markPos:         5,
			cursorPos:       12,
			replacement:     " WORLD",
			expectedContent: "hello WORLD!",
			expectedPos:     11,
			expectError:     false,
		},
		{
			name:            "replace backward selection",
			content:         "hello, world!",
			markPos:         12,
			cursorPos:       5,
			replacement:     " WORLD",
			expectedContent: "hello WORLD!",
			expectedPos:     11,
			expectError:     false,
		},
		{
			name:            "replace with empty string",
			content:         "hello, world!",
			markPos:         5,
			cursorPos:       7,
			replacement:     "",
			expectedContent: "helloworld!",
			expectedPos:     5,
			expectError:     false,
		},
		{
			name:            "replace entire content",
			content:         "hello",
			markPos:         0,
			cursorPos:       5,
			replacement:     "world",
			expectedContent: "world",
			expectedPos:     5,
			expectError:     false,
		},
		{
			name:            "no mark set",
			content:         "hello, world!",
			markPos:         -1, // indicates no mark
			cursorPos:       5,
			replacement:     "test",
			expectedContent: "hello, world!",
			expectedPos:     5,
			expectError:     true,
			errorMessage:    "no mark set",
		},
		{
			name:            "mark and cursor at same position",
			content:         "hello, world!",
			markPos:         5,
			cursorPos:       5,
			replacement:     "X",
			expectedContent: "helloX, world!",
			expectedPos:     6,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New(tt.content)
			
			// Set mark if specified (markPos >= 0)
			if tt.markPos >= 0 {
				editor.setPosition(tt.markPos)
				editor.Do(SetMark())
			}
			
			// Set cursor position
			editor.setPosition(tt.cursorPos)
			
			err := editor.Do(ReplaceRegion(tt.replacement))
			
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
			
			if editor.String() != tt.expectedContent {
				t.Errorf("expected content %q, got %q", tt.expectedContent, editor.String())
			}
			
			if !tt.expectError && editor.Position() != tt.expectedPos {
				t.Errorf("expected position %d, got %d", tt.expectedPos, editor.Position())
			}
			
			// Mark should be cleared after successful replacement
			if !tt.expectError {
				_, marked := editor.Mark()
				if marked {
					t.Error("mark should be cleared after replacement")
				}
			}
		})
	}
}
