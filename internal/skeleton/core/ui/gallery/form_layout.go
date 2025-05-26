package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleFormLayoutDetail handles the form layout component detail page
func HandleFormLayoutDetail(w http.ResponseWriter, r *http.Request) {
	// Create demo content showing form layout components
	demoContent := html.Div(
		ui.CSSClass("space-y-8"),
		
		// Header section
		html.Div(
			ui.CSSClass("mb-8"),
			html.H1(
				ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
				g.Text("Form Layout Components"),
			),
			html.P(
				ui.CSSClass("text-lg", "text-gray-600"),
				g.Text("Form layout components including FormGroup and FieldSet for organizing form elements with proper labeling, validation, and accessibility."),
			),
		),

		// FormGroup Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("FormGroup"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("FormGroup combines labels, inputs, help text, and error messages into a cohesive unit with proper accessibility."),
			),
			
			// Basic FormGroup examples
			html.Div(
				ui.CSSClass("space-y-6", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Basic FormGroup Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-4", "max-w-md"),
					
					// Basic form group with text input
					ui.FormGroup(ui.FormGroupProps{
						Label:    "Email Address",
						HelpText: "We'll never share your email with anyone else.",
						Required: true,
						ID:       "email-basic",
					},
						ui.TextInput(ui.TextInputProps{
							Type:        "email",
							Placeholder: "Enter your email",
							ID:          "email-basic",
						}),
					),
					
					// Form group with textarea
					ui.FormGroup(ui.FormGroupProps{
						Label:    "Description",
						HelpText: "Provide a brief description of your request.",
						ID:       "description-basic",
					},
						ui.TextArea(ui.TextAreaProps{
							Placeholder: "Enter description",
							Rows:        3,
							ID:          "description-basic",
						}),
					),
					
					// Form group with select
					ui.FormGroup(ui.FormGroupProps{
						Label:    "Priority Level",
						Required: true,
						ID:       "priority-basic",
					},
						ui.Select(ui.SelectProps{
							Placeholder: "Select priority",
							Options: []ui.SelectOption{
								{Value: "low", Label: "Low"},
								{Value: "medium", Label: "Medium"},
								{Value: "high", Label: "High"},
								{Value: "urgent", Label: "Urgent"},
							},
							ID: "priority-basic",
						}),
					),
				),
			),
			
			// Error state examples
			html.Div(
				ui.CSSClass("space-y-6", "p-6", "bg-red-50", "rounded-lg", "border", "border-red-200"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4", "text-red-900"),
					g.Text("FormGroup with Validation Errors"),
				),
				html.Div(
					ui.CSSClass("space-y-4", "max-w-md"),
					
					// Form group with error
					ui.FormGroup(ui.FormGroupProps{
						Label:     "Email Address",
						ErrorText: "Please enter a valid email address.",
						Required:  true,
						ID:        "email-error",
					},
						ui.TextInput(ui.TextInputProps{
							Type:            "email",
							Value:           "invalid-email",
							ValidationState: "invalid",
							ID:              "email-error",
						}),
					),
					
					// Form group with password error
					ui.FormGroup(ui.FormGroupProps{
						Label:     "Password",
						ErrorText: "Password must be at least 8 characters long.",
						Required:  true,
						ID:        "password-error",
					},
						ui.TextInput(ui.TextInputProps{
							Type:            "password",
							Value:           "123",
							ValidationState: "invalid",
							ID:              "password-error",
						}),
					),
				),
			),
		),

		// FieldSet Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("FieldSet"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("FieldSet groups related form controls with a legend and provides semantic organization for complex forms."),
			),
			
			// Basic FieldSet example
			html.Div(
				ui.CSSClass("space-y-6", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Contact Information FieldSet"),
				),
				ui.FieldSet(ui.FieldSetProps{
					Legend: "Contact Information",
					ID:     "contact-info",
				},
					ui.FormGroup(ui.FormGroupProps{
						Label:    "First Name",
						Required: true,
						ID:       "first-name",
					},
						ui.TextInput(ui.TextInputProps{
							Type:        "text",
							Placeholder: "Enter first name",
							ID:          "first-name",
						}),
					),
					ui.FormGroup(ui.FormGroupProps{
						Label:    "Last Name",
						Required: true,
						ID:       "last-name",
					},
						ui.TextInput(ui.TextInputProps{
							Type:        "text",
							Placeholder: "Enter last name",
							ID:          "last-name",
						}),
					),
					ui.FormGroup(ui.FormGroupProps{
						Label: "Phone Number",
						ID:    "phone",
					},
						ui.TextInput(ui.TextInputProps{
							Type:        "tel",
							Placeholder: "Enter phone number",
							ID:          "phone",
						}),
					),
				),
			),
			
			// Disabled FieldSet example
			html.Div(
				ui.CSSClass("space-y-6", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Disabled FieldSet"),
				),
				ui.FieldSet(ui.FieldSetProps{
					Legend:   "Billing Address (Disabled)",
					Disabled: true,
					ID:       "billing-disabled",
				},
					ui.FormGroup(ui.FormGroupProps{
						Label: "Street Address",
						ID:    "street-disabled",
					},
						ui.TextInput(ui.TextInputProps{
							Type:        "text",
							Value:       "123 Main St",
							ID:          "street-disabled",
						}),
					),
					ui.FormGroup(ui.FormGroupProps{
						Label: "City",
						ID:    "city-disabled",
					},
						ui.TextInput(ui.TextInputProps{
							Type:        "text",
							Value:       "Anytown",
							ID:          "city-disabled",
						}),
					),
				),
			),
			
			// Nested FieldSet example
			html.Div(
				ui.CSSClass("space-y-6", "p-6", "bg-blue-50", "rounded-lg", "border", "border-blue-200"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4", "text-blue-900"),
					g.Text("Complex Form with Multiple FieldSets"),
				),
				html.Form(
					ui.CSSClass("space-y-6"),
					
					// Personal Information FieldSet
					ui.FieldSet(ui.FieldSetProps{
						Legend: "Personal Information",
						ID:     "personal-info",
					},
						html.Div(
							ui.CSSClass("grid", "grid-cols-1", "md:grid-cols-2", "gap-4"),
							ui.FormGroup(ui.FormGroupProps{
								Label:    "First Name",
								Required: true,
								ID:       "first-name-complex",
							},
								ui.TextInput(ui.TextInputProps{
									Type:        "text",
									Placeholder: "First name",
									ID:          "first-name-complex",
								}),
							),
							ui.FormGroup(ui.FormGroupProps{
								Label:    "Last Name",
								Required: true,
								ID:       "last-name-complex",
							},
								ui.TextInput(ui.TextInputProps{
									Type:        "text",
									Placeholder: "Last name",
									ID:          "last-name-complex",
								}),
							),
						),
						ui.FormGroup(ui.FormGroupProps{
							Label:    "Email",
							Required: true,
							ID:       "email-complex",
						},
							ui.TextInput(ui.TextInputProps{
								Type:        "email",
								Placeholder: "Email address",
								ID:          "email-complex",
							}),
						),
					),
					
					// Address Information FieldSet
					ui.FieldSet(ui.FieldSetProps{
						Legend: "Address Information",
						ID:     "address-info",
					},
						ui.FormGroup(ui.FormGroupProps{
							Label: "Street Address",
							ID:    "street-complex",
						},
							ui.TextInput(ui.TextInputProps{
								Type:        "text",
								Placeholder: "123 Main Street",
								ID:          "street-complex",
							}),
						),
						html.Div(
							ui.CSSClass("grid", "grid-cols-1", "md:grid-cols-3", "gap-4"),
							ui.FormGroup(ui.FormGroupProps{
								Label: "City",
								ID:    "city-complex",
							},
								ui.TextInput(ui.TextInputProps{
									Type:        "text",
									Placeholder: "City",
									ID:          "city-complex",
								}),
							),
							ui.FormGroup(ui.FormGroupProps{
								Label: "State",
								ID:    "state-complex",
							},
								ui.Select(ui.SelectProps{
									Placeholder: "Select state",
									Options: []ui.SelectOption{
										{Value: "ca", Label: "California"},
										{Value: "ny", Label: "New York"},
										{Value: "tx", Label: "Texas"},
									},
									ID: "state-complex",
								}),
							),
							ui.FormGroup(ui.FormGroupProps{
								Label: "ZIP Code",
								ID:    "zip-complex",
							},
								ui.TextInput(ui.TextInputProps{
									Type:        "text",
									Placeholder: "12345",
									ID:          "zip-complex",
								}),
							),
						),
					),
				),
			),
		),
	)

	// Create page content with proper sidebar navigation
	pageContent := core.Page("Form Layout Components",
		html.Div(
			ui.CSSClass("flex", "min-h-screen", "-mx-4", "-mt-4"),
			// Sidebar with full component navigation
			html.Nav(
				ui.CSSClass("w-64", "bg-white", "border-r", "border-gray-200", "p-6", "overflow-y-auto"),
				html.H1(
					ui.CSSClass("text-lg", "font-semibold", "text-gray-900", "mb-6"),
					g.Text("Components"),
				),
				g.Group(BuildSidebar()),
			),
			// Main content
			html.Main(
				ui.CSSClass("flex-1", "p-6", "overflow-y-auto"),
				html.Div(
					ui.CSSClass("max-w-4xl"),
					demoContent,
				),
			),
		),
	)

	// Use existing Layout function
	response := core.Layout(
		"Form Layout Components - UI Gallery",
		pageContent,
	)

	w.Header().Set("Content-Type", "text/html")
	response.Render(w)
}