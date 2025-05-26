package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ToastProps defines the properties for the Toast component
type ToastProps struct {
	Type        string // success, warning, error, info
	Title       string // optional title for the toast
	Message     string // main toast message
	Dismissible bool   // whether the toast can be dismissed
	Position    string // top-right, top-left, bottom-right, bottom-left
}

// Toast creates a temporary notification component that appears overlay on the page
func Toast(props ToastProps) g.Node {
	// Set default type if not specified
	toastType := props.Type
	if toastType == "" {
		toastType = "info"
	}

	// Set default position if not specified
	position := props.Position
	if position == "" {
		position = "top-right"
	}

	// Build CSS classes based on toast type and position
	var classes []string
	classes = append(classes, 
		"fixed", "z-50", "p-4", "max-w-sm", "w-full",
		"bg-white", "rounded-lg", "shadow-lg", "border",
		"transform", "transition-all", "duration-300",
		"flex", "items-start", "space-x-3")

	// Add position-specific styles
	switch position {
	case "top-left":
		classes = append(classes, "top-4", "left-4")
	case "bottom-right":
		classes = append(classes, "bottom-4", "right-4")
	case "bottom-left":
		classes = append(classes, "bottom-4", "left-4")
	default: // "top-right"
		classes = append(classes, "top-4", "right-4")
	}

	// Add type-specific border and accent colors
	switch toastType {
	case "success":
		classes = append(classes, "border-l-4", "border-l-green-400")
	case "warning":
		classes = append(classes, "border-l-4", "border-l-yellow-400")
	case "error":
		classes = append(classes, "border-l-4", "border-l-red-400")
	default: // "info"
		classes = append(classes, "border-l-4", "border-l-blue-400")
	}

	// Create icon based on toast type
	var icon g.Node
	switch toastType {
	case "success":
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-6", "h-6", "text-green-500", "font-bold"),
				g.Text("✓"),
			),
		)
	case "warning":
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-6", "h-6", "text-yellow-500", "font-bold"),
				g.Text("⚠"),
			),
		)
	case "error":
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-6", "h-6", "text-red-500", "font-bold"),
				g.Text("✕"),
			),
		)
	default: // "info"
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-6", "h-6", "text-blue-500", "font-bold"),
				g.Text("ℹ"),
			),
		)
	}

	// Build content
	var content []g.Node
	content = append(content, icon)

	// Add text content
	contentDiv := []g.Node{}
	
	// Add title if provided
	if props.Title != "" {
		contentDiv = append(contentDiv, 
			html.H4(
				CSSClass("text-sm", "font-semibold", "text-gray-900"),
				g.Text(props.Title),
			),
		)
	}

	// Add message
	if props.Message != "" {
		messageClasses := []string{"text-sm", "text-gray-700"}
		if props.Title != "" {
			messageClasses = append(messageClasses, "mt-1")
		}
		contentDiv = append(contentDiv,
			html.P(
				CSSClass(messageClasses...),
				g.Text(props.Message),
			),
		)
	}

	content = append(content, 
		html.Div(
			CSSClass("flex-1", "min-w-0"),
			g.Group(contentDiv),
		),
	)

	// Add dismiss button if dismissible
	if props.Dismissible {
		dismissButton := html.Button(
			CSSClass(
				"flex-shrink-0", "ml-4",
				"text-gray-400", "hover:text-gray-600",
				"focus:outline-none", "focus:text-gray-600",
				"transition-colors", "duration-200",
			),
			html.Type("button"),
			html.Aria("label", "Dismiss"),
			g.Text("×"),
		)
		content = append(content, dismissButton)
	}

	// Add ARIA attributes for accessibility
	var attributes []g.Node
	attributes = append(attributes, CSSClass(classes...))
	attributes = append(attributes, html.Role("alert"))
	attributes = append(attributes, html.Aria("live", "assertive"))

	return html.Div(
		g.Group(append(attributes, content...)),
	)
}