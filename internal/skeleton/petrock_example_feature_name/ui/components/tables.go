package components

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"github.com/petrock/example_module_path/core/ui"
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
)

// ItemsTable renders a responsive table of items
func ItemsTable(items []state.Item) g.Node {
	if len(items) == 0 {
		return ui.Alert(ui.AlertProps{
			Type:    "info",
			Title:   "No Items",
			Message: "No items have been created yet. Get started by creating your first item.",
		})
	}

	// Build table rows
	rows := make([]g.Node, len(items))
	for i, item := range items {
		rows[i] = html.Tr(
			ui.CSSClass("hover:bg-gray-50"),
			html.Td(
				ui.CSSClass("px-6", "py-4", "whitespace-nowrap", "text-sm", "font-medium", "text-gray-900"),
				html.A(
					html.Href("/petrock_example_feature_name/"+item.ID),
					ui.CSSClass("text-indigo-600", "hover:text-indigo-900"),
					g.Text(item.Name),
				),
			),
			html.Td(
				ui.CSSClass("px-6", "py-4", "text-sm", "text-gray-500", "max-w-xs", "truncate"),
				g.Text(item.Description),
			),
			html.Td(
				ui.CSSClass("px-6", "py-4", "whitespace-nowrap", "text-sm", "text-gray-500"),
				g.Text(item.CreatedAt.Format("Jan 2, 2006")),
			),
			html.Td(
				ui.CSSClass("px-6", "py-4", "whitespace-nowrap", "text-right", "text-sm", "font-medium"),
				ui.ButtonGroup(ui.ButtonGroupProps{
					Orientation: "horizontal",
					Spacing:     "small",
				},
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
		)
	}

	return ui.Card(ui.CardProps{Variant: "default", Padding: "none"},
		html.Div(ui.CSSClass("overflow-hidden"),
			html.Table(ui.CSSClass("min-w-full", "divide-y", "divide-gray-200"),
				html.THead(ui.CSSClass("bg-gray-50"),
					html.Tr(
						html.Th(
							ui.CSSClass("px-6", "py-3", "text-left", "text-xs", "font-medium", "text-gray-500", "uppercase", "tracking-wider"),
							g.Text("Name"),
						),
						html.Th(
							ui.CSSClass("px-6", "py-3", "text-left", "text-xs", "font-medium", "text-gray-500", "uppercase", "tracking-wider"),
							g.Text("Description"),
						),
						html.Th(
							ui.CSSClass("px-6", "py-3", "text-left", "text-xs", "font-medium", "text-gray-500", "uppercase", "tracking-wider"),
							g.Text("Created"),
						),
						html.Th(
							ui.CSSClass("relative", "px-6", "py-3"),
							html.Span(ui.CSSClass("sr-only"), g.Text("Actions")),
						),
					),
				),
				html.TBody(ui.CSSClass("bg-white", "divide-y", "divide-gray-200"),
					g.Group(rows),
				),
			),
		),
	)
}
