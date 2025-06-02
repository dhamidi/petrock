package ed

import (
	"fmt"
)

// Editor represents a text editor buffer similar to an Emacs buffer
type Editor struct {
	content  string
	position int    // current cursor position  
	mark     int    // mark position for regions
	marked   bool   // whether mark is set
}

// Error represents an error from an editor operation
type Error struct {
	Operation string
	Message   string
	Position  int
}

func (e Error) Error() string {
	return fmt.Sprintf("ed: %s at position %d: %s", e.Operation, e.Position, e.Message)
}

// Operation represents a single editor operation
type Operation func(*Editor) error

// New creates a new editor with the given initial content
func New(content string) *Editor {
	return &Editor{
		content:  content,
		position: 0,
		mark:     0,
		marked:   false,
	}
}

// Do executes a sequence of operations, stopping at the first error
func (e *Editor) Do(operations ...Operation) error {
	for _, op := range operations {
		if err := op(e); err != nil {
			return err
		}
	}
	return nil
}

// String returns the current content of the editor
func (e *Editor) String() string {
	return e.content
}

// Position returns the current cursor position
func (e *Editor) Position() int {
	return e.position
}

// Mark returns the current mark position and whether it's set
func (e *Editor) Mark() (int, bool) {
	return e.mark, e.marked
}

// setPosition sets the cursor position, clamping to valid bounds
func (e *Editor) setPosition(pos int) {
	if pos < 0 {
		e.position = 0
	} else if pos > len(e.content) {
		e.position = len(e.content)
	} else {
		e.position = pos
	}
}

// replaceRange replaces content between start and end positions
func (e *Editor) replaceRange(start, end int, replacement string) {
	if start < 0 {
		start = 0
	}
	if end > len(e.content) {
		end = len(e.content)
	}
	if start > end {
		start, end = end, start
	}
	
	e.content = e.content[:start] + replacement + e.content[end:]
	e.position = start + len(replacement)
	
	// Adjust mark if it exists and is affected by the replacement
	if e.marked {
		if e.mark > end {
			// Mark is after replacement, adjust by length difference
			e.mark += len(replacement) - (end - start)
		} else if e.mark > start {
			// Mark is within replacement range, move to start
			e.mark = start
		}
	}
}
