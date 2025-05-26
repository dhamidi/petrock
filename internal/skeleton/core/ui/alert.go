package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// AlertProps defines the properties for the Alert component
type AlertProps struct {
	Type        string // success, warning, error, info
	Title       string // optional title for the alert
	Message     string // main alert message
	Dismissible bool   // whether the alert can be dismissed
}

// Alert creates a feedback component for displaying important messages to users
func Alert(props AlertProps) g.Node {
	// Set default type if not specified
	alertType := props.Type
	if alertType == "" {
		alertType = "info"
	}

	// Build CSS classes based on alert type
	var classes []string
	classes = append(classes, 
		"p-4", "rounded-md", "border", "relative",
		"flex", "items-start", "space-x-3")

	// Add type-specific styles
	switch alertType {
	case "success":
		classes = append(classes, 
			"bg-green-50", "border-green-200", "text-green-800")
	case "warning":
		classes = append(classes, 
			"bg-yellow-50", "border-yellow-200", "text-yellow-800")
	case "error":
		classes = append(classes, 
			"bg-red-50", "border-red-200", "text-red-800")
	default: // "info"
		classes = append(classes, 
			"bg-blue-50", "border-blue-200", "text-blue-800")
	}

	// Create icon based on alert type
	var icon g.Node
	switch alertType {
	case "success":
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-5", "h-5", "text-green-400"),
				g.Text("✓"), // Success checkmark
			),
		)
	case "warning":
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-5", "h-5", "text-yellow-400"),
				g.Text("⚠"), // Warning triangle
			),
		)
	case "error":
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-5", "h-5", "text-red-400"),
				g.Text("✕"), // Error X
			),
		)
	default: // "info"
		icon = html.Div(
			CSSClass("flex-shrink-0"),
			html.Div(
				CSSClass("w-5", "h-5", "text-blue-400"),
				g.Text("ℹ"), // Info i
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
			html.H3(
				CSSClass("text-sm", "font-medium"),
				g.Text(props.Title),
			),
		)
	}

	// Add message
	if props.Message != "" {
		messageClasses := []string{"text-sm"}
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
			CSSClass("flex-1"),
			g.Group(contentDiv),
		),
	)

	// Add dismiss button if dismissible
	if props.Dismissible {
		dismissButton := html.Button(
			CSSClass(
				"absolute", "top-2", "right-2",
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

	// Add ARIA role for accessibility
	var attributes []g.Node
	attributes = append(attributes, CSSClass(classes...))
	attributes = append(attributes, html.Role("alert"))
	attributes = append(attributes, html.Aria("live", "polite"))

	return html.Div(
		g.Group(append(attributes, content...)),
	)
}