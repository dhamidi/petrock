package core

import (
	g "maragu.dev/gomponents"            // Use canonical import path
	. "maragu.dev/gomponents/components" // Dot-import for helpers like Classes
	"maragu.dev/gomponents/html"
	// "time" // Removed as it's only used in commented code
)

// Layout renders the full HTML page structure.
// It includes common head elements and wraps the body content.
func Layout(pageTitle string, bodyContent ...g.Node) g.Node {
	return html.HTML( // Corrected casing: html.HTML
		html.Lang("en"),
		html.Head( // Corrected casing: html.Head
			html.Meta(html.Charset("utf-8")), // Corrected casing: html.Meta
			html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1")), // Corrected casing: html.Meta
			html.TitleEl(g.Text(pageTitle)), // TitleEl is correct

			// Link to Tailwind CSS (via CDN for simplicity, replace with local build if needed)
			// Consider using the project's own asset bundling pipeline
			StylesheetLink("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"),

			// Placeholder for other CSS or JS links
			// StylesheetLink("/assets/core/main.css"), // Example for core assets
			// ScriptLink("/assets/core/main.js", false, true), // Example core JS (defer)

			// Include Hotwire/Stimulus if used
			// ScriptLink("https://unpkg.com/@hotwired/turbo@7.1.0/dist/turbo.es2017-umd.js", false, true),
			// ScriptLink("https://unpkg.com/stimulus/dist/stimulus.umd.js", false, true),
		),
		html.Body(
			// Basic body styling (e.g., background color)
			Classes{"bg-gray-100": true, "font-sans": true, "antialiased": true}, // Correct map literal syntax

			// Optional: Add common header/navigation here
			// Navbar(),

			// Main content area
			html.Main(
				g.Group(bodyContent), // Embed the page-specific content
			),

			// Optional: Add common footer here
			// Footer(),

			// Placeholder for page-specific scripts or bottom-of-body includes
		),
	)
}

// --- Optional Common Components ---

// Navbar example (customize with actual links and styling)
// func Navbar() g.Node {
// 	return html.Nav(
// 		Classes{"bg-gray-800": true, "p-4": true, "text-white": true}, // Correct map literal syntax
// 		html.Div(
// 			Classes{"container": true, "mx-auto": true, "flex": true, "justify-between": true}, // Correct map literal syntax
// 			html.A(html.Href("/"), g.Text("petrock_example_project_name")), // Use template var if needed
// 			html.Div(
// 				// Add nav links here
// 				html.A(html.Href("/"), Classes{"px-3": true}, g.Text("Home")), // Correct map literal syntax
// 				// html.A(html.Href("/posts"), Classes{"px-3": true}, g.Text("Posts")), // Example feature link
// 			),
// 		),
// 	)
// }

// Footer example
// func Footer() g.Node {
// 	return html.Footer(
// 		Classes{"bg-gray-200": true, "text-center": true, "p-4": true, "mt-8": true}, // Correct map literal syntax
// 		g.Textf("© %d petrock_example_project_name", time.Now().Year()), // Use template var if needed
// 	)
// }
