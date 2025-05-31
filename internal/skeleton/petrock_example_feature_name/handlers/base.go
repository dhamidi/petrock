package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
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

// RenderPage is a helper function to render a complete HTML page using ui.Layout and ui.Page
func RenderPage(w http.ResponseWriter, pageTitle string, content g.Node) error {
	return RenderPageWithSuccess(w, pageTitle, content, "")
}

// RenderPageWithSuccess renders a complete HTML page with a success message using ui components
func RenderPageWithSuccess(w http.ResponseWriter, pageTitle string, content g.Node, successMsg string) error {
	// Set content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Create success message content if provided
	var pageContent g.Node
	if successMsg != "" {
		successAlert := html.Div(
			ui.CSSClass("mb-6", "rounded-md", "bg-green-50", "p-4", "border", "border-green-200"),
			html.Div(
				ui.CSSClass("flex"),
				html.Div(
					ui.CSSClass("ml-3"),
					html.P(
						ui.CSSClass("text-sm", "font-medium", "text-green-800"),
						html.Span(
							ui.CSSClass("inline-block", "mr-1"),
							g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`),
						),
						g.Text(successMsg),
					),
				),
			),
		)
		pageContent = html.Div(
			successAlert,
			content,
		)
	} else {
		pageContent = content
	}

	// Use ui.Layout and ui.Page for consistent styling
	page := ui.Layout(
		pageTitle,
		ui.Page(pageTitle, pageContent),
	)

	// Render the HTML
	return page.Render(w)
}
