package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleButtonDetail handles the Button component demo page
func HandleButtonDetail(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different button variants and states
		demoContent := html.Div(
			ui.CSSClass("space-y-8"),
			
			// Header section
			html.Div(
				ui.CSSClass("mb-8"),
				html.H1(
					ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
					g.Text("Button Component"),
				),
				html.P(
					ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
					g.Text("The Button component provides interactive elements with multiple variants, sizes, and states for consistent user interactions."),
				),
			),
			
			// Button Variants
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Button Variants"),
				),
				
				// Primary buttons
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Primary Buttons"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-4", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Primary")),
						ui.Button(ui.ButtonProps{Variant: "primary", Disabled: true}, g.Text("Disabled")),
					),
				),
				
				// Secondary buttons
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Secondary Buttons"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-4", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Secondary")),
						ui.Button(ui.ButtonProps{Variant: "secondary", Disabled: true}, g.Text("Disabled")),
					),
				),
				
				// Danger buttons
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Danger Buttons"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-4", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Variant: "danger"}, g.Text("Delete")),
						ui.Button(ui.ButtonProps{Variant: "danger", Disabled: true}, g.Text("Disabled")),
					),
				),
				
				// Link buttons
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Link Buttons"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-4", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Variant: "link"}, g.Text("Link Button")),
						ui.Button(ui.ButtonProps{Variant: "link", Disabled: true}, g.Text("Disabled")),
					),
				),
			),
			
			// Button Sizes
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Button Sizes"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Buttons come in three sizes to fit different contexts:"),
				),
				html.Div(
					ui.CSSClass("flex", "flex-wrap", "items-center", "gap-4", "p-4", "border", "rounded", "bg-gray-50"),
					ui.Button(ui.ButtonProps{Size: "small"}, g.Text("Small")),
					ui.Button(ui.ButtonProps{Size: "medium"}, g.Text("Medium")),
					ui.Button(ui.ButtonProps{Size: "large"}, g.Text("Large")),
				),
			),
			
			// Button Types
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Button Types"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Buttons support different HTML types for forms and interactions:"),
				),
				html.Div(
					ui.CSSClass("flex", "flex-wrap", "gap-4", "p-4", "border", "rounded", "bg-gray-50"),
					ui.Button(ui.ButtonProps{Type: "button"}, g.Text("Button")),
					ui.Button(ui.ButtonProps{Type: "submit", Variant: "primary"}, g.Text("Submit")),
					ui.Button(ui.ButtonProps{Type: "reset", Variant: "secondary"}, g.Text("Reset")),
				),
			),
			
			// Interactive Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Interactive Examples"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Common button combinations and use cases:"),
				),
				
				// Action group
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Action Group"),
					),
					html.Div(
						ui.CSSClass("flex", "gap-3", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Save")),
						ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Cancel")),
					),
				),
				
				// Destructive action
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Destructive Action"),
					),
					html.Div(
						ui.CSSClass("flex", "gap-3", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Variant: "danger"}, g.Text("Delete Account")),
						ui.Button(ui.ButtonProps{Variant: "link"}, g.Text("Never mind")),
					),
				),
				
				// Form buttons
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Form Actions"),
					),
					html.Div(
						ui.CSSClass("flex", "gap-3", "p-4", "border", "rounded", "bg-gray-50"),
						ui.Button(ui.ButtonProps{Type: "submit", Variant: "primary"}, g.Text("Submit Form")),
						ui.Button(ui.ButtonProps{Type: "reset", Variant: "secondary"}, g.Text("Reset")),
						ui.Button(ui.ButtonProps{Type: "button", Variant: "link"}, g.Text("Back")),
					),
				),
			),
			
			// Code Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Usage Examples"),
				),
				html.Pre(
					ui.CSSClass("bg-gray-100", "p-4", "rounded", "text-sm", "overflow-x-auto"),
					g.Text(`// Basic button with default settings
ui.Button(ui.ButtonProps{}, g.Text("Click me"))

// Primary button with specific size
ui.Button(ui.ButtonProps{
    Variant: "primary",
    Size: "large",
}, g.Text("Save Changes"))

// Disabled secondary button
ui.Button(ui.ButtonProps{
    Variant: "secondary",
    Disabled: true,
}, g.Text("Not Available"))

// Submit button for forms
ui.Button(ui.ButtonProps{
    Type: "submit",
    Variant: "primary",
}, g.Text("Submit Form"))

// Danger button for destructive actions
ui.Button(ui.ButtonProps{
    Variant: "danger",
    Size: "small",
}, g.Text("Delete"))`),
				),
			),
			
			// Properties documentation
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Properties"),
				),
				html.Div(
					ui.CSSClass("border", "rounded", "overflow-hidden"),
					html.Table(
						ui.CSSClass("w-full"),
						html.THead(
							ui.CSSClass("bg-gray-50"),
							html.Tr(
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Property")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Type")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Default")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Description")),
							),
						),
						html.TBody(
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Variant")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"primary\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Button style: primary, secondary, danger, link")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Size")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"medium\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Button size: small, medium, large")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Type")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"button\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("HTML button type: button, submit, reset")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Disabled")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("bool")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("false")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Whether the button is disabled")),
							),
						),
					),
				),
			),
		)

		// Create page content with proper sidebar navigation
		pageContent := core.Page("Button Component",
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
			"Button Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
}