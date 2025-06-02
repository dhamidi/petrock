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

// toTitleCase converts a string to TitleCase (e.g., "create_post" -> "CreatePost", "schedule-publication" -> "SchedulePublication")
func toTitleCase(input string) string {
	if input == "" {
		return ""
	}

	// Convert kebab-case to snake_case first
	normalized := kebabToSnakeCase(input)
	
	parts := strings.Split(normalized, "_")
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

// kebabToSnakeCase converts kebab-case to snake_case (e.g., "schedule-publication" -> "schedule_publication")
func kebabToSnakeCase(input string) string {
	return strings.ReplaceAll(input, "-", "_")
}

// isValidKebabName checks if a name follows kebab-case naming pattern
func isValidKebabName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// Must start with letter
	if !isLetter(rune(name[0])) {
		return false
	}

	// Rest can be letters, numbers, or hyphens
	for _, r := range name[1:] {
		if !isLetter(r) && !isDigit(r) && r != '-' {
			return false
		}
	}

	// Cannot end with hyphen
	return name[len(name)-1] != '-'
}

// ValidateEntityName checks if an entity name is valid for component generation
func ValidateEntityName(entityName string) error {
	if entityName == "" {
		return fmt.Errorf("entity name cannot be empty")
	}

	// Check for valid naming pattern (letters, numbers, underscores, or kebab-case)
	if !isValidEntityName(entityName) && !isValidKebabName(entityName) {
		return fmt.Errorf("invalid entity name %q: must contain only letters, numbers, underscores, or hyphens (kebab-case)", entityName)
	}

	return nil
}
