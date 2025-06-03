package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core/ui"
	localUI "github.com/petrock/example_module_path/petrock_example_feature_name/ui"
)

// ItemsListView renders a list of items, potentially with pagination.
func ItemsListView(result ListResult) g.Node {
	// Determine total number of pages
	totalPages := 1
	if result.PageSize > 0 {
		totalPages = (result.TotalCount + result.PageSize - 1) / result.PageSize
	}

	return ui.Container(ui.ContainerProps{Variant: "default"},
		ui.Section(ui.SectionProps{Heading: "Items", Level: 1},
			// Stats and action bar
			html.Div(
				ui.CSSClass("flex", "flex-col", "sm:flex-row", "justify-between", "items-start", "sm:items-center", "gap-4", "mb-6"),

				// Stats badge
				ui.Badge(ui.BadgeProps{
					Variant: "info",
					Size:    "large",
				}, g.Textf("%d item%s", result.TotalCount, localUI.Pluralize(result.TotalCount))),

				// Create new button
				html.A(
					html.Href("/petrock_example_feature_name/new"),
					ui.Button(ui.ButtonProps{
						Variant: "primary",
						Size:    "medium",
					}, g.Text("New Item")),
				),
			),

			// Item list content
			func() g.Node {
				// Empty state
				if len(result.Items) == 0 {
					return ui.Card(ui.CardProps{Variant: "default", Padding: "large"},
						html.Div(
							ui.CSSClass("text-center", "py-12"),
							html.H3(
								ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-2"),
								g.Text("No items yet"),
							),
							html.P(
								ui.CSSClass("text-gray-500", "mb-6"),
								g.Text("Get started by creating your first item."),
							),
							html.A(
								html.Href("/petrock_example_feature_name/new"),
								ui.Button(ui.ButtonProps{
									Variant: "primary",
									Size:    "medium",
								}, g.Text("Create New Item")),
							),
						),
					)
				}

				// Items grid using card components
				items := make([]g.Node, 0, len(result.Items))
				for _, item := range result.Items {
					items = append(items, ui.Card(ui.CardProps{Variant: "default", Padding: "medium"},
						ui.CardHeader(
							html.H3(
								ui.CSSClass("text-lg", "font-medium"),
								html.A(
									html.Href("/petrock_example_feature_name/"+item.ID),
									ui.CSSClass("text-indigo-600", "hover:text-indigo-800"),
									g.Text(item.Name),
								),
							),
							ui.Badge(ui.BadgeProps{
								Variant: "secondary",
								Size:    "small",
							}, g.Text(item.CreatedAt.Format("Jan 2, 2006"))),
						),
						ui.CardBody(
							html.P(
								ui.CSSClass("text-gray-700", "text-sm", "mb-4"),
								g.Text(item.Description),
							),
							// Content preview
							func() g.Node {
								if item.Content == "" {
									return nil
								}
								// Truncate content for preview (first 150 characters)
								contentPreview := item.Content
								if len(contentPreview) > 150 {
									contentPreview = contentPreview[:150] + "..."
								}
								return html.Div(
									ui.CSSClass("border-t", "border-gray-100", "pt-3", "mb-3"),
									html.H4(
										ui.CSSClass("text-xs", "font-medium", "text-gray-500", "mb-1"),
										g.Text("Content"),
									),
									html.P(
										ui.CSSClass("text-sm", "text-gray-600", "bg-gray-50", "p-2", "rounded", "whitespace-pre-wrap"),
										g.Text(contentPreview),
									),
								)
							}(),
							// Summary (if available)
							func() g.Node {
								if item.Summary == "" {
									return nil
								}
								return html.Div(
									ui.CSSClass("border-t", "border-gray-100", "pt-3"),
									html.H4(
										ui.CSSClass("text-xs", "font-medium", "text-gray-500", "mb-1"),
										g.Text("Summary"),
									),
									html.P(
										ui.CSSClass("text-sm", "italic", "text-gray-600"),
										g.Text(item.Summary),
									),
								)
							}(),
						),
						ui.CardFooter(
							ui.ButtonGroup(ui.ButtonGroupProps{
								Orientation: "horizontal",
								Spacing:     "small",
							},
								html.A(
									html.Href("/petrock_example_feature_name/"+item.ID),
									ui.Button(ui.ButtonProps{
										Variant: "secondary",
										Size:    "small",
									}, g.Text("View")),
								),
								html.A(
									html.Href("/petrock_example_feature_name/"+item.ID+"/edit"),
									ui.Button(ui.ButtonProps{
										Variant: "secondary",
										Size:    "small",
									}, g.Text("Edit")),
								),
								html.A(
									html.Href("/petrock_example_feature_name/"+item.ID+"/delete"),
									ui.Button(ui.ButtonProps{
										Variant: "danger",
										Size:    "small",
									}, g.Text("Delete")),
								),
							),
						),
					))
				}

				return ui.Grid(ui.GridProps{
					Columns: "repeat(auto-fit, minmax(320px, 1fr))",
					Gap:     "1.5rem",
				}, g.Group(items))
		}(),

			// Pagination controls (if more than one page)
			func() g.Node {
				if totalPages <= 1 {
					return nil
				}

				return ui.Pagination(ui.PaginationProps{
					CurrentPage: result.Page,
					TotalPages:  totalPages,
					BaseURL:     fmt.Sprintf("/petrock_example_feature_name/?pageSize=%d&page=", result.PageSize),
					ShowEnds:    true,
					MaxVisible:  7,
				})
			}(),
		),
	)
}
