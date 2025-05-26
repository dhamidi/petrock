package gallery

import (
	"net/http"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
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
		// Route to specific component handlers
		switch componentName {
		case "container":
			HandleContainerDetail(app)(w, r)
			return
		case "grid":
			HandleGridDetail(app)(w, r)
			return
		case "section":
			HandleSectionDetail(app)(w, r)
			return
		case "divider":
			HandleDividerDetail(app)(w, r)
			return
		case "card":
			HandleCardDetail(app)(w, r)
			return
		case "button":
			HandleButtonDetail(app)(w, r)
			return
		case "button-group":
			HandleButtonGroupDetail(app)(w, r)
			return
		}
		
		// Component found - show details
			pageContent = core.Page("Component: "+foundComponent.Name,
				html.Div(
					Classes{"flex": true, "min-h-screen": true, "-mx-4": true, "-mt-4": true},
					// Sidebar with back link
					html.Nav(
						Classes{"w-64": true, "bg-white": true, "border-r": true, "border-gray-200": true, "p-6": true},
						html.H1(
							Classes{"text-lg": true, "font-semibold": true, "text-gray-900": true, "mb-6": true},
							html.A(
								html.Href("/_/ui"),
								Classes{"text-blue-600": true, "hover:text-blue-800": true, "no-underline": true},
								g.Text("← Gallery"),
							),
						),
					),
					// Main content
					html.Main(
						Classes{"flex-1": true, "p-6": true},
						html.Div(
							Classes{"max-w-4xl": true},
							html.P(
								Classes{"text-lg": true, "text-gray-600": true, "mb-4": true},
								g.Text(foundComponent.Description),
							),
							html.P(
								Classes{"text-gray-600": true, "mb-4": true},
								g.Text("Interactive examples and documentation will appear here when the component is implemented."),
							),
						),
					),
				),
			)
		} else {
			// Component not found
			pageContent = core.Page("Component Not Found",
				html.Div(
					Classes{"flex": true, "min-h-screen": true, "-mx-4": true, "-mt-4": true},
					// Sidebar with back link
					html.Nav(
						Classes{"w-64": true, "bg-white": true, "border-r": true, "border-gray-200": true, "p-6": true},
						html.H1(
							Classes{"text-lg": true, "font-semibold": true, "text-gray-900": true, "mb-6": true},
							html.A(
								html.Href("/_/ui"),
								Classes{"text-blue-600": true, "hover:text-blue-800": true, "no-underline": true},
								g.Text("← Gallery"),
							),
						),
					),
					// Main content
					html.Main(
						Classes{"flex-1": true, "p-6": true},
						html.Div(
							Classes{"max-w-4xl": true},
							html.P(
								Classes{"text-lg": true, "text-red-600": true, "mb-4": true},
								g.Textf("The component \"%s\" does not exist in the gallery.", componentName),
							),
							html.P(
								Classes{"text-gray-600": true, "mb-4": true},
								html.A(
									html.Href("/_/ui"),
									Classes{"text-blue-600": true, "hover:text-blue-800": true},
									g.Text("← Back to Gallery"),
								),
							),
							html.P(
								Classes{"text-gray-500": true, "italic": true},
								g.Text("No components are available yet."),
							),
						),
					),
				),
			)
		}

		// Use existing Layout function
		layout := core.Layout("UI Component Gallery", pageContent)
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := layout.Render(w)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
	}
}