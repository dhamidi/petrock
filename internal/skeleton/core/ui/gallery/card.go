package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleCardDetail renders the Card component demo page
func HandleCardDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different card layouts
		demoContent := html.Div(
			ui.CSSClass("space-y-8"),
			
			// Header section
			html.Div(
				ui.CSSClass("mb-8"),
				html.H1(
					ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
					g.Text("Card Component"),
				),
				html.P(
					ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
					g.Text("The Card component provides a structured content container with header, body, and footer sections."),
				),
			),
			
			// Basic Card Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Card Variants"),
				),
				
				// Default Card
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Default Card"),
					),
					html.Div(
						ui.CSSClass("max-w-md"),
						ui.Card(ui.CardProps{Variant: "default"},
							ui.CardHeader(
								html.H4(ui.CSSClass("text-lg", "font-semibold"), g.Text("Card Title")),
								html.P(ui.CSSClass("text-sm", "text-gray-600"), g.Text("Card subtitle or description")),
							),
							ui.CardBody(
								html.P(ui.CSSClass("text-gray-700"), g.Text("This is the main content area of the card. It can contain any type of content including text, images, or other components.")),
							),
							ui.CardFooter(
								html.Div(
									ui.CSSClass("flex", "justify-end", "space-x-2"),
									html.Button(
										ui.CSSClass("px-4", "py-2", "text-sm", "text-gray-600", "hover:text-gray-800"),
										g.Text("Cancel"),
									),
									html.Button(
										ui.CSSClass("px-4", "py-2", "bg-blue-500", "text-white", "rounded", "text-sm", "hover:bg-blue-600"),
										g.Text("Action"),
									),
								),
							),
						),
					),
				),
				
				// Outlined Card
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Outlined Card"),
					),
					html.Div(
						ui.CSSClass("max-w-md"),
						ui.Card(ui.CardProps{Variant: "outlined"},
							ui.CardHeader(
								html.H4(ui.CSSClass("text-lg", "font-semibold"), g.Text("Outlined Card")),
							),
							ui.CardBody(
								html.P(ui.CSSClass("text-gray-700"), g.Text("This card uses the outlined variant with a simple border and no shadow.")),
							),
						),
					),
				),
				
				// Elevated Card
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Elevated Card"),
					),
					html.Div(
						ui.CSSClass("max-w-md"),
						ui.Card(ui.CardProps{Variant: "elevated"},
							ui.CardHeader(
								html.H4(ui.CSSClass("text-lg", "font-semibold"), g.Text("Elevated Card")),
							),
							ui.CardBody(
								html.P(ui.CSSClass("text-gray-700"), g.Text("This card uses the elevated variant with a larger shadow for emphasis.")),
							),
						),
					),
				),
			),
			
			// Padding Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Padding Variants"),
				),
				html.Div(
					ui.CSSClass("grid", "grid-cols-1", "md:grid-cols-3", "gap-4"),
					
					// Small padding
					html.Div(
						html.H4(ui.CSSClass("text-sm", "font-medium", "text-gray-900", "mb-2"), g.Text("Small Padding")),
						ui.Card(ui.CardProps{Variant: "outlined", Padding: "small"},
							g.Text("Small padding card"),
						),
					),
					
					// Medium padding (default)
					html.Div(
						html.H4(ui.CSSClass("text-sm", "font-medium", "text-gray-900", "mb-2"), g.Text("Medium Padding")),
						ui.Card(ui.CardProps{Variant: "outlined", Padding: "medium"},
							g.Text("Medium padding card"),
						),
					),
					
					// Large padding
					html.Div(
						html.H4(ui.CSSClass("text-sm", "font-medium", "text-gray-900", "mb-2"), g.Text("Large Padding")),
						ui.Card(ui.CardProps{Variant: "outlined", Padding: "large"},
							g.Text("Large padding card"),
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
					g.Text(`// Basic card with sections
ui.Card(ui.CardProps{Variant: "default"},
    ui.CardHeader(
        html.H4(g.Text("Title")),
        html.P(g.Text("Subtitle")),
    ),
    ui.CardBody(
        html.P(g.Text("Main content")),
    ),
    ui.CardFooter(
        html.Button(g.Text("Action")),
    ),
)

// Simple card with custom padding
ui.Card(ui.CardProps{Variant: "outlined", Padding: "large"},
    g.Text("Card content"),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Variant")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"default\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Card style variant: default, outlined, elevated")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Padding")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"medium\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Padding size: none, small, medium, large")),
							),
						),
					),
				),
			),
		)

		// Create page content with proper sidebar navigation
		pageContent := core.Page("Card Component",
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
			"Card Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
	}
}