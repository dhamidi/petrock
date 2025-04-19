package ui

// pluralize returns "s" if count is not 1, otherwise returns empty string
func Pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}