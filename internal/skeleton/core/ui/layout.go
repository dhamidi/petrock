package ui

import (
	"net/url"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ParseError represents a single validation error
type ParseError struct {
	Field   string
	Message string
	Code    string
	Meta    map[string]interface{}
}

// FormData holds form values and validation errors for template rendering
type FormData struct {
	Values url.Values
	Errors []ParseError
}

// HasError checks if there are any errors for the given field
func (f *FormData) HasError(field string) bool {
	for _, err := range f.Errors {
		if err.Field == field {
			return true
		}
	}
	return false
}

// GetError returns the first error message for the given field
func (f *FormData) GetError(field string) string {
	for _, err := range f.Errors {
		if err.Field == field {
			return err.Message
		}
	}
	return ""
}

// Get returns the value for the given field from form values
func (f *FormData) Get(key string) string {
	if f.Values == nil {
		return ""
	}
	return f.Values.Get(key)
}

// NewFormData creates a new FormData instance
func NewFormData(values url.Values, errors []ParseError) *FormData {
	return &FormData{
		Values: values,
		Errors: errors,
	}
}



// Legacy Form interface for backward compatibility during migration
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

// FormGroupWithValidation creates a form group using FormData with validation state
func FormGroupWithValidation(formData *FormData, fieldName, label string, input g.Node, helpText ...string) g.Node {
	var helpTextStr string
	if len(helpText) > 0 {
		helpTextStr = helpText[0]
	}

	var errorText string
	if formData.HasError(fieldName) {
		errorText = formData.GetError(fieldName)
	}

	return FormGroup(FormGroupProps{
		Label:     label,
		HelpText:  helpTextStr,
		ErrorText: errorText,
		Required:  false, // This could be a parameter if needed
		ID:        fieldName,
	}, input)
}

// FormGroupWithValidationLegacy provides backward compatibility with legacy Form interface
func FormGroupWithValidationLegacy(form Form, fieldName, label string, input g.Node, helpText ...string) g.Node {
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

// TextInputWithValidation creates a TextInput with validation state from FormData
func TextInputWithValidation(formData *FormData, props TextInputProps) g.Node {
	// Set the value from the form if not already set
	if props.Value == "" {
		props.Value = formData.Get(props.Name)
	}

	// Set validation state based on form errors
	if props.ValidationState == "" {
		if formData.HasError(props.Name) {
			props.ValidationState = "invalid"
		}
	}

	// Set ID to Name if not provided (for label association)
	if props.ID == "" && props.Name != "" {
		props.ID = props.Name
	}

	return TextInput(props)
}

// TextInputWithValidationLegacy provides backward compatibility with legacy Form interface
func TextInputWithValidationLegacy(form Form, props TextInputProps) g.Node {
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

// TextAreaWithValidation creates a TextArea with validation state from FormData
func TextAreaWithValidation(formData *FormData, props TextAreaProps) g.Node {
	// Set the value from the form if not already set
	if props.Value == "" {
		props.Value = formData.Get(props.Name)
	}

	// Set validation state based on form errors
	if props.ValidationState == "" {
		if formData.HasError(props.Name) {
			props.ValidationState = "invalid"
		}
	}

	// Set ID to Name if not provided (for label association)
	if props.ID == "" && props.Name != "" {
		props.ID = props.Name
	}

	return TextArea(props)
}

// TextAreaWithValidationLegacy provides backward compatibility with legacy Form interface
func TextAreaWithValidationLegacy(form Form, props TextAreaProps) g.Node {
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

// SelectWithValidation creates a Select with validation state from FormData
func SelectWithValidation(formData *FormData, props SelectProps) g.Node {
	// Set the value from the form if not already set
	if props.Value == "" {
		props.Value = formData.Get(props.Name)
	}

	// Set validation state based on form errors
	if props.ValidationState == "" {
		if formData.HasError(props.Name) {
			props.ValidationState = "invalid"
		}
	}

	// Set ID to Name if not provided (for label association)
	if props.ID == "" && props.Name != "" {
		props.ID = props.Name
	}

	return Select(props)
}

// SelectWithValidationLegacy provides backward compatibility with legacy Form interface
func SelectWithValidationLegacy(form Form, props SelectProps) g.Node {
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

// FormError renders error messages for a specific field from FormData.
// Returns nil if there's no error for the field.
func FormError(formData *FormData, field string) g.Node {
	if !formData.HasError(field) {
		return nil
	}
	return html.Span(
		CSSClass("text-red-600", "text-sm", "mt-1"),
		g.Text(formData.GetError(field)),
	)
}

// FormErrorLegacy provides backward compatibility with legacy Form interface
func FormErrorLegacy(form Form, field string) g.Node {
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
//   func MyFormHandler(formData *ui.FormData) g.Node {
//     return ui.Layout("My Form", 
//       ui.Page("Contact Form",
//         html.Form(
//           html.Method("POST"),
//           ui.CSRFInput("csrf-token-here"),
//           
//           // Text input with integrated validation
//           ui.FormGroupWithValidation(formData, "name", "Full Name",
//             ui.TextInputWithValidation(formData, ui.TextInputProps{
//               Name: "name",
//               Type: "text",
//               Placeholder: "Enter your full name",
//               Required: true,
//             }),
//             "Please enter your full name",
//           ),
//           
//           // Email input with validation
//           ui.FormGroupWithValidation(formData, "email", "Email Address",
//             ui.TextInputWithValidation(formData, ui.TextInputProps{
//               Name: "email", 
//               Type: "email",
//               Placeholder: "your.email@example.com",
//               Required: true,
//             }),
//           ),
//           
//           // Textarea with validation
//           ui.FormGroupWithValidation(formData, "message", "Message",
//             ui.TextAreaWithValidation(formData, ui.TextAreaProps{
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
//
// Legacy compatibility example:
//
//   func MyLegacyFormHandler(form *core.Form) g.Node {
//     return ui.Layout("My Form", 
//       ui.Page("Contact Form",
//         html.Form(
//           html.Method("POST"),
//           ui.CSRFInput("csrf-token-here"),
//           
//           // Use legacy functions for backward compatibility
//           ui.FormGroupWithValidationLegacy(form, "name", "Full Name",
//             ui.TextInputWithValidationLegacy(form, ui.TextInputProps{
//               Name: "name",
//               Type: "text",
//               Placeholder: "Enter your full name",
//               Required: true,
//             }),
//             "Please enter your full name",
//           ),
//         ),
//       ),
//     )
//   }