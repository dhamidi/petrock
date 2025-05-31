package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"strconv"
)

// BadgeProps defines the properties for the Badge component
type BadgeProps struct {
	Variant string // primary, secondary, success, warning, error, info
	Size    string // small, medium, large
	Count   int    // numerical count to display (0 means no count shown)
}

// Badge creates a small status indicator or count display component
func Badge(props BadgeProps, children ...g.Node) g.Node {
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

	// Build CSS classes based on variant and size
	var classes []string
	classes = append(classes, 
		"inline-flex", "items-center", "justify-center",
		"font-medium", "rounded-full", "text-center")

	// Add variant-specific styles
	switch variant {
	case "secondary":
		classes = append(classes, 
			"bg-gray-100", "text-gray-800")
	case "success":
		classes = append(classes, 
			"bg-green-100", "text-green-800")
	case "warning":
		classes = append(classes, 
			"bg-yellow-100", "text-yellow-800")
	case "error":
		classes = append(classes, 
			"bg-red-100", "text-red-800")
	case "info":
		classes = append(classes, 
			"bg-blue-100", "text-blue-800")
	default: // "primary"
		classes = append(classes, 
			"bg-blue-600", "text-white")
	}

	// Add size-specific styles
	switch size {
	case "small":
		classes = append(classes, "px-2", "py-0.5", "text-xs", "min-w-[1.25rem]", "h-5")
	case "large":
		classes = append(classes, "px-4", "py-1", "text-base", "min-w-[2rem]", "h-8")
	default: // "medium"
		classes = append(classes, "px-3", "py-1", "text-sm", "min-w-[1.5rem]", "h-6")
	}

	// Determine content to display
	var content []g.Node
	if props.Count > 0 {
		// Display count
		countText := strconv.Itoa(props.Count)
		// Show 99+ for counts over 99
		if props.Count > 99 {
			countText = "99+"
		}
		content = append(content, g.Text(countText))
	} else if len(children) > 0 {
		// Display provided children
		content = children
	} else {
		// Empty badge (just a dot)
		switch size {
		case "small":
			classes = append(classes, "w-2", "h-2")
		case "large":
			classes = append(classes, "w-4", "h-4")
		default: // "medium"
			classes = append(classes, "w-3", "h-3")
		}
		// Remove padding for dot badges
		classes = filterClasses(classes, "px-2", "px-3", "px-4", "py-0.5", "py-1", "min-w-[1.25rem]", "min-w-[1.5rem]", "min-w-[2rem]")
	}

	return html.Span(
		CSSClass(classes...),
		g.Group(content),
	)
}

// filterClasses removes specified classes from a slice
func filterClasses(classes []string, toRemove ...string) []string {
	var filtered []string
	removeMap := make(map[string]bool)
	for _, class := range toRemove {
		removeMap[class] = true
	}
	
	for _, class := range classes {
		if !removeMap[class] {
			filtered = append(filtered, class)
		}
	}
	return filtered
}