package gallery

import (
	"net/http"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
)

// ComponentInfo represents metadata about a UI component
type ComponentInfo struct {
	Name        string
	Description string
	Category    string
}

// GetAllComponents returns a list of all available UI components
// Initially returns an empty list as components will be added in later steps
func GetAllComponents() []ComponentInfo {
	return []ComponentInfo{}
}

// galleryCSS contains the CSS styles for the gallery
const galleryCSS = `
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background-color: #f8f9fa;
    color: #212529;
}

.gallery-container {
    display: flex;
    min-height: 100vh;
}

.sidebar {
    width: 280px;
    background-color: #ffffff;
    border-right: 1px solid #e9ecef;
    padding: 24px;
    overflow-y: auto;
}

.sidebar h1 {
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 24px;
    color: #495057;
}

.content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
}

.category {
    margin-bottom: 32px;
}

.category h2 {
    font-size: 18px;
    font-weight: 500;
    margin-bottom: 12px;
    color: #6c757d;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.component-list {
    list-style: none;
}

.component-item {
    margin-bottom: 8px;
}

.component-link {
    display: block;
    padding: 12px 16px;
    text-decoration: none;
    color: #495057;
    border-radius: 6px;
    transition: background-color 0.2s ease;
}

.component-link:hover {
    background-color: #f8f9fa;
    color: #007bff;
}

.component-name {
    font-weight: 500;
    margin-bottom: 4px;
}

.component-description {
    font-size: 14px;
    color: #6c757d;
    line-height: 1.4;
}

.empty-state {
    text-align: center;
    color: #6c757d;
    font-style: italic;
    margin-top: 48px;
}

.welcome-content {
    max-width: 600px;
}

.welcome-content h1 {
    font-size: 32px;
    font-weight: 700;
    margin-bottom: 16px;
    color: #212529;
}

.welcome-content p {
    font-size: 16px;
    line-height: 1.6;
    color: #6c757d;
    margin-bottom: 24px;
}
`

// HandleGallery returns an HTTP handler for the main gallery page
func HandleGallery(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Get all available components
		components := GetAllComponents()

		// Group components by category
		categories := make(map[string][]ComponentInfo)
		for _, comp := range components {
			if comp.Category == "" {
				comp.Category = "General"
			}
			categories[comp.Category] = append(categories[comp.Category], comp)
		}

		// Render page using gomponents
		page := galleryPage(categories)
		err := page.Render(w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			return
		}
	}
}

// galleryPage renders the main gallery page using gomponents
func galleryPage(categories map[string][]ComponentInfo) g.Node {
	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1.0")),
			html.TitleEl(g.Text("UI Component Gallery")),
			html.Style(galleryCSS),
		),
		html.Body(
			html.Div(
				Classes{"gallery-container": true},
				gallerySidebar(categories),
				galleryMainContent(categories),
			),
		),
	)
}

// gallerySidebar renders the sidebar navigation
func gallerySidebar(categories map[string][]ComponentInfo) g.Node {
	var categoryNodes []g.Node

	if len(categories) > 0 {
		for categoryName, components := range categories {
			var componentItems []g.Node
			for _, comp := range components {
				componentItems = append(componentItems, html.Li(
					Classes{"component-item": true},
					html.A(
						html.Href("/_/ui/"+comp.Name),
						Classes{"component-link": true},
						html.Div(
							Classes{"component-name": true},
							g.Text(comp.Name),
						),
						html.Div(
							Classes{"component-description": true},
							g.Text(comp.Description),
						),
					),
				))
			}

			categoryNodes = append(categoryNodes, html.Div(
				Classes{"category": true},
				html.H2(g.Text(categoryName)),
				html.Ul(
					Classes{"component-list": true},
					g.Group(componentItems),
				),
			))
		}
	} else {
		categoryNodes = append(categoryNodes, html.Div(
			Classes{"empty-state": true},
			g.Text("No components available yet"),
		))
	}

	return html.Nav(
		Classes{"sidebar": true},
		html.H1(g.Text("Components")),
		g.Group(categoryNodes),
	)
}

// galleryMainContent renders the main content area
func galleryMainContent(categories map[string][]ComponentInfo) g.Node {
	var additionalMessage g.Node
	if len(categories) == 0 {
		additionalMessage = html.P(
			g.Text("Components will appear in the sidebar as they are implemented. "),
			g.Text("The gallery will be populated as the design system grows."),
		)
	}

	return html.Main(
		Classes{"content": true},
		html.Div(
			Classes{"welcome-content": true},
			html.H1(g.Text("UI Component Gallery")),
			html.P(
				g.Text("Welcome to the UI component gallery. This is your central place to explore, "),
				g.Text("test, and understand all available UI components in the design system."),
			),
			html.P(
				g.Text("Each component includes interactive examples, usage guidelines, and accessibility "),
				g.Text("information to help you build consistent and accessible user interfaces."),
			),
			g.If(additionalMessage != nil, additionalMessage),
		),
	)
}