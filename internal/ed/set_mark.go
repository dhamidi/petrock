package ed

// SetMark sets the mark at the current cursor position
func SetMark() Operation {
	return func(e *Editor) error {
		e.mark = e.position
		e.marked = true
		return nil
	}
}
