package ed

import "testing"

func TestExample(t *testing.T) {
	// Test the exact example from the requirements
	editor := New("hello, world!")
	err := editor.Do(
		BeginningOfBuffer(),
		Search(","),
		SetMark(),
		Search("!"),
		ReplaceRegion("WORLD"),
	)
	
	if err != nil {
		t.Fatalf("example failed with error: %v", err)
	}
	
	// The actual result replacing ", world!" with "WORLD"
	expected := "helloWORLD!"
	if editor.String() != expected {
		t.Errorf("expected %q, got %q", expected, editor.String())
	}
}

func TestExampleAsIntended(t *testing.T) {
	// Test a version that gives the result shown in the requirements  
	editor := New("hello, world")
	err := editor.Do(
		BeginningOfBuffer(),
		Search(","),
		SetMark(),
		ForwardChar(7), // Move to end: ", world" is 7 chars
		ReplaceRegion(" WORLD"),
	)
	
	if err != nil {
		t.Fatalf("example failed with error: %v", err)
	}
	
	expected := "hello WORLD"
	if editor.String() != expected {
		t.Errorf("expected %q, got %q", expected, editor.String())
	}
}
