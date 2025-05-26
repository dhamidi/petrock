package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleGridDetail renders the Grid component demo page
func HandleGridDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different grid layouts
		demoContent := html.Div(
			ui.CSSClass("space-y-8"),
			
			// Header section
			html.Div(
				ui.CSSClass("mb-8"),
				html.H1(
					ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
					g.Text("Grid Component"),
				),
				html.P(
					ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
					g.Text("The Grid component provides flexible CSS Grid layouts with customizable columns, gaps, and grid template areas."),
				),
				html.A(
					html.Href("/_/ui"),
					ui.CSSClass("inline-flex", "items-center", "text-blue-600", "hover:text-blue-800", "text-sm", "font-medium"),
					g.Text("← Back to Gallery"),
				),
			),
			
			// Basic Grid Example
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Basic Two-Column Grid"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("A simple two-column grid with equal widths:"),
				),
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Grid(ui.GridProps{Columns: "1fr 1fr", Gap: "1rem"},
						html.Div(ui.CSSClass("bg-blue-100", "p-4", "text-center"), g.Text("Column 1")),
						html.Div(ui.CSSClass("bg-green-100", "p-4", "text-center"), g.Text("Column 2")),
					),
				),
				html.Pre(
					ui.CSSClass("bg-gray-100", "p-4", "rounded", "text-sm"),
					g.Text(`ui.Grid(ui.GridProps{Columns: "1fr 1fr", Gap: "1rem"},
    html.Div(g.Text("Column 1")),
    html.Div(g.Text("Column 2")),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Columns")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"1fr\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("CSS grid-template-columns value")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Gap")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"1rem\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("CSS gap value for spacing")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Areas")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("CSS grid-template-areas for named regions")),
							),
						),
					),
				),
			),
		)

		// Create page content with proper sidebar navigation
		pageContent := core.Page("Grid Component",
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
			"Grid Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
	}
}