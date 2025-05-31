package gallery

import (
	"net/http"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
)

// HandleComponentDetail returns an HTTP handler for individual component detail pages
func HandleComponentDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		componentName := r.PathValue("component")
		if componentName == "" {
			http.Error(w, "Component name is required", http.StatusBadRequest)
			return
		}

		components := GetAllComponents()
		var foundComponent *ComponentInfo
		
		// Find the requested component
		for _, comp := range components {
			if comp.Name == componentName {
				foundComponent = &comp
				break
			}
		}

		var pageContent g.Node
	
	if foundComponent != nil {
		// Use the handler from ComponentInfo
		if foundComponent.Handler != nil {
			foundComponent.Handler(w, r)
			return
		}
		
		// Component found - show details
		pageContent = ui.Page("Component: "+foundComponent.Name,
		html.Div(
		ui.CSSClass("flex", "min-h-screen", "-mx-4", "-mt-4"),
		// Sidebar with back link
		html.Nav(
		ui.CSSClass("w-64", "bg-white", "border-r", "border-gray-200", "p-6"),
		html.H1(
		ui.CSSClass("text-lg", "font-semibold", "text-gray-900", "mb-6"),
		html.A(
		html.Href("/_/ui"),
		ui.CSSClass("text-blue-600", "hover:text-blue-800", "no-underline"),
		g.Text("← Gallery"),
		),
		),
		),
		// Main content
		html.Main(
		ui.CSSClass("flex-1", "p-6"),
		html.Div(
		ui.CSSClass("max-w-4xl"),
		html.P(
		ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
		g.Text(foundComponent.Description),
		),
		html.P(
		ui.CSSClass("text-gray-600", "mb-4"),
		g.Text("Interactive examples and documentation will appear here when the component is implemented."),
		),
		),
		),
		),
		)
		} else {
			// Component not found
			pageContent = ui.Page("Component Not Found",
			html.Div(
			ui.CSSClass("flex", "min-h-screen", "-mx-4", "-mt-4"),
			// Sidebar with back link
			html.Nav(
			ui.CSSClass("w-64", "bg-white", "border-r", "border-gray-200", "p-6"),
			html.H1(
			ui.CSSClass("text-lg", "font-semibold", "text-gray-900", "mb-6"),
			html.A(
			html.Href("/_/ui"),
			ui.CSSClass("text-blue-600", "hover:text-blue-800", "no-underline"),
			g.Text("← Gallery"),
			),
			),
			),
			// Main content
			html.Main(
			ui.CSSClass("flex-1", "p-6"),
			html.Div(
			ui.CSSClass("max-w-4xl"),
			html.P(
			ui.CSSClass("text-lg", "text-red-600", "mb-4"),
			g.Textf("The component \"%s\" does not exist in the gallery.", componentName),
			),
			html.P(
			ui.CSSClass("text-gray-600", "mb-4"),
			html.A(
			html.Href("/_/ui"),
			ui.CSSClass("text-blue-600", "hover:text-blue-800"),
			g.Text("← Back to Gallery"),
			),
			),
			html.P(
			ui.CSSClass("text-gray-500", "italic"),
			g.Text("No components are available yet."),
			),
			),
			),
			),
			)
		}

		// Use existing Layout function
		layout := ui.Layout("UI Component Gallery", pageContent)
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := layout.Render(w)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
	}
}