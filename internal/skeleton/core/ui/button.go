package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ButtonProps defines the properties for the Button component
type ButtonProps struct {
	Variant  string // primary, secondary, danger, link
	Size     string // small, medium, large
	Type     string // button, submit, reset
	Disabled bool   // whether the button is disabled
}

// Button creates an interactive button component with various styles and states
func Button(props ButtonProps, children ...g.Node) g.Node {
	// Set default variant if not specified
	variant := props.Variant
	if variant == "" {
		variant = "primary"
	}

	// Set default size if not specified
	size := props.Size
	if size == "" {
		size = "medium"
	}

	// Set default type if not specified
	buttonType := props.Type
	if buttonType == "" {
		buttonType = "button"
	}

	// Build CSS classes based on variant, size, and state
	var classes []string
	
	// Base button styles
	classes = append(classes, 
		"inline-flex", "items-center", "justify-center", 
		"font-medium", "rounded-md", "transition-colors", 
		"focus:outline-none", "focus:ring-2", "focus:ring-offset-2")

	// Add variant-specific styles
	switch variant {
	case "secondary":
		if props.Disabled {
			classes = append(classes, "bg-gray-100", "text-gray-400", "cursor-not-allowed")
		} else {
			classes = append(classes, 
				"bg-white", "text-gray-700", "border", "border-gray-300",
				"hover:bg-gray-50", "focus:ring-gray-500")
		}
	case "danger":
		if props.Disabled {
			classes = append(classes, "bg-red-300", "text-red-100", "cursor-not-allowed")
		} else {
			classes = append(classes, 
				"bg-red-600", "text-white", "hover:bg-red-700", "focus:ring-red-500")
		}
	case "link":
		if props.Disabled {
			classes = append(classes, "text-gray-400", "cursor-not-allowed")
		} else {
			classes = append(classes, 
				"text-blue-600", "hover:text-blue-800", "underline", 
				"bg-transparent", "focus:ring-blue-500")
		}
	default: // "primary"
		if props.Disabled {
			classes = append(classes, "bg-blue-300", "text-blue-100", "cursor-not-allowed")
		} else {
			classes = append(classes, 
				"bg-blue-600", "text-white", "hover:bg-blue-700", "focus:ring-blue-500")
		}
	}

	// Add size-specific styles
	switch size {
	case "small":
		classes = append(classes, "px-3", "py-1.5", "text-sm")
	case "large":
		classes = append(classes, "px-6", "py-3", "text-lg")
	default: // "medium"
		classes = append(classes, "px-4", "py-2", "text-base")
	}

	// Build attributes
	var attributes []g.Node
	attributes = append(attributes, CSSClass(classes...))
	attributes = append(attributes, html.Type(buttonType))

	// Add disabled attribute if disabled
	if props.Disabled {
		attributes = append(attributes, html.Disabled())
	}

	// Add ARIA attributes for accessibility
	if props.Disabled {
		attributes = append(attributes, html.Aria("disabled", "true"))
	}

	// Combine attributes and children
	var allNodes []g.Node
	allNodes = append(allNodes, attributes...)
	allNodes = append(allNodes, children...)

	return html.Button(allNodes...)
}