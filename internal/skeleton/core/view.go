package core

import (
	g "maragu.dev/gomponents" // Use canonical import path
	. "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"
)

// --- Basic Elements with Tailwind ---

// Button renders a styled button.
func Button(text string, attrs ...g.Node) g.Node {
	return html.Button(
		// Example Tailwind classes - customize as needed
		Classes{ // Correct map literal syntax
			"py-2": true, "px-4": true, "bg-blue-500": true, "text-white": true, "font-semibold": true,
			"rounded-lg": true, "shadow-md": true, "hover:bg-blue-700": true, "focus:outline-none": true,
			"focus:ring-2": true, "focus:ring-blue-400": true, "focus:ring-opacity-75": true,
		},
		g.Text(text),
		g.Group(attrs), // g.Group is correct here as attrs is already a slice
	)
}

// Input renders a styled input field.
func Input(inputType, name, value string, attrs ...g.Node) g.Node {
	return html.Input(
		html.Type(inputType),
		html.Name(name),
		html.Value(value),
		// Example Tailwind classes
		Classes{ // Correct map literal syntax
			"mt-1": true, "block": true, "w-full": true, "rounded-md": true, "border-gray-300": true,
			"shadow-sm": true, "focus:border-indigo-300": true, "focus:ring": true,
			"focus:ring-indigo-200": true, "focus:ring-opacity-50": true,
		},
		g.Group(attrs), // g.Group is correct here as attrs is already a slice
	)
}

// TextArea renders a styled textarea.
func TextArea(name, value string, attrs ...g.Node) g.Node {
	return html.Textarea(
		html.Name(name),
		// Example Tailwind classes
		Classes{ // Correct map literal syntax
			"mt-1": true, "block": true, "w-full": true, "rounded-md": true, "border-gray-300": true,
			"shadow-sm": true, "focus:border-indigo-300": true, "focus:ring": true,
			"focus:ring-indigo-200": true, "focus:ring-opacity-50": true,
		},
		g.Group(attrs), // g.Group is correct here as attrs is already a slice
		g.Text(value),  // Text content for textarea
	)
}

// Select renders a styled select dropdown.
func Select(name string, options map[string]string, selectedValue string, attrs ...g.Node) g.Node {
	opts := make([]g.Node, 0, len(options))
	for val, text := range options {
		opt := html.Option(html.Value(val), g.Text(text))
		if val == selectedValue {
			// Correctly pass a slice to g.Group
			opt = g.Group([]g.Node{opt, html.Selected()})
		}
		opts = append(opts, opt)
	}

	return html.Select(
		html.Name(name),
		// Example Tailwind classes
		Classes{ // Correct map literal syntax
			"mt-1": true, "block": true, "w-full": true, "rounded-md": true, "border-gray-300": true,
			"shadow-sm": true, "focus:border-indigo-300": true, "focus:ring": true,
			"focus:ring-indigo-200": true, "focus:ring-opacity-50": true,
		},
		g.Group(attrs), // g.Group is correct here as attrs is already a slice
		g.Group(opts),  // g.Group is correct here as opts is already a slice
	)
}

// --- Form Handling ---

// FormError renders error messages for a specific field from a core.Form instance.
// Returns nil if there's no error for the field.
func FormError(form *Form, field string) g.Node {
	if !form.HasError(field) {
		return nil
	}
	// Example Tailwind classes for error message
	return html.Span(Classes{"text-red-600": true, "text-sm": true, "mt-1": true}, g.Text(form.GetError(field))) // Correct map literal syntax
}

// CSRFTokenInput renders a hidden input field for CSRF token protection.
// Assumes you have a way to get the current CSRF token.
func CSRFTokenInput(token string) g.Node {
	// The name "csrf_token" is common, adjust if your CSRF library expects differently.
	return html.Input(html.Type("hidden"), html.Name("csrf_token"), html.Value(token))
}

// FieldGroup combines a label, input/textarea/select, and error message.
func FieldGroup(label, fieldName string, inputElement g.Node, form *Form) g.Node {
	return html.Div(
		Classes{"mb-4": true}, // Correct map literal syntax
		html.Label(
			html.For(fieldName),
			Classes{"block": true, "text-sm": true, "font-medium": true, "text-gray-700": true}, // Correct map literal syntax
			g.Text(label),
		),
		inputElement,
		FormError(form, fieldName), // Display error if present
	)
}

// --- Page Structure ---

// Page component (can be used within Layout)
func Page(title string, children ...g.Node) g.Node {
	return html.Div(
		Classes{"container": true, "mx-auto": true, "p-4": true},                           // Correct map literal syntax
		html.H1(Classes{"text-2xl": true, "font-bold": true, "mb-4": true}, g.Text(title)), // Correct map literal syntax
		g.Group(children), // g.Group is correct here as children is already a slice
	)
}

// --- Asset Handling ---

// StylesheetLink creates a <link> tag for a CSS stylesheet.
func StylesheetLink(href string) g.Node {
	return html.Link(html.Rel("stylesheet"), html.Href(href))
}

// ScriptLink creates a <script> tag for a JavaScript file.
// useDefer indicates if the 'defer' attribute should be added.
func ScriptLink(src string, async bool, useDefer bool) g.Node {
	attrs := []g.Node{html.Src(src)}
	if async {
		attrs = append(attrs, html.Async())
	}
	if useDefer {
		attrs = append(attrs, DeferAttr()) // Use renamed helper function
	}
	return html.Script(attrs...)
}

// --- Helper Components (from gomponents/components) ---

// DeferAttr returns a gomponents node that adds the defer attribute.
func DeferAttr() g.Node {
	return g.Attr("defer")
}
