package gallery

import (
	"net/http"
	"strings"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
)

// componentCSS contains the CSS styles for the component detail page
const componentCSS = `
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

.component-container {
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
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 16px;
}

.sidebar a {
    color: #007bff;
    text-decoration: none;
}

.sidebar a:hover {
    text-decoration: underline;
}

.content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
}

.component-header {
    margin-bottom: 32px;
    padding-bottom: 16px;
    border-bottom: 1px solid #e9ecef;
}

.component-title {
    font-size: 32px;
    font-weight: 700;
    margin-bottom: 8px;
    color: #212529;
}

.component-description {
    font-size: 16px;
    color: #6c757d;
    line-height: 1.6;
}

.component-category {
    display: inline-block;
    background-color: #e7f3ff;
    color: #0056b3;
    padding: 4px 12px;
    border-radius: 16px;
    font-size: 12px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 16px;
}

.error-state {
    text-align: center;
    color: #dc3545;
    margin-top: 48px;
}

.error-state h1 {
    font-size: 24px;
    margin-bottom: 16px;
}

.error-state p {
    font-size: 16px;
    margin-bottom: 24px;
}

.back-link {
    display: inline-block;
    padding: 8px 16px;
    background-color: #007bff;
    color: white;
    text-decoration: none;
    border-radius: 6px;
    font-weight: 500;
}

.back-link:hover {
    background-color: #0056b3;
    text-decoration: none;
}

.placeholder-content {
    background-color: #ffffff;
    border: 2px dashed #dee2e6;
    border-radius: 8px;
    padding: 48px;
    text-align: center;
    color: #6c757d;
    margin-top: 24px;
}

.placeholder-content h3 {
    margin-bottom: 12px;
    color: #495057;
}

.component-list {
    list-style: none;
    margin-top: 16px;
}

.component-item {
    margin-bottom: 4px;
}

.component-link {
    display: block;
    padding: 8px 12px;
    text-decoration: none;
    color: #495057;
    border-radius: 4px;
    transition: background-color 0.2s ease;
}

.component-link:hover {
    background-color: #f8f9fa;
    color: #007bff;
}

.component-link.active {
    background-color: #e7f3ff;
    color: #0056b3;
    font-weight: 500;
}
`

// HandleComponentDetail returns an HTTP handler for individual component detail pages
func HandleComponentDetail(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract component name from URL path
		// Expected pattern: /_/ui/{component}
		path := strings.TrimPrefix(r.URL.Path, "/_/ui/")
		componentName := strings.Trim(path, "/")

		// Validate component name
		if componentName == "" {
			http.Error(w, "Component name is required", http.StatusBadRequest)
			return
		}

		// Check if component exists
		components := GetAllComponents()
		var component *ComponentInfo
		for _, comp := range components {
			if strings.EqualFold(comp.Name, componentName) {
				component = &comp
				break
			}
		}

		// Set content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Render page using gomponents
		page := componentDetailPage(component, componentName, components)
		err := page.Render(w)
		if err != nil {
			http.Error(w, "Rendering error", http.StatusInternalServerError)
			return
		}
	}
}

// componentDetailPage renders the component detail page using gomponents
func componentDetailPage(component *ComponentInfo, componentName string, allComponents []ComponentInfo) g.Node {
	title := "UI Component Gallery"
	if component != nil {
		title = component.Name + " - " + title
	}

	return html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1.0")),
			html.TitleEl(g.Text(title)),
			html.Style(componentCSS),
		),
		html.Body(
			html.Div(
				Classes{"component-container": true},
				componentSidebar(componentName, allComponents),
				componentMainContent(component, componentName, allComponents),
			),
		),
	)
}

// componentSidebar renders the sidebar navigation
func componentSidebar(currentComponentName string, allComponents []ComponentInfo) g.Node {
	var componentItems []g.Node
	
	if len(allComponents) > 0 {
		componentItems = append(componentItems, html.H2(
			html.Style("font-size: 14px; font-weight: 500; color: #6c757d; margin: 16px 0 8px 0; text-transform: uppercase; letter-spacing: 0.5px;"),
			g.Text("Components"),
		))

		var items []g.Node
		for _, comp := range allComponents {
			linkClasses := Classes{"component-link": true}
			if strings.EqualFold(comp.Name, currentComponentName) {
				linkClasses["active"] = true
			}

			items = append(items, html.Li(
				Classes{"component-item": true},
				html.A(
					html.Href("/_/ui/"+comp.Name),
					linkClasses,
					g.Text(comp.Name),
				),
			))
		}

		componentItems = append(componentItems, html.Ul(
			Classes{"component-list": true},
			g.Group(items),
		))
	}

	return html.Nav(
		Classes{"sidebar": true},
		html.H1(
			html.A(
				html.Href("/_/ui"),
				g.Text("← Gallery"),
			),
		),
		g.Group(componentItems),
	)
}

// componentMainContent renders the main content area
func componentMainContent(component *ComponentInfo, componentName string, allComponents []ComponentInfo) g.Node {
	if component != nil {
		// Component exists, show its details
		var categoryBadge g.Node
		if component.Category != "" {
			categoryBadge = html.Div(
				Classes{"component-category": true},
				g.Text(component.Category),
			)
		}

		return html.Main(
			Classes{"content": true},
			html.Div(
				Classes{"component-header": true},
				categoryBadge,
				html.H1(
					Classes{"component-title": true},
					g.Text(component.Name),
				),
				html.P(
					Classes{"component-description": true},
					g.Text(component.Description),
				),
			),
			html.Div(
				Classes{"placeholder-content": true},
				html.H3(g.Text("Component Implementation")),
				html.P(
					g.Text("This component will be implemented in a future step. "),
					g.Text("Interactive examples and documentation will appear here."),
				),
			),
		)
	}

	// Component not found, show error state
	var availableComponents g.Node
	if len(allComponents) > 0 {
		var items []g.Node
		for _, comp := range allComponents {
			items = append(items, html.Li(
				html.Style("margin-bottom: 8px;"),
				html.A(
					html.Href("/_/ui/"+comp.Name),
					html.Style("color: #007bff; text-decoration: none;"),
					g.Text(comp.Name),
				),
			))
		}

		availableComponents = html.Div(
			html.Style("margin-top: 32px;"),
			html.H3(
				html.Style("color: #495057; margin-bottom: 16px;"),
				g.Text("Available Components:"),
			),
			html.Ul(
				html.Style("list-style: none; display: inline-block; text-align: left;"),
				g.Group(items),
			),
		)
	} else {
		availableComponents = html.Div(
			html.Style("margin-top: 32px;"),
			html.P(
				html.Style("color: #6c757d; font-style: italic;"),
				g.Text("No components are available yet."),
			),
		)
	}

	return html.Main(
		Classes{"content": true},
		html.Div(
			Classes{"error-state": true},
			html.H1(g.Text("Component Not Found")),
			html.P(g.Text("The component \""+componentName+"\" does not exist in the gallery.")),
			html.A(
				html.Href("/_/ui"),
				Classes{"back-link": true},
				g.Text("← Back to Gallery"),
			),
			availableComponents,
		),
	)
}