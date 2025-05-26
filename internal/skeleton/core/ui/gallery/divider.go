package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleDividerDetail renders the Divider component demo page
func HandleDividerDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different divider variants
		demoContent := html.Div(
			ui.CSSClass("space-y-8"),
			
			// Header section
			html.Div(
				ui.CSSClass("mb-8"),
				html.H1(
					ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
					g.Text("Divider Component"),
				),
				html.P(
					ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
					g.Text("The Divider component creates horizontal separators with different styles and spacing options."),
				),
			),
			
			// Variant Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Divider Variants"),
				),
				
				// Solid divider
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Solid Divider (Default)"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content above the divider")),
						ui.Divider(ui.DividerProps{Variant: "solid"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content below the divider")),
					),
				),
				
				// Dashed divider
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Dashed Divider"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content above the divider")),
						ui.Divider(ui.DividerProps{Variant: "dashed"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content below the divider")),
					),
				),
				
				// Dotted divider
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Dotted Divider"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content above the divider")),
						ui.Divider(ui.DividerProps{Variant: "dotted"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content below the divider")),
					),
				),
			),
			
			// Margin Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Spacing Variants"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Dividers support different spacing options to control the vertical space around them:"),
				),
				
				// Small margin
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Small Margin"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with small spacing")),
						ui.Divider(ui.DividerProps{Margin: "small"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with small spacing")),
					),
				),
				
				// Medium margin (default)
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Medium Margin (Default)"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with medium spacing")),
						ui.Divider(ui.DividerProps{Margin: "medium"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with medium spacing")),
					),
				),
				
				// Large margin
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Large Margin"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with large spacing")),
						ui.Divider(ui.DividerProps{Margin: "large"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with large spacing")),
					),
				),
				
				// No margin
				html.Div(
					ui.CSSClass("mb-6"),
					html.H3(
						ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
						g.Text("No Margin"),
					),
					html.Div(
						ui.CSSClass("border", "p-4", "rounded", "bg-gray-50"),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with no spacing")),
						ui.Divider(ui.DividerProps{Margin: "none"}),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Content with no spacing")),
					),
				),
			),
			
			// Usage in Content
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Usage in Content"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Dividers are commonly used to separate content sections:"),
				),
				html.Div(
					ui.CSSClass("border", "p-6", "rounded", "bg-white"),
					html.H3(ui.CSSClass("text-lg", "font-semibold", "text-gray-900"), g.Text("Article Title")),
					html.P(ui.CSSClass("text-gray-700"), g.Text("This is the first paragraph of the article content.")),
					ui.Divider(ui.DividerProps{}),
					html.H4(ui.CSSClass("text-base", "font-semibold", "text-gray-900"), g.Text("Subsection")),
					html.P(ui.CSSClass("text-gray-700"), g.Text("This content is separated by a divider from the previous section.")),
					ui.Divider(ui.DividerProps{Variant: "dashed", Margin: "large"}),
					html.P(ui.CSSClass("text-gray-700", "italic"), g.Text("Final thoughts or conclusion section.")),
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
					g.Text(`// Basic divider with default settings
ui.Divider(ui.DividerProps{})

// Dashed divider with large margins
ui.Divider(ui.DividerProps{
    Variant: "dashed",
    Margin: "large",
})

// Dotted divider with no margins
ui.Divider(ui.DividerProps{
    Variant: "dotted",
    Margin: "none",
})

// In content layout
html.Div(
    html.P(g.Text("First section")),
    ui.Divider(ui.DividerProps{}),
    html.P(g.Text("Second section")),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"solid\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Divider style: solid, dashed, dotted")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Margin")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"medium\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Vertical spacing: none, small, medium, large")),
							),
						),
					),
				),
			),
		)

		// Create page content with proper sidebar navigation
		pageContent := core.Page("Divider Component",
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
			"Divider Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
	}
}