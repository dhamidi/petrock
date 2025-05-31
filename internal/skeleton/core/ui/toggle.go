package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ToggleProps defines the properties for a toggle switch component
type ToggleProps struct {
	Checked  bool
	Disabled bool
	Label    string
	ID       string
	Name     string
	Value    string
}

// Toggle creates a styled toggle switch component using Tailwind CSS peer utilities
func Toggle(props ToggleProps) g.Node {
	// Generate ID if not provided
	id := props.ID
	if id == "" {
		id = "toggle-" + props.Name
	}

	// Build checkbox input attributes that powers the toggle
	var inputAttrs []g.Node
	inputAttrs = append(inputAttrs,
		html.Type("checkbox"),
		html.ID(id),
		CSSClass("peer", "sr-only"), // peer class enables peer-* selectors, sr-only hides visually
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

	if props.Disabled {
		inputAttrs = append(inputAttrs, html.Disabled())
	}

	// Toggle switch background - uses peer-* selectors to respond to checkbox state
	backgroundClasses := []string{
		"absolute", "top-0", "left-0", "h-6", "w-11", "rounded-full", "border-2", "border-transparent", 
		"transition-colors", "duration-200", "ease-in-out",
		// Default state (unchecked)
		"bg-gray-200",
		// Checked state (when peer checkbox is checked)
		"peer-checked:bg-blue-600",
		// Focus state
		"peer-focus:ring-2", "peer-focus:ring-blue-500", "peer-focus:ring-offset-2",
	}

	// Disabled state styling
	if props.Disabled {
		backgroundClasses = append(backgroundClasses, "opacity-50", "cursor-not-allowed")
	}

	// Toggle thumb (circle) classes - position responds to checkbox state via peer-*
	// Must be direct sibling of input for peer-* selectors to work
	thumbClasses := []string{
		"absolute", "top-0.5", "left-0.5", "h-5", "w-5", "rounded-full",
		"bg-white", "shadow", "transform", "ring-0", "transition", "duration-200", "ease-in-out",
		"pointer-events-none",
		// Default position (unchecked)
		"translate-x-0",
		// Checked position (when peer checkbox is checked)
		"peer-checked:translate-x-5",
	}

	// Build the toggle switch container with proper sibling structure for peer selectors
	toggleSwitch := html.Label(
		CSSClass("relative", "inline-flex", "h-6", "w-11", "flex-shrink-0", "cursor-pointer"),
		html.For(id),
		html.Input(inputAttrs...),
		// Background and thumb must be direct siblings of input for peer-* to work
		html.Span(CSSClass(backgroundClasses...)),
		html.Span(CSSClass(thumbClasses...)),
	)

	// If no label text, return just the toggle
	if props.Label == "" {
		return toggleSwitch
	}

	// Container with label text and toggle
	containerClasses := []string{"flex", "items-center", "text-sm", "font-medium", "text-gray-700"}
	if props.Disabled {
		containerClasses = append(containerClasses, "opacity-50")
	}

	return html.Div(
		CSSClass(containerClasses...),
		html.Span(CSSClass("mr-3"), g.Text(props.Label)),
		toggleSwitch,
	)
}

