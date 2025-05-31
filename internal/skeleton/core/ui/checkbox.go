package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// CheckboxProps defines the properties for a checkbox component
type CheckboxProps struct {
	Checked  bool
	Required bool
	Disabled bool
	Value    string
	Label    string
	ID       string
	Name     string
}

// Checkbox creates a styled checkbox component with label
func Checkbox(props CheckboxProps) g.Node {
	// Generate ID if not provided
	id := props.ID
	if id == "" {
		id = "checkbox-" + props.Value
	}

	// Build checkbox input attributes
	var inputAttrs []g.Node
	inputAttrs = append(inputAttrs,
		html.Type("checkbox"),
		html.ID(id),
		CSSClass("mr-2", "h-4", "w-4", "text-blue-600", "focus:ring-blue-500", "border-gray-300", "rounded"),
	)

	if props.Value != "" {
		inputAttrs = append(inputAttrs, html.Value(props.Value))
	}

	if props.Name != "" {
		inputAttrs = append(inputAttrs, html.Name(props.Name))
	}

	if props.Checked {
		inputAttrs = append(inputAttrs, html.Checked())
	}

	if props.Required {
		inputAttrs = append(inputAttrs, html.Required())
	}

	if props.Disabled {
		inputAttrs = append(inputAttrs, html.Disabled())
		// Add disabled styling
		inputAttrs = append(inputAttrs, CSSClass("opacity-50", "cursor-not-allowed"))
	}

	// Build label classes
	labelClasses := []string{"flex", "items-center", "text-sm", "font-medium", "text-gray-700"}
	if props.Disabled {
		labelClasses = append(labelClasses, "opacity-50", "cursor-not-allowed")
	} else {
		labelClasses = append(labelClasses, "cursor-pointer")
	}

	// Create the checkbox with label
	return html.Label(
		CSSClass(labelClasses...),
		html.For(id),
		html.Input(inputAttrs...),
		g.Text(props.Label),
	)
}

