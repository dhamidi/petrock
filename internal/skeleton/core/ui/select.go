package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// SelectOption represents an option in a select dropdown
type SelectOption struct {
	Value    string
	Label    string
	Disabled bool
	Selected bool
}

// SelectProps defines the properties for a select component
type SelectProps struct {
	Value           string
	Options         []SelectOption
	Required        bool
	Disabled        bool
	ValidationState string // "valid", "invalid", "pending", or empty string for default
	ID              string
	Name            string
	Placeholder     string // Used for an empty first option
	Multiple        bool
}

// Select creates a styled select dropdown component
func Select(props SelectProps) g.Node {
	// Build CSS classes based on state
	classes := []string{
		"px-3", "py-2", "border", "rounded-md", "text-sm", "leading-5",
		"transition-colors", "duration-200", "focus:outline-none", "focus:ring-2",
		"focus:ring-blue-500", "focus:border-blue-500", "bg-white",
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
		classes = append(classes, "border-gray-300")
	}

	// Add disabled state classes
	if props.Disabled {
		classes = append(classes, "opacity-50", "cursor-not-allowed", "bg-gray-100")
	}

	// Build select attributes
	var attrs []g.Node
	attrs = append(attrs,
		CSSClass(classes...),
	)

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

	if props.Multiple {
		attrs = append(attrs, html.Multiple())
	}

	// Add ARIA attributes for validation state
	switch props.ValidationState {
	case "invalid":
		attrs = append(attrs, html.Aria("invalid", "true"))
	case "valid":
		attrs = append(attrs, html.Aria("invalid", "false"))
	}

	// Build options
	var options []g.Node

	// Add placeholder option if provided
	if props.Placeholder != "" && !props.Multiple {
		options = append(options, html.Option(
			html.Value(""),
			html.Disabled(),
			g.Text(props.Placeholder),
		))
	}

	// Add regular options
	for _, option := range props.Options {
		var optionAttrs []g.Node
		optionAttrs = append(optionAttrs, html.Value(option.Value))

		if option.Disabled {
			optionAttrs = append(optionAttrs, html.Disabled())
		}

		// Check if this option should be selected
		selected := option.Selected
		if !selected && props.Value != "" && option.Value == props.Value {
			selected = true
		}

		if selected {
			optionAttrs = append(optionAttrs, html.Selected())
		}

		options = append(options, html.Option(
			append(optionAttrs, g.Text(option.Label))...,
		))
	}

	return html.Select(append(attrs, options...)...)
}

// SelectStyles provides CSS for select components
const SelectStyles = `
.select-group {
	display: flex;
	flex-direction: column;
	gap: 0.25rem;
}

.select-label {
	font-size: 0.875rem;
	font-weight: 500;
	color: #374151;
}

.select-label.required::after {
	content: ' *';
	color: #ef4444;
}

.select-help {
	font-size: 0.75rem;
	color: #6b7280;
}

.select-error {
	font-size: 0.75rem;
	color: #ef4444;
}

/* Custom arrow for select */
select {
	appearance: none;
	background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='m6 8 4 4 4-4'/%3e%3c/svg%3e");
	background-position: right 0.5rem center;
	background-repeat: no-repeat;
	background-size: 1.5em 1.5em;
	padding-right: 2.5rem;
}

select[multiple] {
	background-image: none;
	padding-right: 0.75rem;
}
`