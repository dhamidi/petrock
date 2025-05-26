package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleFormInputsDetail handles the form inputs component detail page
func HandleFormInputsDetail(w http.ResponseWriter, r *http.Request) {
	// Create demo content showing different form input components
	demoContent := html.Div(
		ui.CSSClass("space-y-8"),
		
		// Header section
		html.Div(
			ui.CSSClass("mb-8"),
			html.H1(
				ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
				g.Text("Form Inputs Component"),
			),
			html.P(
				ui.CSSClass("text-lg", "text-gray-600"),
				g.Text("Essential form input components including text inputs, textareas, and select dropdowns with validation states."),
			),
		),

		// Text Input Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Text Input"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("Text inputs support various types, validation states, and accessibility features."),
			),
			
			// Basic inputs
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Basic Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-4"),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Email Address"),
						),
						ui.TextInput(ui.TextInputProps{
							Type:        "email",
							Placeholder: "Enter your email",
							ID:          "email-basic",
						}),
					),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Password"),
						),
						ui.TextInput(ui.TextInputProps{
							Type:        "password",
							Placeholder: "Enter your password",
							ID:          "password-basic",
							Required:    true,
						}),
					),
				),
			),
			
			// Validation states
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Validation States"),
				),
				html.Div(
					ui.CSSClass("space-y-4"),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Valid Input"),
						),
						ui.TextInput(ui.TextInputProps{
							Type:            "text",
							Value:           "john@example.com",
							ValidationState: "valid",
							ID:              "valid-input",
						}),
						html.Div(
							ui.CSSClass("text-sm", "text-green-600", "mt-1"),
							g.Text("✓ Email format is correct"),
						),
					),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Invalid Input"),
						),
						ui.TextInput(ui.TextInputProps{
							Type:            "email",
							Value:           "invalid-email",
							ValidationState: "invalid",
							ID:              "invalid-input",
						}),
						html.Div(
							ui.CSSClass("text-sm", "text-red-600", "mt-1"),
							g.Text("Please enter a valid email address"),
						),
					),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Disabled Input"),
						),
						ui.TextInput(ui.TextInputProps{
							Type:     "text",
							Value:    "Disabled field",
							Disabled: true,
							ID:       "disabled-input",
						}),
					),
				),
			),
		),

		// TextArea Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("TextArea"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("Multi-line text inputs with configurable rows and resize behavior."),
			),
			
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("TextArea Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-4"),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Description"),
						),
						ui.TextArea(ui.TextAreaProps{
							Placeholder: "Enter a description",
							Rows:        3,
							ID:          "description",
						}),
					),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Comments (No Resize)"),
						),
						ui.TextArea(ui.TextAreaProps{
							Placeholder: "Enter your comments",
							Rows:        4,
							Resize:      "none",
							ID:          "comments",
						}),
					),
				),
			),
		),

		// Select Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Select"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("Dropdown select components with options and validation states."),
			),
			
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Select Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-4"),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Country"),
						),
						ui.Select(ui.SelectProps{
							Placeholder: "Select a country",
							Options: []ui.SelectOption{
								{Value: "us", Label: "United States"},
								{Value: "ca", Label: "Canada"},
								{Value: "uk", Label: "United Kingdom"},
								{Value: "de", Label: "Germany"},
								{Value: "fr", Label: "France"},
							},
							ID: "country",
						}),
					),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Priority"),
						),
						ui.Select(ui.SelectProps{
							Value: "medium",
							Options: []ui.SelectOption{
								{Value: "low", Label: "Low Priority"},
								{Value: "medium", Label: "Medium Priority"},
								{Value: "high", Label: "High Priority"},
								{Value: "urgent", Label: "Urgent"},
							},
							Required: true,
							ID:       "priority",
						}),
					),
					html.Div(
						html.Label(
							ui.CSSClass("block", "text-sm", "font-medium", "text-gray-700", "mb-1"),
							g.Text("Multi-Select Tags"),
						),
						ui.Select(ui.SelectProps{
							Options: []ui.SelectOption{
								{Value: "react", Label: "React", Selected: true},
								{Value: "vue", Label: "Vue.js"},
								{Value: "angular", Label: "Angular"},
								{Value: "svelte", Label: "Svelte", Selected: true},
								{Value: "typescript", Label: "TypeScript"},
								{Value: "javascript", Label: "JavaScript"},
							},
							Multiple: true,
							ID:       "tags",
						}),
						html.Div(
							ui.CSSClass("text-sm", "text-gray-500", "mt-1"),
							g.Text("Hold Ctrl/Cmd to select multiple options"),
						),
					),
				),
			),
		),
	)

	// Create page content with proper sidebar navigation
	pageContent := core.Page("Form Inputs Component",
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
		"Form Inputs Component - UI Gallery",
		pageContent,
	)

	w.Header().Set("Content-Type", "text/html")
	response.Render(w)
}