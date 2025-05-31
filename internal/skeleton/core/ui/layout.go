package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Form is the core Form type that will be imported from parent core package
type Form interface {
	HasError(field string) bool
	GetError(field string) string
	Get(key string) string
}

// Layout renders the full HTML page structure.
// It includes common head elements and wraps the body content.
// This replaces core.Layout() with the same HTML structure and Tailwind CDN.
func Layout(pageTitle string, bodyContent ...g.Node) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text(pageTitle)),

			// Link to Tailwind CSS (via CDN for simplicity)
			html.Script(
				html.Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"),
				html.Async(),
				html.Defer(),
			),
		),
		html.Body(
			CSSClass("bg-gray-100", "font-sans", "antialiased"),

			// Main content area
			html.Main(
				g.Group(bodyContent),
			),
		),
	)
}

// Page renders a page container with title and children.
// This replaces core.Page() with container, title, and children pattern.
func Page(title string, children ...g.Node) g.Node {
	return html.Div(
		CSSClass("container", "mx-auto", "p-4"),
		html.H1(CSSClass("text-2xl", "font-bold", "mb-4"), g.Text(title)),
		g.Group(children),
	)
}

// FormGroupWithValidation creates a form group using the existing FormGroup component
// and integrates it with core.Form validation. It uses form.HasError(), form.GetError(), 
// and form.Get() methods directly to provide validation state and error messages.
func FormGroupWithValidation(form Form, fieldName, label string, input g.Node, helpText ...string) g.Node {
	var helpTextStr string
	if len(helpText) > 0 {
		helpTextStr = helpText[0]
	}

	var errorText string
	if form.HasError(fieldName) {
		errorText = form.GetError(fieldName)
	}

	return FormGroup(FormGroupProps{
		Label:     label,
		HelpText:  helpTextStr,
		ErrorText: errorText,
		Required:  false, // This could be a parameter if needed
		ID:        fieldName,
	}, input)
}

// TextInputWithValidation creates a TextInput with validation state from core.Form
func TextInputWithValidation(form Form, props TextInputProps) g.Node {
	// Set the value from the form if not already set
	if props.Value == "" {
		props.Value = form.Get(props.Name)
	}

	// Set validation state based on form errors
	if props.ValidationState == "" {
		if form.HasError(props.Name) {
			props.ValidationState = "invalid"
		}
	}

	// Set ID to Name if not provided (for label association)
	if props.ID == "" && props.Name != "" {
		props.ID = props.Name
	}

	return TextInput(props)
}

// TextAreaWithValidation creates a TextArea with validation state from core.Form
func TextAreaWithValidation(form Form, props TextAreaProps) g.Node {
	// Set the value from the form if not already set
	if props.Value == "" {
		props.Value = form.Get(props.Name)
	}

	// Set validation state based on form errors
	if props.ValidationState == "" {
		if form.HasError(props.Name) {
			props.ValidationState = "invalid"
		}
	}

	// Set ID to Name if not provided (for label association)
	if props.ID == "" && props.Name != "" {
		props.ID = props.Name
	}

	return TextArea(props)
}

// SelectWithValidation creates a Select with validation state from core.Form
func SelectWithValidation(form Form, props SelectProps) g.Node {
	// Set the value from the form if not already set
	if props.Value == "" {
		props.Value = form.Get(props.Name)
	}

	// Set validation state based on form errors
	if props.ValidationState == "" {
		if form.HasError(props.Name) {
			props.ValidationState = "invalid"
		}
	}

	// Set ID to Name if not provided (for label association)
	if props.ID == "" && props.Name != "" {
		props.ID = props.Name
	}

	return Select(props)
}

// FormError renders error messages for a specific field from a core.Form instance.
// Returns nil if there's no error for the field.
func FormError(form Form, field string) g.Node {
	if !form.HasError(field) {
		return nil
	}
	return html.Span(
		CSSClass("text-red-600", "text-sm", "mt-1"),
		g.Text(form.GetError(field)),
	)
}

// CSRFInput renders a hidden input field for CSRF token protection.
func CSRFInput(token string) g.Node {
	return html.Input(
		html.Type("hidden"),
		html.Name("csrf_token"),
		html.Value(token),
	)
}

// Example usage:
//
// Creating a complete form with validation integration:
//
//   func MyFormHandler(form *core.Form) g.Node {
//     return ui.Layout("My Form", 
//       ui.Page("Contact Form",
//         html.Form(
//           html.Method("POST"),
//           ui.CSRFInput("csrf-token-here"),
//           
//           // Text input with integrated validation
//           ui.FormGroupWithValidation(form, "name", "Full Name",
//             ui.TextInputWithValidation(form, ui.TextInputProps{
//               Name: "name",
//               Type: "text",
//               Placeholder: "Enter your full name",
//               Required: true,
//             }),
//             "Please enter your full name",
//           ),
//           
//           // Email input with validation
//           ui.FormGroupWithValidation(form, "email", "Email Address",
//             ui.TextInputWithValidation(form, ui.TextInputProps{
//               Name: "email", 
//               Type: "email",
//               Placeholder: "your.email@example.com",
//               Required: true,
//             }),
//           ),
//           
//           // Textarea with validation
//           ui.FormGroupWithValidation(form, "message", "Message",
//             ui.TextAreaWithValidation(form, ui.TextAreaProps{
//               Name: "message",
//               Placeholder: "Enter your message...",
//               Rows: 4,
//               Required: true,
//             }),
//           ),
//           
//           // Submit button
//           ui.Button(ui.ButtonProps{
//             Type: "submit",
//             Variant: "primary",
//           }, g.Text("Send Message")),
//         ),
//       ),
//     )
//   }