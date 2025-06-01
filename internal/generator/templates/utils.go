package templates

import (
	"fmt"
	"strings"
)

// isValidEntityName checks if an entity name follows valid patterns
func isValidEntityName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// Must start with letter
	if !isLetter(rune(name[0])) {
		return false
	}

	// Rest can be letters, numbers, or underscores
	for _, r := range name[1:] {
		if !isLetter(r) && !isDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

// toTitleCase converts a string to TitleCase (e.g., "create_post" -> "CreatePost")
func toTitleCase(input string) string {
	if input == "" {
		return ""
	}

	parts := strings.Split(input, "_")
	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}
	return result
}

// isLetter checks if a rune is a letter
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isDigit checks if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// ValidateEntityName checks if an entity name is valid for component generation
func ValidateEntityName(entityName string) error {
	if entityName == "" {
		return fmt.Errorf("entity name cannot be empty")
	}

	// Check for valid naming pattern (letters, numbers, underscores)
	if !isValidEntityName(entityName) {
		return fmt.Errorf("invalid entity name %q: must contain only letters, numbers, and underscores", entityName)
	}

	return nil
}
