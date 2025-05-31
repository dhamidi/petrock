package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ButtonGroupProps defines the properties for the ButtonGroup component
type ButtonGroupProps struct {
	Orientation string // horizontal, vertical
	Spacing     string // none, small, medium, large
}

// ButtonGroup creates a container for grouping related buttons with consistent spacing and orientation
func ButtonGroup(props ButtonGroupProps, children ...g.Node) g.Node {
	// Set default orientation if not specified
	orientation := props.Orientation
	if orientation == "" {
		orientation = "horizontal"
	}

	// Set default spacing if not specified
	spacing := props.Spacing
	if spacing == "" {
		spacing = "medium"
	}

	// Build CSS classes based on orientation and spacing
	var classes []string
	
	// Base button group styles
	classes = append(classes, "inline-flex")

	// Add orientation-specific styles
	switch orientation {
	case "vertical":
		classes = append(classes, "flex-col")
		// Add vertical spacing based on spacing prop
		switch spacing {
		case "none":
			// No spacing classes
		case "small":
			classes = append(classes, "space-y-1")
		case "large":
			classes = append(classes, "space-y-4")
		default: // "medium"
			classes = append(classes, "space-y-2")
		}
	default: // "horizontal"
		classes = append(classes, "flex-row")
		// Add horizontal spacing based on spacing prop
		switch spacing {
		case "none":
			// No spacing classes
		case "small":
			classes = append(classes, "space-x-1")
		case "large":
			classes = append(classes, "space-x-4")
		default: // "medium"
			classes = append(classes, "space-x-2")
		}
	}

	// Build attributes
	var attributes []g.Node
	attributes = append(attributes, CSSClass(classes...))

	// Add ARIA role for accessibility
	attributes = append(attributes, html.Role("group"))

	// Combine attributes and children
	var allNodes []g.Node
	allNodes = append(allNodes, attributes...)
	allNodes = append(allNodes, children...)

	return html.Div(allNodes...)
}