package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleSectionDetail renders the Section component demo page
func HandleSectionDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different section layouts
		demoContent := html.Div(
			ui.CSSClass("space-y-8"),
			
			// Header section
			html.Div(
				ui.CSSClass("mb-8"),
				html.H1(
					ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
					g.Text("Section Component"),
				),
				html.P(
					ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
					g.Text("The Section component provides semantic sectioning with proper heading hierarchy and accessibility attributes."),
				),
			),
			
			// Heading Level Examples
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Heading Levels"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Sections support semantic heading levels 1-6 with appropriate typography:"),
				),
				
				// Level 1 example
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Section(ui.SectionProps{Heading: "Level 1 Section", Level: 1},
						html.P(ui.CSSClass("text-gray-700"), g.Text("This is a level 1 section with the largest heading.")),
					),
				),
				
				// Level 2 example
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Section(ui.SectionProps{Heading: "Level 2 Section", Level: 2},
						html.P(ui.CSSClass("text-gray-700"), g.Text("This is a level 2 section with a large heading.")),
					),
				),
				
				// Level 3 example
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Section(ui.SectionProps{Heading: "Level 3 Section", Level: 3},
						html.P(ui.CSSClass("text-gray-700"), g.Text("This is a level 3 section with a medium heading.")),
					),
				),
				
				// Level 4 example
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Section(ui.SectionProps{Heading: "Level 4 Section", Level: 4},
						html.P(ui.CSSClass("text-gray-700"), g.Text("This is a level 4 section with a small heading.")),
					),
				),
			),
			
			// Nested Content Example
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Nested Content"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Sections can contain multiple children with proper spacing:"),
				),
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Section(ui.SectionProps{Heading: "Article Section", Level: 2},
						html.P(ui.CSSClass("text-gray-700"), g.Text("This section contains multiple paragraphs and elements.")),
						html.P(ui.CSSClass("text-gray-700"), g.Text("Each child element is properly spaced using the space-y-4 utility.")),
						html.Ul(
							ui.CSSClass("list-disc", "pl-6", "text-gray-700"),
							html.Li(g.Text("First list item")),
							html.Li(g.Text("Second list item")),
							html.Li(g.Text("Third list item")),
						),
					),
				),
			),
			
			// Section without heading
			html.Section(
				ui.CSSClass("space-y-6"),
				html.H2(
					ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
					g.Text("Section without Heading"),
				),
				html.P(
					ui.CSSClass("text-gray-600", "mb-4"),
					g.Text("Sections can be used without headings for semantic grouping:"),
				),
				html.Div(
					ui.CSSClass("border", "p-4", "rounded", "mb-4"),
					ui.Section(ui.SectionProps{},
						html.P(ui.CSSClass("text-gray-700"), g.Text("This section has no heading but still provides semantic structure.")),
						html.P(ui.CSSClass("text-gray-700"), g.Text("It's useful for grouping related content without visual hierarchy.")),
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
					g.Text(`// Basic section with heading
ui.Section(ui.SectionProps{Heading: "Section Title", Level: 2},
    html.P(g.Text("Section content")),
    html.P(g.Text("More content")),
)

// Section without heading (semantic grouping)
ui.Section(ui.SectionProps{},
    html.P(g.Text("Content without heading")),
)

// Section with different heading levels
ui.Section(ui.SectionProps{Heading: "Main Section", Level: 1},
    ui.Section(ui.SectionProps{Heading: "Subsection", Level: 2},
        html.P(g.Text("Nested content")),
    ),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Heading")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Optional heading text for the section")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Level")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("int")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("2")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Heading level (1-6) for proper semantic hierarchy")),
							),
						),
					),
				),
			),
		)

		// Create page content with proper sidebar navigation
		pageContent := core.Page("Section Component",
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
			"Section Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
	}
}