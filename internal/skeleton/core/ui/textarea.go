package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"strconv"
)

// TextAreaProps defines the properties for a textarea component
type TextAreaProps struct {
	Placeholder     string
	Value           string
	Rows            int
	Required        bool
	Disabled        bool
	ValidationState string // "valid", "invalid", "pending", or empty string for default
	ID              string
	Name            string
	Resize          string // "none", "vertical", "horizontal", "both", or empty for default
}

// TextArea creates a styled textarea component
func TextArea(props TextAreaProps) g.Node {
	// Set default rows if not provided
	rows := props.Rows
	if rows <= 0 {
		rows = 3
	}

	// Build CSS classes based on state
	classes := []string{
		"px-3", "py-2", "border", "rounded-md", "text-sm", "leading-5",
		"transition-colors", "duration-200", "focus:outline-none", "focus:ring-2",
		"focus:ring-blue-500", "focus:border-blue-500", "resize-vertical",
	}

	// Add resize classes
	switch props.Resize {
	case "none":
		classes = append(classes, "resize-none")
	case "vertical":
		classes = append(classes, "resize-y")
	case "horizontal":
		classes = append(classes, "resize-x")
	case "both":
		classes = append(classes, "resize")
	default:
		classes = append(classes, "resize-y")
	}

	// Add validation state classes
	switch props.ValidationState {
	case "valid":
		classes = append(classes, "border-green-500", "bg-green-50")
	case "invalid":
		classes = append(classes, "border-red-500", "bg-red-50")
	case "pending":
		classes = append(classes, "border-yellow-500", "bg-yellow-50")
	default:
		classes = append(classes, "border-gray-300", "bg-white")
	}

	// Add disabled state classes
	if props.Disabled {
		classes = append(classes, "opacity-50", "cursor-not-allowed", "bg-gray-100")
	}

	// Build textarea attributes
	var attrs []g.Node
	attrs = append(attrs,
		html.Rows(strconv.Itoa(rows)),
		CSSClass(classes...),
	)

	if props.Placeholder != "" {
		attrs = append(attrs, html.Placeholder(props.Placeholder))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.Required {
		attrs = append(attrs, html.Required())
	}

	if props.Disabled {
		attrs = append(attrs, html.Disabled())
	}

	// Add ARIA attributes for validation state
	switch props.ValidationState {
	case "invalid":
		attrs = append(attrs, html.Aria("invalid", "true"))
	case "valid":
		attrs = append(attrs, html.Aria("invalid", "false"))
	}

	// Add value as text content if provided
	var content []g.Node
	if props.Value != "" {
		content = append(content, g.Text(props.Value))
	}

	return html.Textarea(append(attrs, content...)...)
}

// TextAreaStyles provides CSS for textarea components
const TextAreaStyles = `
.textarea-group {
	display: flex;
	flex-direction: column;
	gap: 0.25rem;
}

.textarea-label {
	font-size: 0.875rem;
	font-weight: 500;
	color: #374151;
}

.textarea-label.required::after {
	content: ' *';
	color: #ef4444;
}

.textarea-help {
	font-size: 0.75rem;
	color: #6b7280;
}

.textarea-error {
	font-size: 0.75rem;
	color: #ef4444;
}
`