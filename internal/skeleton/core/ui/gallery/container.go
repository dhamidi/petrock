package gallery

import (
	"net/http"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	html "maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
)

// HandleContainerDetail returns an HTTP handler for the container component demo page
func HandleContainerDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create demo content showing different container variants
		content := html.Div(
			Classes{"space-y-8": true},
			
			// Header section
			html.Div(
				Classes{"mb-8": true},
				html.H1(
					Classes{"text-3xl": true, "font-bold": true, "text-gray-900": true, "mb-4": true},
					g.Text("Container Component"),
				),
				html.P(
					Classes{"text-lg": true, "text-gray-600": true, "mb-4": true},
					g.Text("The Container component provides responsive width constraints and consistent horizontal margins. It's the foundation for page layouts."),
				),
				html.A(
					html.Href("/_/ui"),
					Classes{"inline-flex": true, "items-center": true, "text-blue-600": true, "hover:text-blue-800": true, "text-sm": true, "font-medium": true},
					g.Text("← Back to Gallery"),
				),
			),
			
			// Variants section
			html.Section(
				Classes{"space-y-6": true},
				html.H2(
					Classes{"text-2xl": true, "font-semibold": true, "text-gray-900": true, "mb-4": true},
					g.Text("Container Variants"),
				),
				
				// Default container
				html.Div(
					Classes{"space-y-4": true},
					html.H3(
						Classes{"text-lg": true, "font-medium": true, "text-gray-900": true},
						g.Text("Default Container (max-width: 896px)"),
					),
					html.Div(
						Classes{"border": true, "border-gray-200": true, "bg-gray-50": true, "p-1": true},
						ui.Container(ui.ContainerProps{Variant: "default"},
							html.Div(
								Classes{"bg-blue-100": true, "border": true, "border-blue-300": true, "p-4": true, "text-center": true},
								g.Text("Default container content - good for most pages"),
							),
						),
					),
				),
				
				// Narrow container
				html.Div(
					Classes{"space-y-4": true},
					html.H3(
						Classes{"text-lg": true, "font-medium": true, "text-gray-900": true},
						g.Text("Narrow Container (max-width: 672px)"),
					),
					html.Div(
						Classes{"border": true, "border-gray-200": true, "bg-gray-50": true, "p-1": true},
						ui.Container(ui.ContainerProps{Variant: "narrow"},
							html.Div(
								Classes{"bg-green-100": true, "border": true, "border-green-300": true, "p-4": true, "text-center": true},
								g.Text("Narrow container - perfect for reading content"),
							),
						),
					),
				),
				
				// Wide container
				html.Div(
					Classes{"space-y-4": true},
					html.H3(
						Classes{"text-lg": true, "font-medium": true, "text-gray-900": true},
						g.Text("Wide Container (max-width: 1280px)"),
					),
					html.Div(
						Classes{"border": true, "border-gray-200": true, "bg-gray-50": true, "p-1": true},
						ui.Container(ui.ContainerProps{Variant: "wide"},
							html.Div(
								Classes{"bg-purple-100": true, "border": true, "border-purple-300": true, "p-4": true, "text-center": true},
								g.Text("Wide container - great for dashboards and data-heavy pages"),
							),
						),
					),
				),
				
				// Full width container
				html.Div(
					Classes{"space-y-4": true},
					html.H3(
						Classes{"text-lg": true, "font-medium": true, "text-gray-900": true},
						g.Text("Full Width Container"),
					),
					html.Div(
						Classes{"border": true, "border-gray-200": true, "bg-gray-50": true, "p-1": true},
						ui.Container(ui.ContainerProps{Variant: "full"},
							html.Div(
								Classes{"bg-red-100": true, "border": true, "border-red-300": true, "p-4": true, "text-center": true},
								g.Text("Full width container - spans the entire viewport width"),
							),
						),
					),
				),
			),
			
			// Usage section
			html.Section(
				Classes{"space-y-4": true, "mt-12": true},
				html.H2(
					Classes{"text-2xl": true, "font-semibold": true, "text-gray-900": true, "mb-4": true},
					g.Text("Usage"),
				),
				html.Pre(
					Classes{"bg-gray-100": true, "p-4": true, "rounded": true, "text-sm": true, "overflow-x-auto": true, "font-mono": true},
					g.Text(`// Default container
ui.Container(ui.ContainerProps{}, 
    // your content here
)

// Narrow container for reading
ui.Container(ui.ContainerProps{Variant: "narrow"}, 
    // your content here
)

// Wide container for data-heavy layouts
ui.Container(ui.ContainerProps{Variant: "wide"}, 
    // your content here
)

// Full width container
ui.Container(ui.ContainerProps{Variant: "full"}, 
    // your content here
)

// Custom max-width
ui.Container(ui.ContainerProps{MaxWidth: "600px"}, 
    // your content here
)`),
				),
			),
			
			// Properties section
			html.Section(
				Classes{"space-y-4": true, "mt-8": true},
				html.H2(
					Classes{"text-2xl": true, "font-semibold": true, "text-gray-900": true, "mb-4": true},
					g.Text("Properties"),
				),
				html.Div(
					Classes{"overflow-x-auto": true},
					html.Table(
						Classes{"min-w-full": true, "border": true, "border-gray-200": true},
						html.THead(
							Classes{"bg-gray-50": true},
							html.Tr(
								html.Th(Classes{"px-4": true, "py-2": true, "border-b": true, "text-left": true, "font-medium": true}, g.Text("Property")),
								html.Th(Classes{"px-4": true, "py-2": true, "border-b": true, "text-left": true, "font-medium": true}, g.Text("Type")),
								html.Th(Classes{"px-4": true, "py-2": true, "border-b": true, "text-left": true, "font-medium": true}, g.Text("Default")),
								html.Th(Classes{"px-4": true, "py-2": true, "border-b": true, "text-left": true, "font-medium": true}, g.Text("Description")),
							),
						),
						html.TBody(
							html.Tr(
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true, "font-mono": true, "text-sm": true}, g.Text("Variant")),
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true}, g.Text("string")),
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true}, g.Text("\"default\"")),
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true}, g.Text("Container width variant: default, narrow, wide, full")),
							),
							html.Tr(
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true, "font-mono": true, "text-sm": true}, g.Text("MaxWidth")),
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true}, g.Text("string")),
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true}, g.Text("\"\"")),
								html.Td(Classes{"px-4": true, "py-2": true, "border-b": true}, g.Text("Custom max-width CSS value (overrides variant)")),
							),
						),
					),
				),
			),
		)

		// Create page using existing Page component
		pageContent := core.Page("Container Component",
			content,
		)

		// Use existing Layout function
		response := core.Layout(
			"Container Component - UI Gallery",
			pageContent,
		)

		w.Header().Set("Content-Type", "text/html")
		response.Render(w)
	}
}