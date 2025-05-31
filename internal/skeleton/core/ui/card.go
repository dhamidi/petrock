package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// CardProps defines the properties for the Card component
type CardProps struct {
	Variant string // default, outlined, elevated
	Padding string // none, small, medium, large
}

// Card creates a structured content container with header, body, and footer sections
func Card(props CardProps, children ...g.Node) g.Node {
	// Set default variant if not specified
	variant := props.Variant
	if variant == "" {
		variant = "default"
	}

	// Set default padding if not specified
	padding := props.Padding
	if padding == "" {
		padding = "medium"
	}

	// Build CSS classes based on variant and padding
	var classes []string
	classes = append(classes, "bg-white", "rounded-lg") // Base card styles

	// Add variant-specific styles
	switch variant {
	case "outlined":
		classes = append(classes, "border", "border-gray-200")
	case "elevated":
		classes = append(classes, "shadow-lg", "border", "border-gray-100")
	default: // "default"
		classes = append(classes, "shadow", "border", "border-gray-200")
	}

	// Add padding styles
	switch padding {
	case "none":
		// No padding classes
	case "small":
		classes = append(classes, "p-3")
	case "large":
		classes = append(classes, "p-8")
	default: // "medium"
		classes = append(classes, "p-6")
	}

	return html.Div(
		CSSClass(classes...),
		g.Group(children),
	)
}

// CardHeader creates a card header section
func CardHeader(children ...g.Node) g.Node {
	return html.Div(
		CSSClass("border-b", "border-gray-200", "pb-4", "mb-4"),
		g.Group(children),
	)
}

// CardBody creates a card body section
func CardBody(children ...g.Node) g.Node {
	return html.Div(
		CSSClass("flex-1"),
		g.Group(children),
	)
}

// CardFooter creates a card footer section
func CardFooter(children ...g.Node) g.Node {
	return html.Div(
		CSSClass("border-t", "border-gray-200", "pt-4", "mt-4"),
		g.Group(children),
	)
}