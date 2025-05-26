package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleButtonGroupDetail handles the ButtonGroup component demo page
func HandleButtonGroupDetail(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different button group layouts
		demoContent := html.Div(
			ui.CSSClass("space-y-8"),
			
			// Header section
			html.Div(
				ui.CSSClass("mb-8"),
				html.H1(
					ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
					g.Text("ButtonGroup Component"),
				),
				html.P(
					ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
					g.Text("The ButtonGroup component provides a container for grouping related buttons with consistent spacing and orientation."),
				),
			),
			
			// Horizontal Button Groups
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Horizontal Button Groups"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Horizontal button groups arrange buttons side by side with consistent spacing:"),
				),
				
				// Default horizontal group
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Default Horizontal (Medium Spacing)"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal", Spacing: "medium"},
							ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Save")),
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Cancel")),
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Help")),
						),
					),
				),
				
				// Small spacing horizontal
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Small Spacing"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal", Spacing: "small"},
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("Bold")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("Italic")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("Underline")),
						),
					),
				),
				
				// Large spacing horizontal
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Large Spacing"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal", Spacing: "large"},
							ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Previous")),
							ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Next")),
						),
					),
				),
				
				// No spacing horizontal
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("No Spacing (Connected)"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal", Spacing: "none"},
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("1")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("2")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("3")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("4")),
						),
					),
				),
			),
			
			// Vertical Button Groups
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Vertical Button Groups"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Vertical button groups stack buttons vertically with consistent spacing:"),
				),
				
				html.Div(
					ui.CSSClass("grid", "grid-cols-1", "md:grid-cols-3", "gap-6"),
					
					// Default vertical group
					html.Div(
						html.H3(
							ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
							g.Text("Medium Spacing"),
						),
						html.Div(
							ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
							ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "vertical", Spacing: "medium"},
								ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("New File")),
								ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Open")),
								ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Save")),
								ui.Button(ui.ButtonProps{Variant: "danger"}, g.Text("Delete")),
							),
						),
					),
					
					// Small spacing vertical
					html.Div(
						html.H3(
							ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
							g.Text("Small Spacing"),
						),
						html.Div(
							ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
							ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "vertical", Spacing: "small"},
								ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("Option 1")),
								ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("Option 2")),
								ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("Option 3")),
							),
						),
					),
					
					// Large spacing vertical
					html.Div(
						html.H3(
							ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
							g.Text("Large Spacing"),
						),
						html.Div(
							ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
							ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "vertical", Spacing: "large"},
								ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Submit")),
								ui.Button(ui.ButtonProps{Variant: "link"}, g.Text("Cancel")),
							),
						),
					),
				),
			),
			
			// Interactive Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Common Use Cases"),
				),
				
				// Form actions
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Form Actions"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal"},
							ui.Button(ui.ButtonProps{Type: "submit", Variant: "primary"}, g.Text("Submit")),
							ui.Button(ui.ButtonProps{Type: "reset", Variant: "secondary"}, g.Text("Reset")),
							ui.Button(ui.ButtonProps{Type: "button", Variant: "link"}, g.Text("Cancel")),
						),
					),
				),
				
				// Navigation controls
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Navigation Controls"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal"},
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("← Previous")),
							ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Next →")),
						),
					),
				),
				
				// Toolbar buttons
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Toolbar Buttons"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "horizontal", Spacing: "small"},
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("📋")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("✂️")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("📋")),
							ui.Button(ui.ButtonProps{Variant: "secondary", Size: "small"}, g.Text("🗑️")),
						),
					),
				),
				
				// Sidebar actions
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Sidebar Actions"),
					),
					html.Div(
						ui.CSSClass("p-4", "border", "rounded", "bg-gray-50", "max-w-xs"),
						ui.ButtonGroup(ui.ButtonGroupProps{Orientation: "vertical"},
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("📄 Documents")),
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("📁 Projects")),
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("⚙️ Settings")),
							ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("👤 Profile")),
						),
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
					g.Text(`// Horizontal button group (default)
ui.ButtonGroup(ui.ButtonGroupProps{},
    ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Save")),
    ui.Button(ui.ButtonProps{Variant: "secondary"}, g.Text("Cancel")),
)

// Vertical button group with large spacing
ui.ButtonGroup(ui.ButtonGroupProps{
    Orientation: "vertical",
    Spacing: "large",
},
    ui.Button(ui.ButtonProps{Variant: "primary"}, g.Text("Submit")),
    ui.Button(ui.ButtonProps{Variant: "link"}, g.Text("Cancel")),
)

// Compact horizontal toolbar
ui.ButtonGroup(ui.ButtonGroupProps{
    Orientation: "horizontal",
    Spacing: "small",
},
    ui.Button(ui.ButtonProps{Size: "small"}, g.Text("Cut")),
    ui.Button(ui.ButtonProps{Size: "small"}, g.Text("Copy")),
    ui.Button(ui.ButtonProps{Size: "small"}, g.Text("Paste")),
)

// Connected buttons (no spacing)
ui.ButtonGroup(ui.ButtonGroupProps{Spacing: "none"},
    ui.Button(ui.ButtonProps{}, g.Text("1")),
    ui.Button(ui.ButtonProps{}, g.Text("2")),
    ui.Button(ui.ButtonProps{}, g.Text("3")),
)`),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Orientation")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"horizontal\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Button arrangement: horizontal, vertical")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Spacing")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"medium\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Space between buttons: none, small, medium, large")),
							),
						),
					),
				),
			),
		)

		// Create page content with proper sidebar navigation
		pageContent := core.Page("ButtonGroup Component",
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
			"ButtonGroup Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
}