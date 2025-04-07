package core

import (
	"net/http" // Added for http.ResponseWriter, http.Request

	g "maragu.dev/gomponents"            // Use canonical import path
	. "maragu.dev/gomponents/components" // Dot-import for helpers like Classes
	"maragu.dev/gomponents/html"
)

// IndexPage renders the main content for the home page, displaying registered commands and queries.
// This component is intended to be passed to the Layout function.
func IndexPage(commandNames, queryNames []string) g.Node {
	commandListItems := make([]g.Node, len(commandNames))
	for i, name := range commandNames {
		commandListItems[i] = html.Li(html.Code(g.Text(name)))
	}

	queryListItems := make([]g.Node, len(queryNames))
	for i, name := range queryNames {
		queryListItems[i] = html.Li(html.Code(g.Text(name)))
	}

	return Page("Welcome!", // Use the Page component from core/view.go
		html.P(
			Classes{"text-lg": true, "mb-4": true},
			g.Text("Welcome to your Petrock-generated application!"),
		),
		html.H2(Classes{"text-xl": true, "font-semibold": true, "mt-6": true, "mb-2": true}, g.Text("Available Commands")),
		html.Ul(
			Classes{"list-disc": true, "list-inside": true, "space-y-1": true, "mb-4": true},
			g.Group(commandListItems),
		),
		html.P(
			g.Text("Execute commands via: "),
			html.Code(Classes{"bg-gray-200": true, "p-1": true, "rounded": true}, g.Text("POST /commands")),
			g.Text(" with JSON payload: "),
			html.Code(Classes{"bg-gray-200": true, "p-1": true, "rounded": true}, g.Text(`{"type": "CommandName", "payload": {...}}`)),
		),

		html.H2(Classes{"text-xl": true, "font-semibold": true, "mt-6": true, "mb-2": true}, g.Text("Available Queries")),
		html.Ul(
			Classes{"list-disc": true, "list-inside": true, "space-y-1": true, "mb-4": true},
			g.Group(queryListItems),
		),
		html.P(
			g.Text("Execute queries via: "),
			html.Code(Classes{"bg-gray-200": true, "p-1": true, "rounded": true}, g.Text("GET /queries/{QueryName}?param1=value1&...")),
		),
	)
}

// HandleIndex creates an http.HandlerFunc for the index page.
// It fetches registered command/query names and renders the IndexPage component wrapped in the main Layout.
func HandleIndex(commandRegistry *CommandRegistry, queryRegistry *QueryRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch registered names
		commandNames := commandRegistry.RegisteredCommandNames()
		queryNames := queryRegistry.RegisteredQueryNames()

		component := IndexPage(commandNames, queryNames)
		layout := Layout("Home - petrock_example_project_name", component) // Use project name in title

		// Set content type and render
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := layout.Render(w)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			// Consider logging the error: slog.Error("Failed to render index page", "error", err)
		}
	}
}
