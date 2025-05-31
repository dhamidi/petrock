package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// TextInputProps defines the properties for a text input component
type TextInputProps struct {
	Type            string // "text", "email", "password", "search", "tel", "url"
	Placeholder     string
	Value           string
	Required        bool
	Disabled        bool
	ValidationState string // "valid", "invalid", "pending", or empty string for default
	ID              string
	Name            string
	AutoComplete    string
}

// TextInput creates a styled text input component
func TextInput(props TextInputProps) g.Node {
	// Set default type if not provided
	inputType := props.Type
	if inputType == "" {
		inputType = "text"
	}

	// Build CSS classes based on state
	classes := []string{
		"px-3", "py-2", "border", "rounded-md", "text-sm", "leading-5",
		"transition-colors", "duration-200", "focus:outline-none", "focus:ring-2",
		"focus:ring-blue-500", "focus:border-blue-500",
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

	// Build input attributes
	var attrs []g.Node
	attrs = append(attrs,
		html.Type(inputType),
		CSSClass(classes...),
	)

	if props.Placeholder != "" {
		attrs = append(attrs, html.Placeholder(props.Placeholder))
	}

	if props.Value != "" {
		attrs = append(attrs, html.Value(props.Value))
	}

	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}

	if props.Name != "" {
		attrs = append(attrs, html.Name(props.Name))
	}

	if props.AutoComplete != "" {
		attrs = append(attrs, html.AutoComplete(props.AutoComplete))
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

	return html.Input(attrs...)
}

// TextInputStyles provides CSS for text input components
const TextInputStyles = `
.text-input-group {
	display: flex;
	flex-direction: column;
	gap: 0.25rem;
}

.text-input-label {
	font-size: 0.875rem;
	font-weight: 500;
	color: #374151;
}

.text-input-label.required::after {
	content: ' *';
	color: #ef4444;
}

.text-input-help {
	font-size: 0.75rem;
	color: #6b7280;
}

.text-input-error {
	font-size: 0.75rem;
	color: #ef4444;
}
`