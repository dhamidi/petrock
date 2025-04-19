package pages

import (
	"fmt"

	// "github.com/petrock/example_module_path/core"
	g "maragu.dev/gomponents"
	// . "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/petrock_example_feature_name/ui"
)

// ItemsListView renders a list of items, potentially with pagination.
func ItemsListView(result ListResult) g.Node {
	// Determine total number of pages
	totalPages := 1
	if result.PageSize > 0 {
		totalPages = (result.TotalCount + result.PageSize - 1) / result.PageSize
	}

	return html.Div(
		g.Attr("class", "space-y-8"),

		// List header with actions
		html.Div(
			g.Attr("class", "flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"),

			// Stats section
			html.Div(
				g.Attr("class", "text-sm text-slate-500"),
				g.Textf(
					"%d item%s • Page %d of %d",
					result.TotalCount,
					ui.Pluralize(result.TotalCount),
					result.Page,
					totalPages,
				),
			),

			// Create new button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/new"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
				g.Text("New Item"),
			),
		),

		// Item list with item cards
		func() g.Node {
			// Empty state
			if len(result.Items) == 0 {
				return html.Div(
					g.Attr("class", "text-center p-12 border border-dashed border-slate-300 rounded-lg"),
					html.Div(
						g.Attr("class", "mx-auto h-12 w-12 text-slate-400 flex items-center justify-center border-2 border-dashed border-slate-300 rounded-full"),
						// Simple placeholder text instead of SVG
						html.Span(
							g.Attr("class", "text-2xl font-light"),
							g.Text("!"),
						),
					),
					html.H3(
						g.Attr("class", "mt-2 text-sm font-medium text-slate-900"),
						g.Text("No items"),
					),
					html.P(
						g.Attr("class", "mt-1 text-sm text-slate-500"),
						g.Text("Get started by creating a new item."),
					),
					html.Div(
						g.Attr("class", "mt-6"),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/new"),
							g.Attr("class", "inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
							g.Text("Create New Item"),
						),
					),
				)
			}

			// Items list in a responsive grid
			items := make([]g.Node, 0, len(result.Items))
			for _, item := range result.Items {
				items = append(items, html.Div(
					g.Attr("class", "bg-white shadow-sm rounded-lg overflow-hidden border border-slate-200 hover:shadow-md transition-shadow duration-200"),

					// Card header with name and created date
					html.Div(
						g.Attr("class", "px-4 py-5 sm:px-6 border-b border-slate-200 bg-slate-50"),
						html.Div(
							g.Attr("class", "flex justify-between items-center flex-wrap gap-2"),
							html.H3(
								g.Attr("class", "text-lg font-medium leading-6 text-slate-900 break-all"),
								html.A(
									g.Attr("href", "/petrock_example_feature_name/"+item.ID),
									g.Attr("class", "hover:text-indigo-600"),
									g.Text(item.Name),
								),
							),
							html.Span(
								g.Attr("class", "inline-flex items-center rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-600"),
								g.Text(item.CreatedAt.Format("Jan 2, 2006")),
							),
						),
					),

					// Card body with description and summary
					html.Div(
						g.Attr("class", "px-4 py-5 sm:p-6 space-y-3"),
						// Description
						html.P(
							g.Attr("class", "text-sm text-slate-700 break-words line-clamp-3"),
							g.Text(item.Description),
						),
						// Summary (if available)
						func() g.Node {
							if item.Summary == "" {
								return nil
							}
							return html.Div(
								g.Attr("class", "mt-2 border-t border-slate-100 pt-2"),
								html.H4(
									g.Attr("class", "text-xs font-medium text-slate-500 mb-1"),
									g.Text("Summary"),
								),
								html.P(
									g.Attr("class", "text-sm italic text-slate-600 break-words"),
									g.Text(item.Summary),
								),
							)
						}(),
					),

					// Card footer with actions
					html.Div(
						g.Attr("class", "border-t border-slate-200 bg-slate-50 px-4 py-4 sm:px-6 flex justify-end space-x-3"),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/edit"),
							g.Attr("class", "text-sm font-medium text-indigo-600 hover:text-indigo-500"),
							g.Text("Edit"),
						),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/delete"),
							g.Attr("class", "text-sm font-medium text-red-600 hover:text-red-500"),
							g.Text("Delete"),
						),
					),
				))
			}

			return html.Div(
				g.Attr("class", "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"),
				g.Group(items),
			)
		}(),

		// Pagination controls (if more than one page)
		func() g.Node {
			if totalPages <= 1 {
				return nil
			}

			// Simple pagination with prev/next buttons
			return html.Nav(
				g.Attr("class", "border-t border-slate-200 px-4 flex items-center justify-between sm:px-0"),

				// Previous page button
				html.Div(func() g.Node {
					if result.Page <= 1 {
						return html.Span(
							g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-300 bg-white border border-slate-300 cursor-not-allowed rounded-md"),
							g.Text("Previous"),
						)
					}

					prevPage := result.Page - 1
					prevURL := fmt.Sprintf("/petrock_example_feature_name/?page=%d&pageSize=%d", prevPage, result.PageSize)
					return html.A(
						g.Attr("href", prevURL),
						g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50"),
						g.Text("Previous"),
					)
				}()),

				// Next page button
				html.Div(func() g.Node {
					if result.Page >= totalPages {
						return html.Span(
							g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-300 bg-white border border-slate-300 cursor-not-allowed rounded-md"),
							g.Text("Next"),
						)
					}

					nextPage := result.Page + 1
					nextURL := fmt.Sprintf("/petrock_example_feature_name/?page=%d&pageSize=%d", nextPage, result.PageSize)
					return html.A(
						g.Attr("href", nextURL),
						g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50"),
						g.Text("Next"),
					)
				}()),
			)
		}(),
	)
}