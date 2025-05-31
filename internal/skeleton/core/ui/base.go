package ui

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// CSSClass creates a CSS class attribute from multiple class names
// It joins non-empty class names with spaces and returns a class attribute
func CSSClass(classes ...string) g.Node {
	var validClasses []string
	for _, class := range classes {
		if trimmed := strings.TrimSpace(class); trimmed != "" {
			validClasses = append(validClasses, trimmed)
		}
	}
	if len(validClasses) == 0 {
		return nil
	}
	return html.Class(strings.Join(validClasses, " "))
}

// Style creates an inline style attribute from a map of CSS properties
// Keys should be CSS property names and values should be CSS values
func Style(props map[string]string) g.Node {
	if len(props) == 0 {
		return nil
	}

	var styles []string
	for property, value := range props {
		if property != "" && value != "" {
			styles = append(styles, property+":"+value)
		}
	}

	if len(styles) == 0 {
		return nil
	}

	return html.Style(strings.Join(styles, ";"))
}



// Attrs combines multiple attributes into a slice for easy spreading
func Attrs(attrs ...g.Node) []g.Node {
	var result []g.Node
	for _, attr := range attrs {
		// Only add non-nil attributes
		if attr != nil {
			result = append(result, attr)
		}
	}
	return result
}

// ConditionalAttr returns the attribute if the condition is true, otherwise returns nil
func ConditionalAttr(condition bool, attr g.Node) g.Node {
	if condition {
		return attr
	}
	return nil
}

// Text creates a text node, handling empty strings gracefully
func Text(content string) g.Node {
	return g.Text(content)
}

// If conditionally renders content based on a boolean condition
func If(condition bool, content g.Node) g.Node {
	if condition {
		return content
	}
	return g.Text("")
}