package ed

// BeginningOfBuffer moves the cursor to the beginning of the buffer
func BeginningOfBuffer() Operation {
	return func(e *Editor) error {
		e.setPosition(0)
		return nil
	}
}
