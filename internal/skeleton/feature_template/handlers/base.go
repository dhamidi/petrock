package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/petrock/example_module_path/core"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// FeatureServer holds dependencies required by the feature's HTTP handlers.
// This struct is initialized in register.go and passed to RegisterRoutes.
type FeatureServer struct {
	app     *core.App // The central App instance with all core dependencies
	querier *Querier  // Query execution logic
	state   *State    // Direct state access (use querier/executor preferably)
}

// NewFeatureServer creates and initializes the FeatureServer with its dependencies.
// Note: It now receives the central App instance.
func NewFeatureServer(
	app *core.App, // The central App instance
	querier *Querier,
	state *State,
) *FeatureServer {
	// Basic validation
	if app == nil || querier == nil || state == nil {
		panic("missing required dependencies for FeatureServer")
	}
	return &FeatureServer{
		app:     app,
		querier: querier,
		state:   state,
	}
}

// --- Helper Functions ---

// respondJSON is a utility to send JSON responses.
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// Log error, but can't change status code now
			slog.Error("Failed to encode JSON response", "error", err)
		}
	}
}

// parseIntParam is a helper to parse integer query parameters with a default value.
func ParseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return val
}

// --- View Helper Functions ---

// RenderPage is a helper function to render a complete HTML page with proper layout
func RenderPage(w http.ResponseWriter, pageTitle string, content g.Node) error {
	return RenderPageWithSuccess(w, pageTitle, content, "")
}

// RenderPageWithSuccess renders a complete HTML page with a success message
func RenderPageWithSuccess(w http.ResponseWriter, pageTitle string, content g.Node, successMsg string) error {
	// Set content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Create the page using a modern layout
	html := html.HTML(
		html.Lang("en"),
		html.Head(
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")),
			html.TitleEl(g.Text(pageTitle)),
			// Link to Tailwind CSS (modern version)
			html.Script(
				html.Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"),
				html.Async(),
				html.Defer(),
			),
			// Add a modern font
			html.Link(
				html.Rel("stylesheet"),
				html.Href("https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap"),
			),
		),
		html.Body(
			// Modern styling
			g.Attr("class", "bg-gradient-to-br from-slate-50 to-slate-100 min-h-screen font-sans antialiased text-slate-800"),

			// Header - full width
			html.Header(
				g.Attr("class", "bg-white shadow-sm border-b border-slate-200"),
				html.Div(
					g.Attr("class", "container mx-auto px-4 sm:px-6 lg:px-8 py-4"),
					html.Div(
						g.Attr("class", "flex justify-between items-center"),
						html.Div(
							g.Attr("class", "flex items-center"),
							html.A(
								g.Attr("href", "/"),
								g.Attr("class", "text-xl font-semibold text-indigo-600"),
								g.Text("Petrock App"),
							),
						),
						html.Nav(
							g.Attr("class", "flex space-x-4"),
							html.A(
								g.Attr("href", "/petrock_example_feature_name"),
								g.Attr("class", "text-sm font-medium text-slate-700 hover:text-indigo-600"),
								g.Text("Items"),
							),
							html.A(
								g.Attr("href", "/petrock_example_feature_name/new"),
								g.Attr("class", "text-sm font-medium text-slate-700 hover:text-indigo-600"),
								g.Text("New Item"),
							),
						),
					),
				),
			),

			// Main content - centered on larger screens
			html.Main(
				g.Attr("class", "container mx-auto px-4 sm:px-6 lg:px-8 py-8"),
				html.Div(
					g.Attr("class", "max-w-4xl mx-auto"),
					// Page title
					html.H1(
						g.Attr("class", "text-2xl font-bold text-slate-900 mb-6"),
						g.Text(pageTitle),
					),
					// Success message (if any)
					func() g.Node {
						if successMsg == "" {
							return nil
						}
						return html.Div(
							g.Attr("class", "mb-6 rounded-md bg-green-50 p-4 border border-green-200"),
							html.Div(
								g.Attr("class", "flex"),
								html.Div(
									g.Attr("class", "ml-3"),
									html.P(
										g.Attr("class", "text-sm font-medium text-green-800"),
										g.Text("u2713 "+successMsg),
									),
								),
							),
						)
					}(),
					// Page content
					html.Div(
						g.Attr("class", "bg-white shadow-sm rounded-lg border border-slate-200 p-6"),
						content,
					),
				),
			),

			// Footer - full width
			html.Footer(
				g.Attr("class", "bg-white border-t border-slate-200 mt-auto"),
				html.Div(
					g.Attr("class", "container mx-auto px-4 sm:px-6 lg:px-8 py-4"),
					html.Div(
						g.Attr("class", "text-center text-sm text-slate-500"),
						g.Text("u00a9 2025 Petrock App - Built with petrock"),
					),
				),
			),
		),
	)

	// Render the HTML
	return html.Render(w)
}
