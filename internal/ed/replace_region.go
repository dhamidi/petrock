package ed

// ReplaceRegion replaces the text between mark and current position with the given text
func ReplaceRegion(replacement string) Operation {
	return func(e *Editor) error {
		if !e.marked {
			return Error{
				Operation: "ReplaceRegion",
				Message:   "no mark set",
				Position:  e.position,
			}
		}
		
		start := e.mark
		end := e.position
		
		// Ensure start <= end
		if start > end {
			start, end = end, start
		}
		
		e.replaceRange(start, end, replacement)
		
		// Clear the mark after replacement
		e.marked = false
		
		return nil
	}
}
