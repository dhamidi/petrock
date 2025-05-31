package core

import (
	"net/http"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"github.com/petrock/example_module_path/core/ui"
)

// IndexPage renders the main content for the home page, displaying registered commands and queries.
// This component uses the new UI components from the design system.
func IndexPage(commandNames, queryNames []string) g.Node {
	commandListItems := make([]g.Node, len(commandNames))
	for i, name := range commandNames {
		commandListItems[i] = html.Li(html.Code(g.Text(name)))
	}

	queryListItems := make([]g.Node, len(queryNames))
	for i, name := range queryNames {
		queryListItems[i] = html.Li(html.Code(g.Text(name)))
	}

	return ui.Container(ui.ContainerProps{Variant: "default"},
		ui.Section(ui.SectionProps{Heading: "Welcome!", Level: 1},
			html.P(
				ui.CSSClass("text-lg", "mb-4"),
				g.Text("Welcome to your Petrock-generated application!"),
			),
		),

		ui.Section(ui.SectionProps{Heading: "Available Commands", Level: 2},
			html.Ul(
				ui.CSSClass("list-disc", "list-inside", "space-y-1", "mb-4"),
				g.Group(commandListItems),
			),
			html.P(
				g.Text("Execute commands via: "),
				html.Code(ui.CSSClass("bg-gray-200", "p-1", "rounded"), g.Text("POST /commands")),
				g.Text(" with JSON payload: "),
				html.Code(ui.CSSClass("bg-gray-200", "p-1", "rounded"), g.Text(`{"type": "CommandName", "payload": {...}}`)),
			),
		),

		ui.Section(ui.SectionProps{Heading: "Available Queries", Level: 2},
			html.Ul(
				ui.CSSClass("list-disc", "list-inside", "space-y-1", "mb-4"),
				g.Group(queryListItems),
			),
			html.P(
				g.Text("Execute queries via: "),
				html.Code(ui.CSSClass("bg-gray-200", "p-1", "rounded"), g.Text("GET /queries/{QueryName}?param1=value1&...")),
			),
		),
	)
}

// HandleIndex creates an http.HandlerFunc for the index page.
// It fetches registered command/query names and renders the IndexPage component wrapped in the UI Layout.
func HandleIndex(commandRegistry *CommandRegistry, queryRegistry *QueryRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch registered names
		commandNames := commandRegistry.RegisteredCommandNames()
		queryNames := queryRegistry.RegisteredQueryNames()

		component := IndexPage(commandNames, queryNames)
		layout := ui.Layout("Home - petrock_example_project_name", component)

		// Set content type and render
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := layout.Render(w)
		if err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			// Consider logging the error: slog.Error("Failed to render index page", "error", err)
		}
	}
}
