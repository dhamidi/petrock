package pages

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"github.com/petrock/example_module_path/core/ui"
)

// ItemView renders the HTML representation of a single item.
// Adapt the fields and structure based on the 'Result' type in queries.go.
func ItemView(item Result) g.Node {
	return ui.Container(ui.ContainerProps{Variant: "default"},
		// Navigation breadcrumbs
		ui.Breadcrumbs(ui.BreadcrumbsProps{
			Items: []ui.BreadcrumbItem{
				{Label: "Items", Href: "/petrock_example_feature_name"},
				{Label: item.Name, Current: true},
			},
		}),

		ui.Section(ui.SectionProps{Heading: item.Name, Level: 1},
			// Main content in a card
			ui.Card(ui.CardProps{Variant: "default", Padding: "large"},
				ui.CardHeader(
					html.Div(
						ui.CSSClass("flex", "items-center", "justify-between"),
						html.H2(ui.CSSClass("text-xl", "font-semibold"), g.Text("Item Details")),
						ui.Badge(ui.BadgeProps{
							Variant: "info",
							Size:    "medium",
						}, g.Textf("Version %d", item.Version)),
					),
				),
				ui.CardBody(
					// Item details in a two-column grid
					ui.Grid(ui.GridProps{
						Columns: "repeat(auto-fit, minmax(250px, 1fr))",
						Gap:     "1.5rem",
					},
						// Left column - Basic Info
						html.Div(
							ui.CSSClass("space-y-4"),
							html.H3(ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"), g.Text("Basic Information")),
							
							html.Div(
								html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("ID")),
								html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900", "font-mono"), g.Text(item.ID)),
							),
							html.Div(
								html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Name")),
								html.Dd(ui.CSSClass("mt-1", "text-lg", "font-semibold", "text-gray-900"), g.Text(item.Name)),
							),
							html.Div(
								html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Description")),
								html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900"), g.Text(item.Description)),
							),
						),

						// Right column - Timestamps
						html.Div(
							ui.CSSClass("space-y-4"),
							html.H3(ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"), g.Text("Timestamps")),
							
							html.Div(
								html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Created")),
								html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900"), g.Text(item.CreatedAt.Format("Jan 2, 2006 at 15:04"))),
							),
							html.Div(
								html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Last Updated")),
								html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900"), g.Text(item.UpdatedAt.Format("Jan 2, 2006 at 15:04"))),
							),
						),
					),

					ui.Divider(ui.DividerProps{Variant: "default", Margin: "large"}),

					// Content section
					html.Div(
						ui.CSSClass("space-y-4"),
						html.H3(ui.CSSClass("text-lg", "font-medium", "text-gray-900"), g.Text("Content")),
						html.Pre(
							ui.CSSClass("text-sm", "text-gray-900", "whitespace-pre-wrap", "p-4", "bg-gray-50", "rounded-lg", "border"),
							g.Text(item.Content),
						),
					),

					// Summary section (if available)
					func() g.Node {
						if item.Summary == "" {
							return nil
						}
						return html.Div(
							ui.CSSClass("space-y-4", "mt-6"),
							html.H3(ui.CSSClass("text-lg", "font-medium", "text-gray-900"), g.Text("AI Summary")),
							html.Div(
								ui.CSSClass("p-4", "bg-blue-50", "rounded-lg", "border", "border-blue-200"),
								html.P(ui.CSSClass("text-sm", "italic", "text-blue-800"), g.Text(item.Summary)),
							),
						)
					}(),
				),
				ui.CardFooter(
					ui.ButtonGroup(ui.ButtonGroupProps{
						Orientation: "horizontal",
						Spacing:     "medium",
					},
						html.A(
							html.Href("/petrock_example_feature_name/"+item.ID+"/edit"),
							ui.Button(ui.ButtonProps{
								Variant: "primary",
								Size:    "medium",
							}, g.Text("Edit Item")),
						),
						html.A(
							html.Href("/petrock_example_feature_name/"+item.ID+"/delete"),
							ui.Button(ui.ButtonProps{
								Variant: "danger",
								Size:    "medium",
							}, g.Text("Delete Item")),
						),
						html.A(
							html.Href("/petrock_example_feature_name"),
							ui.Button(ui.ButtonProps{
								Variant: "secondary",
								Size:    "medium",
							}, g.Text("Back to List")),
						),
					),
				),
			),
		),
	)
}
