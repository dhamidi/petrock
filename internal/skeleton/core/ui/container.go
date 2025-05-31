package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ContainerProps defines the properties for the Container component
type ContainerProps struct {
	Variant  string // default, narrow, wide, full
	MaxWidth string // custom max-width value (optional)
}

// Container creates a responsive container component with various width variants
func Container(props ContainerProps, children ...g.Node) g.Node {
	// Set default variant if not specified
	variant := props.Variant
	if variant == "" {
		variant = "default"
	}

	// Build CSS classes based on variant
	var classes []string
	classes = append(classes, "mx-auto", "px-4") // Base container styles

	switch variant {
	case "narrow":
		classes = append(classes, "max-w-2xl") // ~672px
	case "wide":
		classes = append(classes, "max-w-7xl") // ~1280px
	case "full":
		classes = append(classes, "max-w-none") // Full width
	default: // "default"
		classes = append(classes, "max-w-4xl") // ~896px
	}

	var attributes []g.Node
	
	// Add class attribute
	if len(classes) > 0 {
		attributes = append(attributes, CSSClass(classes...))
	}

	// Add custom max-width if provided
	if props.MaxWidth != "" {
		attributes = append(attributes, Style(map[string]string{
			"max-width": props.MaxWidth,
		}))
	}

	// Combine attributes and children
	var allNodes []g.Node
	allNodes = append(allNodes, attributes...)
	allNodes = append(allNodes, children...)

	return html.Div(allNodes...)
}