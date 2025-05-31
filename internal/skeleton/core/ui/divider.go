package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// DividerProps defines the properties for the Divider component
type DividerProps struct {
	Variant string // solid, dashed, dotted
	Margin  string // none, small, medium, large
}

// Divider creates a horizontal divider component with spacing variants
func Divider(props DividerProps) g.Node {
	// Set default variant if not specified
	variant := props.Variant
	if variant == "" {
		variant = "solid"
	}

	// Set default margin if not specified
	margin := props.Margin
	if margin == "" {
		margin = "medium"
	}

	// Build CSS classes based on variant and margin
	var classes []string
	classes = append(classes, "border-gray-200", "border-0") // Base divider styles

	// Add variant-specific styles
	switch variant {
	case "dashed":
		classes = append(classes, "border-t", "border-dashed")
	case "dotted":
		classes = append(classes, "border-t", "border-dotted")
	default: // "solid"
		classes = append(classes, "border-t")
	}

	// Add margin styles
	switch margin {
	case "none":
		// No margin classes
	case "small":
		classes = append(classes, "my-2")
	case "large":
		classes = append(classes, "my-8")
	default: // "medium"
		classes = append(classes, "my-4")
	}

	return html.Hr(
		CSSClass(classes...),
		html.Role("separator"),
	)
}