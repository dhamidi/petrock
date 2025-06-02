package ed

import "testing"

func TestSetMark(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		position    int
		expectedMark int
		expectedMarked bool
	}{
		{
			name:          "set mark at beginning",
			content:       "hello world",
			position:      0,
			expectedMark:  0,
			expectedMarked: true,
		},
		{
			name:          "set mark in middle",
			content:       "hello world",
			position:      5,
			expectedMark:  5,
			expectedMarked: true,
		},
		{
			name:          "set mark at end",
			content:       "hello world",
			position:      11,
			expectedMark:  11,
			expectedMarked: true,
		},
		{
			name:          "set mark in empty buffer",
			content:       "",
			position:      0,
			expectedMark:  0,
			expectedMarked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New(tt.content)
			editor.setPosition(tt.position)
			
			err := editor.Do(SetMark())
			if err != nil {
				t.Fatalf("SetMark() returned error: %v", err)
			}
			
			mark, marked := editor.Mark()
			if mark != tt.expectedMark {
				t.Errorf("expected mark %d, got %d", tt.expectedMark, mark)
			}
			if marked != tt.expectedMarked {
				t.Errorf("expected marked %t, got %t", tt.expectedMarked, marked)
			}
		})
	}
}
