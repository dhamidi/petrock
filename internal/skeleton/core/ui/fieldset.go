package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// FieldSetProps defines the properties for a fieldset component
type FieldSetProps struct {
	Legend   string // The legend text for the fieldset
	Disabled bool   // Whether the entire fieldset is disabled
	ID       string // Optional ID for the fieldset
	Class    string // Additional CSS classes
}

// FieldSet creates a semantic fieldset with legend and proper accessibility
func FieldSet(props FieldSetProps, children ...g.Node) g.Node {
	// Build fieldset classes
	classes := []string{
		"fieldset", "border", "border-gray-300", "rounded-md", "p-4", "space-y-4",
	}
	
	if props.Disabled {
		classes = append(classes, "opacity-50", "cursor-not-allowed")
	}
	
	if props.Class != "" {
		classes = append(classes, props.Class)
	}

	// Build fieldset attributes
	var attrs []g.Node
	attrs = append(attrs, CSSClass(classes...))
	
	if props.ID != "" {
		attrs = append(attrs, html.ID(props.ID))
	}
	
	if props.Disabled {
		attrs = append(attrs, html.Disabled())
	}

	// Build legend if provided
	var legendNode g.Node
	if props.Legend != "" {
		legendNode = html.Legend(
			CSSClass("fieldset-legend", "text-sm", "font-semibold", "text-gray-900", "px-2", "-mx-1"),
			g.Text(props.Legend),
		)
	}

	return html.FieldSet(
		append(attrs,
			legendNode,
			html.Div(
				CSSClass("fieldset-content", "space-y-4"),
				g.Group(children),
			),
		)...,
	)
}

// FieldSetGroup creates a group of related form controls within a fieldset
func FieldSetGroup(title string, children ...g.Node) g.Node {
	return html.Div(
		CSSClass("fieldset-group", "space-y-3"),
		g.If(title != "",
			html.Div(
				CSSClass("fieldset-group-title", "text-sm", "font-medium", "text-gray-900", "mb-2"),
				g.Text(title),
			),
		),
		html.Div(
			CSSClass("fieldset-group-content", "space-y-2"),
			g.Group(children),
		),
	)
}

// FieldSetStyles provides CSS for fieldset components
const FieldSetStyles = `
.fieldset {
	position: relative;
	margin: 0;
	padding: 1rem;
	border: 1px solid #d1d5db;
	border-radius: 0.375rem;
}

.fieldset-legend {
	position: relative;
	background: white;
	padding: 0 0.5rem;
	margin: 0 -0.25rem;
	font-weight: 600;
	color: #111827;
}

.fieldset:disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

.fieldset:disabled * {
	pointer-events: none;
}

.fieldset:focus-within {
	border-color: #3b82f6;
	box-shadow: 0 0 0 1px #3b82f6;
}

.fieldset-content {
	margin-top: 0.75rem;
}

.fieldset-group {
	border: none;
	margin: 0;
	padding: 0;
}

.fieldset-group-title {
	margin-bottom: 0.5rem;
}

.fieldset-group-content {
	padding-left: 0.5rem;
	border-left: 2px solid #e5e7eb;
}

/* Responsive fieldset */
@media (max-width: 640px) {
	.fieldset {
		padding: 0.75rem;
	}
	
	.fieldset-legend {
		font-size: 0.75rem;
	}
}

/* Dark mode support for fieldset */
@media (prefers-color-scheme: dark) {
	.fieldset {
		border-color: #374151;
		background-color: #1f2937;
	}
	
	.fieldset-legend {
		background: #1f2937;
		color: #f9fafb;
	}
}
`