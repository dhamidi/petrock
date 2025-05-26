package gallery

import (
	"net/http"
	"sort"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	html "maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
)

// ComponentInfo represents metadata about a UI component
type ComponentInfo struct {
	Name        string
	Description string
	Category    string
}

// GetAllComponents returns a list of all available UI components
func GetAllComponents() []ComponentInfo {
	return []ComponentInfo{
		{
			Name:        "container",
			Description: "Responsive container with different width variants",
			Category:    "Layout",
		},
		{
			Name:        "grid",
			Description: "Flexible CSS Grid container for complex layouts",
			Category:    "Layout",
		},
		{
			Name:        "section",
			Description: "Semantic section component with proper heading hierarchy",
			Category:    "Layout",
		},
		{
			Name:        "divider",
			Description: "Horizontal separator with different styles and spacing",
			Category:    "Layout",
		},
		{
			Name:        "card",
			Description: "Structured content container with header, body, and footer sections",
			Category:    "Content",
		},
		{
			Name:        "button",
			Description: "Interactive button with multiple variants, sizes, and states",
			Category:    "Interactive",
		},
		{
			Name:        "button-group",
			Description: "Container for grouping related buttons with consistent spacing",
			Category:    "Interactive",
		},
		{
		Name:        "form-inputs",
		Description: "Essential form input components including text inputs, textareas, and select dropdowns",
		Category:    "Form",
		},
	{
		Name:        "form-controls",
		Description: "Interactive form controls including checkboxes, radio buttons, and toggle switches",
		Category:    "Form",
	},
}
}

// BuildSidebar creates the component navigation sidebar content
func BuildSidebar() []g.Node {
	components := GetAllComponents()

	var sidebarContent []g.Node
	if len(components) == 0 {
		sidebarContent = append(sidebarContent,
			html.Div(
				html.Class("text-gray-500 italic"),
				g.Text("No components available yet"),
			),
		)
	} else {
		// Group components by category and create navigation
		categories := make(map[string][]ComponentInfo)
		for _, comp := range components {
			categories[comp.Category] = append(categories[comp.Category], comp)
		}

		// Sort categories alphabetically
		var categoryNames []string
		for category := range categories {
			categoryNames = append(categoryNames, category)
		}
		sort.Strings(categoryNames)

		// Iterate over sorted categories
		for _, category := range categoryNames {
			comps := categories[category]
			categorySection := []g.Node{
				html.H2(
					html.Class("text-sm font-medium text-gray-600 uppercase tracking-wide mb-3"),
					g.Text(category),
				),
			}

			var compLinks []g.Node
			for _, comp := range comps {
				compLinks = append(compLinks,
					html.Li(
						html.A(
							html.Href("/_/ui/"+comp.Name),
							html.Class("block px-3 py-2 text-blue-600 hover:bg-blue-50 rounded"),
							g.Text(comp.Name),
						),
					),
				)
			}

			categorySection = append(categorySection, html.Ul(
				html.Class("space-y-1 mb-6"),
				g.Group(compLinks),
			))

			sidebarContent = append(sidebarContent, html.Div(
				g.Group(categorySection),
			))
		}
	}

	return sidebarContent
}

// HandleGallery returns an HTTP handler for the main gallery page
func HandleGallery(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create sidebar content
		sidebarContent := BuildSidebar()

		// Create main content using existing Page component
		pageContent := core.Page("UI Component Gallery",
			html.Div(
				Classes{"flex": true, "min-h-screen": true, "-mx-4": true, "-mt-4": true},
				// Sidebar
				html.Nav(
					Classes{"w-64": true, "bg-white": true, "border-r": true, "border-gray-200": true, "p-6": true, "overflow-y-auto": true},
					html.H1(
						Classes{"text-lg": true, "font-semibold": true, "text-gray-900": true, "mb-6": true},
						g.Text("Components"),
					),
					g.Group(sidebarContent),
				),
				// Main content
				html.Main(
					Classes{"flex-1": true, "p-6": true, "overflow-y-auto": true},
					html.Div(
						Classes{"max-w-4xl": true},
						html.P(
							Classes{"text-lg": true, "text-gray-600": true, "mb-4": true},
							g.Text("Welcome to the UI component gallery. This is your central place to explore, test, and understand all available UI components in the design system."),
						),
						html.P(
							Classes{"text-gray-600": true, "mb-4": true},
							g.Text("Each component includes interactive examples, usage guidelines, and accessibility information to help you build consistent and accessible user interfaces."),
						),
						html.P(
							Classes{"text-gray-600": true},
							g.Text("Components will appear in the sidebar as they are implemented. The gallery will be populated as the design system grows."),
						),
					),
				),
			),
		)

		// Use existing Layout function
		layout := core.Layout("UI Component Gallery", pageContent)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := layout.Render(w)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
	}
}
