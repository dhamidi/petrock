package ed

import "strings"

// Search searches for the given text starting from the current position
// and moves the cursor to the beginning of the match
func Search(text string) Operation {
	return func(e *Editor) error {
		if text == "" {
			return Error{
				Operation: "Search",
				Message:   "search text cannot be empty",
				Position:  e.position,
			}
		}
		
		// Search from current position onwards
		remaining := e.content[e.position:]
		index := strings.Index(remaining, text)
		
		if index == -1 {
			return Error{
				Operation: "Search",
				Message:   "text not found: " + text,
				Position:  e.position,
			}
		}
		
		// Move cursor to the beginning of the found text
		newPos := e.position + index
		e.setPosition(newPos)
		
		return nil
	}
}
