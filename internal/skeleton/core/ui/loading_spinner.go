package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// LoadingSpinnerProps defines the properties for the LoadingSpinner component
type LoadingSpinnerProps struct {
	Size  string // small, medium, large
	Color string // primary, secondary, white
	Label string // optional accessible label
}

// LoadingSpinner creates a CSS-only animated loading indicator
func LoadingSpinner(props LoadingSpinnerProps) g.Node {
	// Set default size if not specified
	size := props.Size
	if size == "" {
		size = "medium"
	}

	// Set default color if not specified
	color := props.Color
	if color == "" {
		color = "primary"
	}

	// Set default label if not specified
	label := props.Label
	if label == "" {
		label = "Loading..."
	}

	// Build CSS classes based on size and color
	var classes []string
	classes = append(classes, 
		"inline-block", "animate-spin", "rounded-full", "border-2")

	// Add size-specific styles
	switch size {
	case "small":
		classes = append(classes, "w-4", "h-4", "border-2")
	case "large":
		classes = append(classes, "w-8", "h-8", "border-2")
	default: // "medium"
		classes = append(classes, "w-6", "h-6", "border-2")
	}

	// Add color-specific styles
	switch color {
	case "secondary":
		classes = append(classes, 
			"border-gray-200", "border-t-gray-600")
	case "white":
		classes = append(classes, 
			"border-white/20", "border-t-white")
	default: // "primary"
		classes = append(classes, 
			"border-blue-200", "border-t-blue-600")
	}

	return html.Div(
		CSSClass(classes...),
		html.Role("status"),
		html.Aria("label", label),
		// Visually hidden text for screen readers
		html.Span(
			CSSClass("sr-only"),
			g.Text(label),
		),
	)
}

// LoadingSpinnerWithText creates a loading spinner with visible text label
func LoadingSpinnerWithText(props LoadingSpinnerProps) g.Node {
	// Set default label if not specified
	label := props.Label
	if label == "" {
		label = "Loading..."
	}

	return html.Div(
		CSSClass("flex", "items-center", "space-x-2"),
		LoadingSpinner(props),
		html.Span(
			CSSClass("text-sm", "text-gray-600"),
			g.Text(label),
		),
	)
}

// LoadingOverlay creates a full-screen loading overlay
func LoadingOverlay(props LoadingSpinnerProps) g.Node {
	// Set default label if not specified
	label := props.Label
	if label == "" {
		label = "Loading..."
	}

	// Create a larger spinner for overlay
	overlayProps := props
	if overlayProps.Size == "" || overlayProps.Size == "small" || overlayProps.Size == "medium" {
		overlayProps.Size = "large"
	}

	return html.Div(
		CSSClass(
			"fixed", "inset-0", "z-50", 
			"flex", "items-center", "justify-center",
			"bg-white", "bg-opacity-75",
		),
		html.Div(
			CSSClass("text-center"),
			LoadingSpinner(overlayProps),
			html.P(
				CSSClass("mt-4", "text-sm", "text-gray-600"),
				g.Text(label),
			),
		),
	)
}