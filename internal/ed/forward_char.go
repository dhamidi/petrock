package ed

// ForwardChar moves the cursor forward by the specified number of characters
func ForwardChar(n int) Operation {
	return func(e *Editor) error {
		newPos := e.position + n
		e.setPosition(newPos)
		return nil
	}
}
