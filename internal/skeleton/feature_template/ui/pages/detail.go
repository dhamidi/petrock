package pages

import (
	g "maragu.dev/gomponents"
	// . "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"
)

// ItemView renders the HTML representation of a single item.
// Adapt the fields and structure based on the 'Result' type in queries.go.
func ItemView(item Result) g.Node {
	return html.Div(
		// Item information card
		html.Div(
			g.Attr("class", "space-y-6"),

			// Item metadata in a grid layout - responsive
			html.Dl(
				g.Attr("class", "grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-6"),

				// ID field
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("ID")),
					html.Dd(g.Attr("class", "mt-1 text-sm text-slate-900"), g.Text(item.ID)),
				),

				// Name field
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Name")),
					html.Dd(g.Attr("class", "mt-1 text-sm text-slate-900 font-medium"), g.Text(item.Name)),
				),

				// Description field - spans full width
				html.Div(
					g.Attr("class", "col-span-1 sm:col-span-2"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Description")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-900 whitespace-pre-wrap"),
						g.Text(item.Description),
					),
				),

				// Content field - spans full width
				html.Div(
					g.Attr("class", "col-span-1 sm:col-span-2"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Content")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-900 whitespace-pre-wrap p-3 bg-slate-50 rounded border border-slate-100"),
						g.Text(item.Content),
					),
				),

				// Summary field (if available) - spans full width
				func() g.Node {
					if item.Summary == "" {
						return nil
					}
					return html.Div(
						g.Attr("class", "col-span-1 sm:col-span-2"),
						html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Summary")),
						html.Dd(
							g.Attr("class", "mt-1 text-sm italic text-slate-700 whitespace-pre-wrap p-3 bg-indigo-50 rounded border border-indigo-100"),
							g.Text(item.Summary),
						),
					)
				}(),

				// Created date
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Created")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-500"),
						g.Text(item.CreatedAt.Format("Jan 2, 2006 at 15:04")),
					),
				),

				// Updated date
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Last Updated")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-500"),
						g.Text(item.UpdatedAt.Format("Jan 2, 2006 at 15:04")),
					),
				),

				// Version number
				html.Div(
					g.Attr("class", "col-span-1 sm:col-span-2"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Version")),
					html.Dd(g.Attr("class", "mt-1 text-sm text-slate-500"), g.Textf("%d", item.Version)),
				),
			),
		),

		// Actions section - responsive
		html.Div(
			g.Attr("class", "mt-8 flex flex-col sm:flex-row gap-3"),
			// Edit button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/edit"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:w-auto w-full"),
				g.Text("Edit Item"),
			),
			// Delete button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/delete"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:w-auto w-full"),
				g.Text("Delete Item"),
			),
			// Back to list
			html.A(
				g.Attr("href", "/petrock_example_feature_name"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-slate-300 text-sm font-medium rounded-md shadow-sm text-slate-700 bg-white hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:w-auto w-full"),
				g.Text("Back to List"),
			),
		),
	)
}
