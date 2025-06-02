package ed

import "testing"

func TestBeginningOfBuffer(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		initialPos     int
		expectedPos    int
	}{
		{
			name:        "from middle to beginning",
			content:     "hello world",
			initialPos:  5,
			expectedPos: 0,
		},
		{
			name:        "from end to beginning",
			content:     "hello world",
			initialPos:  11,
			expectedPos: 0,
		},
		{
			name:        "already at beginning",
			content:     "hello world",
			initialPos:  0,
			expectedPos: 0,
		},
		{
			name:        "empty buffer",
			content:     "",
			initialPos:  0,
			expectedPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New(tt.content)
			editor.setPosition(tt.initialPos)
			
			err := editor.Do(BeginningOfBuffer())
			if err != nil {
				t.Fatalf("BeginningOfBuffer() returned error: %v", err)
			}
			
			if editor.Position() != tt.expectedPos {
				t.Errorf("expected position %d, got %d", tt.expectedPos, editor.Position())
			}
		})
	}
}
