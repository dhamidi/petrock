package core

import (
	"net/http" // Added for http.ResponseWriter, http.Request

	g "github.com/maragudk/gomponents"            // Use canonical import path
	. "github.com/maragudk/gomponents/components" // Dot-import for helpers like Classes
	"github.com/maragudk/gomponents/html"
)

// IndexPage renders the main content for the home page.
// This component is intended to be passed to the Layout function.
func IndexPage() g.Node {
	return Page("Welcome!", // Use the Page component from core/view.go
		html.P(
			Classes{"text-lg": true, "mb-4": true}, // Correct map literal syntax
			g.Text("Welcome to your new Petrock-generated application!"),
		),
		html.P(
			g.Text("You can start by adding features using: "),
			html.Code(Classes{"bg-gray-200": true, "p-1": true, "rounded": true}, g.Text("petrock feature <feature_name>")), // Correct map literal syntax
		),
		// Add more introductory content or links here
	)
}

// HandleIndex creates an http.HandlerFunc for the index page.
// It renders the IndexPage component wrapped in the main Layout.
func HandleIndex( /* Pass dependencies like QueryRegistry if needed */ ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here you might fetch data using the QueryRegistry if the index page needs dynamic content

		component := IndexPage()
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
