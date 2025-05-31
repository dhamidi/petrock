package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// RadioProps defines the properties for a radio button component
type RadioProps struct {
	Checked  bool
	Required bool
	Disabled bool
	Value    string
	Name     string
	Label    string
	ID       string
}

// Radio creates a styled radio button component with label
func Radio(props RadioProps) g.Node {
	// Generate ID if not provided
	id := props.ID
	if id == "" {
		id = "radio-" + props.Name + "-" + props.Value
	}

	// Build radio input attributes
	var inputAttrs []g.Node
	inputAttrs = append(inputAttrs,
		html.Type("radio"),
		html.ID(id),
		CSSClass("mr-2", "h-4", "w-4", "text-blue-600", "focus:ring-blue-500", "border-gray-300"),
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

	// Create the radio button with label
	return html.Label(
		CSSClass(labelClasses...),
		html.For(id),
		html.Input(inputAttrs...),
		g.Text(props.Label),
	)
}

// RadioGroup creates a group of radio buttons with shared name
type RadioGroupProps struct {
	Name     string
	Value    string   // Currently selected value
	Options  []RadioOption
	Required bool
	Disabled bool
	Label    string
	Vertical bool // If true, stack vertically; otherwise horizontal
}

// RadioOption represents a single option in a radio group
type RadioOption struct {
	Value    string
	Label    string
	Disabled bool
}

// RadioGroup creates a group of related radio buttons
func RadioGroup(props RadioGroupProps) g.Node {
	// Container classes
	containerClasses := []string{"space-y-2"}
	if !props.Vertical {
		containerClasses = []string{"flex", "flex-wrap", "gap-4"}
	}

	// Build radio buttons
	var radios []g.Node
	for _, option := range props.Options {
		radios = append(radios, Radio(RadioProps{
			Checked:  option.Value == props.Value,
			Required: props.Required,
			Disabled: props.Disabled || option.Disabled,
			Value:    option.Value,
			Name:     props.Name,
			Label:    option.Label,
		}))
	}

	// If there's a group label, wrap in fieldset
	if props.Label != "" {
		return html.FieldSet(
			CSSClass("border", "border-gray-300", "rounded-md", "p-4"),
			html.Legend(
				CSSClass("text-sm", "font-medium", "text-gray-900", "px-2"),
				g.Text(props.Label),
			),
			html.Div(
				CSSClass(containerClasses...),
				g.Group(radios),
			),
		)
	}

	// Just return the radio group without fieldset
	return html.Div(
		CSSClass(containerClasses...),
		g.Group(radios),
	)
}

