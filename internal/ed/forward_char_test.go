package ed

import "testing"

func TestForwardChar(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		initialPos  int
		moveBy      int
		expectedPos int
	}{
		{
			name:        "move forward by 1",
			content:     "hello world",
			initialPos:  0,
			moveBy:      1,
			expectedPos: 1,
		},
		{
			name:        "move forward by 5",
			content:     "hello world",
			initialPos:  0,
			moveBy:      5,
			expectedPos: 5,
		},
		{
			name:        "move past end",
			content:     "hello",
			initialPos:  3,
			moveBy:      10,
			expectedPos: 5, // clamped to end
		},
		{
			name:        "move backward with negative",
			content:     "hello world",
			initialPos:  5,
			moveBy:      -2,
			expectedPos: 3,
		},
		{
			name:        "move before beginning with negative",
			content:     "hello world",
			initialPos:  2,
			moveBy:      -5,
			expectedPos: 0, // clamped to beginning
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New(tt.content)
			editor.setPosition(tt.initialPos)
			
			err := editor.Do(ForwardChar(tt.moveBy))
			if err != nil {
				t.Fatalf("ForwardChar() returned error: %v", err)
			}
			
			if editor.Position() != tt.expectedPos {
				t.Errorf("expected position %d, got %d", tt.expectedPos, editor.Position())
			}
		})
	}
}
