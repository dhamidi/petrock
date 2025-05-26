package gallery

import (
	"net/http"
	"sort"

	g "maragu.dev/gomponents"
	html "maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
)

// ComponentInfo represents metadata about a UI component
type ComponentInfo struct {
	Name        string
	Description string
	Category    string
	Handler     http.HandlerFunc
}

// GetAllComponents returns a list of all available UI components
func GetAllComponents() []ComponentInfo {
	return []ComponentInfo{
		{
			Name:        "container",
			Description: "Responsive container with different width variants",
			Category:    "Layout",
			Handler:     HandleContainerDetail,
		},
		{
			Name:        "grid",
			Description: "Flexible CSS Grid container for complex layouts",
			Category:    "Layout", 
			Handler:     HandleGridDetail,
		},
		{
			Name:        "section",
			Description: "Semantic section component with proper heading hierarchy",
			Category:    "Layout",
			Handler:     HandleSectionDetail,
		},
		{
			Name:        "divider",
			Description: "Horizontal separator with different styles and spacing",
			Category:    "Layout",
			Handler:     HandleDividerDetail,
		},
		{
			Name:        "card",
			Description: "Structured content container with header, body, and footer sections",
			Category:    "Content",
			Handler:     HandleCardDetail,
		},
		{
			Name:        "button",
			Description: "Interactive button with multiple variants, sizes, and states",
			Category:    "Interactive",
			Handler:     HandleButtonDetail,
		},
		{
			Name:        "button-group",
			Description: "Container for grouping related buttons with consistent spacing",
			Category:    "Interactive",
			Handler:     HandleButtonGroupDetail,
		},
		{
			Name:        "form-inputs",
			Description: "Essential form input components including text inputs, textareas, and select dropdowns",
			Category:    "Form",
			Handler:     HandleFormInputsDetail,
		},
		{
			Name:        "form-controls",
			Description: "Interactive form controls including checkboxes, radio buttons, and toggle switches",
			Category:    "Form",
			Handler:     HandleFormControlsDetail,
		},
		{
			Name:        "form-layout",
			Description: "Form layout components including FormGroup and FieldSet for organizing form elements",
			Category:    "Form",
			Handler:     HandleFormLayoutDetail,
		},
		{
			Name:        "navigation",
			Description: "Navigation components including navigation bars, sidebars, tabs, breadcrumbs, and pagination",
			Category:    "Navigation",
			Handler:     HandleNavigationDetail,
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
				ui.CSSClass("flex", "min-h-screen", "-mx-4", "-mt-4"),
				// Sidebar
				html.Nav(
					ui.CSSClass("w-64", "bg-white", "border-r", "border-gray-200", "p-6", "overflow-y-auto"),
					html.H1(
						ui.CSSClass("text-lg", "font-semibold", "text-gray-900", "mb-6"),
						g.Text("Components"),
					),
					g.Group(sidebarContent),
				),
				// Main content
				html.Main(
					ui.CSSClass("flex-1", "p-6", "overflow-y-auto"),
					html.Div(
						ui.CSSClass("max-w-4xl"),
						html.P(
							ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
							g.Text("Welcome to the UI component gallery. This is your central place to explore, test, and understand all available UI components in the design system."),
						),
						html.P(
							ui.CSSClass("text-gray-600", "mb-4"),
							g.Text("Each component includes interactive examples, usage guidelines, and accessibility information to help you build consistent and accessible user interfaces."),
						),
						html.P(
							ui.CSSClass("text-gray-600"),
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
