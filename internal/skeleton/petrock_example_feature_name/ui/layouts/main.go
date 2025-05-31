package layouts

import (
	g "maragu.dev/gomponents"
	"github.com/petrock/example_module_path/core/ui"
)

// FeatureLayout renders a consistent layout for feature pages with navigation
func FeatureLayout(title string, content g.Node) g.Node {
	return ui.Layout(title, content)
}

// DashboardLayout renders a layout with sidebar navigation for admin/dashboard pages
func DashboardLayout(title string, content g.Node) g.Node {
	return ui.Layout(title,
		ui.Container(ui.ContainerProps{Variant: "wide"},
			ui.Grid(ui.GridProps{
				Columns: "250px 1fr",
				Gap:     "2rem",
			},
				// Sidebar
				ui.Card(ui.CardProps{Variant: "default", Padding: "medium"},
					ui.SideNav(ui.SideNavProps{
						Items: []ui.NavItem{
							{Label: "Dashboard", Href: "/admin", Icon: g.Text("📊")},
							{Label: "Items", Href: "/petrock_example_feature_name", Icon: g.Text("📝")},
							{Label: "Settings", Href: "/admin/settings", Icon: g.Text("⚙️")},
						},
					}),
				),
				// Main content
				content,
			),
		),
	)
}
