package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// FormGroupProps defines the properties for a form group component
type FormGroupProps struct {
	Label     string // The label text for the form group
	HelpText  string // Optional help text shown below the input
	ErrorText string // Error message shown when validation fails
	Required  bool   // Whether the field is required
	ID        string // ID for the input element (used for label association)
}

// FormGroup creates a form group with label, input wrapper, help text, and error text
func FormGroup(props FormGroupProps, children ...g.Node) g.Node {
	// Build label attributes
	var labelAttrs []g.Node
	labelAttrs = append(labelAttrs, CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"))
	
	if props.ID != "" {
		labelAttrs = append(labelAttrs, html.For(props.ID))
	}
	
	// Add required indicator classes
	if props.Required {
		labelAttrs = append(labelAttrs, CSSClass("required"))
	}

	// Determine if we're in error state
	hasError := props.ErrorText != ""

	return html.Div(
		CSSClass("form-group", "space-y-1"),
		
		// Label with optional required indicator
		g.If(props.Label != "", 
			html.Label(
				append(labelAttrs, g.Text(props.Label))...,
			),
		),
		
		// Input wrapper with error state styling
		html.Div(
			CSSClass("relative"),
			g.Group(children),
		),
		
		// Help text or error text
		g.If(hasError,
			html.Div(
				CSSClass("text-sm", "text-red-600", "mt-1"),
				html.Span(
					html.Aria("live", "polite"),
					html.Role("alert"),
					g.Text(props.ErrorText),
				),
			),
		),
		g.If(!hasError && props.HelpText != "",
			html.Div(
				CSSClass("text-sm", "text-gray-500", "mt-1"),
				g.Text(props.HelpText),
			),
		),
	)
}

// FormGroupStyles provides CSS for form group components
const FormGroupStyles = `
.form-group {
	margin-bottom: 1rem;
}

.form-group .required::after {
	content: ' *';
	color: #ef4444;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
	outline: none;
	ring: 2px;
	ring-color: #3b82f6;
	border-color: #3b82f6;
}

.form-group input:invalid,
.form-group textarea:invalid,
.form-group select:invalid {
	border-color: #ef4444;
	background-color: #fef2f2;
}

.form-group input:valid,
.form-group textarea:valid,
.form-group select:valid {
	border-color: #10b981;
	background-color: #f0fdf4;
}

/* Error state styling */
.form-group:has([aria-invalid="true"]) input,
.form-group:has([aria-invalid="true"]) textarea,
.form-group:has([aria-invalid="true"]) select {
	border-color: #ef4444;
	background-color: #fef2f2;
}

/* Focus within form group */
.form-group:focus-within label {
	color: #3b82f6;
}
`