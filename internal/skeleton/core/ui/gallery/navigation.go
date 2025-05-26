package gallery

import (
	"net/http"
	"strings"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleNavigationDetail handles the Navigation components demo page
func HandleNavigationDetail(w http.ResponseWriter, r *http.Request) {
	// Create demo content showing different navigation components
	demoContent := html.Div(
		ui.CSSClass("space-y-8"),
		
		// Header section
		html.Div(
			ui.CSSClass("mb-8"),
			html.H1(
				ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
				g.Text("Navigation Components"),
			),
			html.P(
				ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
				g.Text("Navigation components provide consistent patterns for moving through your application, including navigation bars, sidebars, tabs, breadcrumbs, and pagination."),
			),
		),
		
		// NavBar Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Navigation Bar"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("A responsive top navigation bar with brand and navigation items:"),
			),
			
			// NavBar examples
			html.Div(
				ui.CSSClass("space-y-4"),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.NavBar(ui.NavBarProps{
						Brand: "My App",
						Items: []ui.NavItem{
							{Label: "Home", Href: "/", Active: true},
							{Label: "Products", Href: "/products"},
							{Label: "About", Href: "/about"},
							{Label: "Contact", Href: "/contact"},
						},
					}),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.NavBar(ui.NavBarProps{
						Brand: "Dashboard",
						Items: []ui.NavItem{
							{Label: "Dashboard", Href: "/dashboard", Active: true},
							{Label: "Users", Href: "/users"},
							{Label: "Reports", Href: "/reports"},
							{Label: "Settings", Href: "/settings", Disabled: true},
						},
					}),
				),
			),
		),
		
		// SideNav Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Side Navigation"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Collapsible sidebar navigation for dashboard-style layouts:"),
			),
			
			// SideNav examples
			html.Div(
				ui.CSSClass("space-y-4"),
				html.Div(
					ui.CSSClass("flex", "h-64", "border", "rounded", "bg-gray-50"),
					html.Div(
						ui.CSSClass("bg-white"),
						ui.SideNav(ui.SideNavProps{
						Items: createNavItemsWithIcons([]string{"Dashboard", "Projects", "Team", "Settings"}, "/dashboard"),
						Collapsed: false,
						}),
					),
					html.Div(
						ui.CSSClass("flex-1", "p-4", "text-gray-500"),
						g.Text("Main content area..."),
					),
				),
				html.Div(
					ui.CSSClass("flex", "h-64", "border", "rounded", "bg-gray-50"),
					html.Div(
						ui.CSSClass("bg-white"),
						ui.SideNav(ui.SideNavProps{
						Items: createNavItemsWithIcons([]string{"Dashboard", "Projects", "Team", "Settings"}, "/dashboard"),
						Collapsed: true,
						}),
					),
					html.Div(
						ui.CSSClass("flex-1", "p-4", "text-gray-500"),
						g.Text("Collapsed sidebar with more space for content..."),
					),
				),
			),
		),
		
		// Tabs Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Tabs"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Accessible tab interface with ARIA roles and keyboard navigation:"),
			),
			
			// Tabs examples
			html.Div(
				ui.CSSClass("space-y-6"),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Tabs(ui.TabsProps{
						ActiveTab: "overview",
						Items: []ui.TabItem{
							{
								ID:    "overview",
								Label: "Overview",
								Content: html.Div(
									ui.CSSClass("p-4", "bg-white", "rounded", "border"),
									g.Text("This is the overview tab content. It contains general information about the topic."),
								),
							},
							{
								ID:    "details",
								Label: "Details",
								Content: html.Div(
									ui.CSSClass("p-4", "bg-white", "rounded", "border"),
									g.Text("This is the details tab content. It provides in-depth information and specifications."),
								),
							},
							{
								ID:    "settings",
								Label: "Settings",
								Content: html.Div(
									ui.CSSClass("p-4", "bg-white", "rounded", "border"),
									g.Text("This is the settings tab content. Configure options and preferences here."),
								),
							},
							{
								ID:       "disabled",
								Label:    "Disabled",
								Disabled: true,
								Content: html.Div(
									g.Text("This tab is disabled."),
								),
							},
						},
					}),
				),
			),
		),
		
		// Breadcrumbs Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Breadcrumbs"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Breadcrumb navigation shows the current page location within a hierarchy:"),
			),
			
			// Breadcrumbs examples
			html.Div(
				ui.CSSClass("space-y-4"),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Breadcrumbs(ui.BreadcrumbsProps{
						Items: []ui.BreadcrumbItem{
							{Label: "Home", Href: "/"},
							{Label: "Products", Href: "/products"},
							{Label: "Laptops", Href: "/products/laptops"},
							{Label: "MacBook Pro", Current: true},
						},
					}),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Breadcrumbs(ui.BreadcrumbsProps{
						Items: []ui.BreadcrumbItem{
							{Label: "Dashboard", Href: "/dashboard"},
							{Label: "Settings", Href: "/settings"},
							{Label: "Account Settings", Current: true},
						},
					}),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Breadcrumbs(ui.BreadcrumbsProps{
						Items: []ui.BreadcrumbItem{
							{Label: "Home", Href: "/"},
							{Label: "Blog", Href: "/blog"},
							{Label: "Technology", Href: "/blog/technology"},
							{Label: "Web Development", Href: "/blog/technology/web-dev"},
							{Label: "Current Article", Current: true},
						},
					}),
				),
			),
		),
		
		// Pagination Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Pagination"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Pagination component for navigating through large datasets:"),
			),
			
			// Pagination examples
			html.Div(
				ui.CSSClass("space-y-4"),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Pagination(ui.PaginationProps{
						CurrentPage: 1,
						TotalPages:  10,
						BaseURL:     "/products?page=",
						ShowEnds:    true,
						MaxVisible:  7,
					}),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Pagination(ui.PaginationProps{
						CurrentPage: 5,
						TotalPages:  10,
						BaseURL:     "/articles?page=",
						ShowEnds:    true,
						MaxVisible:  7,
					}),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Pagination(ui.PaginationProps{
						CurrentPage: 10,
						TotalPages:  10,
						BaseURL:     "/users?page=",
						ShowEnds:    false,
						MaxVisible:  5,
					}),
				),
			),
		),
		
		// Code Examples
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Usage Examples"),
			),
			html.Pre(
				ui.CSSClass("bg-gray-100", "p-4", "rounded", "text-sm", "overflow-x-auto"),
				g.Text(`// Navigation Bar
ui.NavBar(ui.NavBarProps{
    Brand: "My App",
    Items: []ui.NavItem{
        {Label: "Home", Href: "/", Active: true},
        {Label: "Products", Href: "/products"},
        {Label: "About", Href: "/about"},
    },
})

// Side Navigation
ui.SideNav(ui.SideNavProps{
    Items: []ui.NavItem{
        {Label: "Dashboard", Href: "/dashboard", Active: true},
        {Label: "Settings", Href: "/settings"},
    },
    Collapsed: false,
})

// Tabs
ui.Tabs(ui.TabsProps{
    ActiveTab: "overview",
    Items: []ui.TabItem{
        {
            ID: "overview",
            Label: "Overview", 
            Content: html.Div(g.Text("Overview content")),
        },
        {
            ID: "details",
            Label: "Details",
            Content: html.Div(g.Text("Details content")),
        },
    },
})

// Breadcrumbs
ui.Breadcrumbs(ui.BreadcrumbsProps{
    Items: []ui.BreadcrumbItem{
        {Label: "Home", Href: "/"},
        {Label: "Products", Href: "/products"},
        {Label: "Current Page", Current: true},
    },
})

// Pagination
ui.Pagination(ui.PaginationProps{
    CurrentPage: 5,
    TotalPages: 20,
    BaseURL: "/items?page=",
    ShowEnds: true,
    MaxVisible: 7,
})`),
			),
		),
		
		// Properties documentation
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Component Properties"),
			),
			
			// NavBar Properties
			html.Div(
				ui.CSSClass("mb-6"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("NavBar Properties"),
				),
				html.Div(
					ui.CSSClass("border", "rounded", "overflow-hidden"),
					html.Table(
						ui.CSSClass("w-full"),
						html.THead(
							ui.CSSClass("bg-gray-50"),
							html.Tr(
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Property")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Type")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Description")),
							),
						),
						html.TBody(
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Brand")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Brand text or logo")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Items")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("[]NavItem")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Navigation items with Label, Href, Active, Disabled")),
							),
						),
					),
				),
			),
			
			// Pagination Properties
			html.Div(
				ui.CSSClass("mb-6"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("Pagination Properties"),
				),
				html.Div(
					ui.CSSClass("border", "rounded", "overflow-hidden"),
					html.Table(
						ui.CSSClass("w-full"),
						html.THead(
							ui.CSSClass("bg-gray-50"),
							html.Tr(
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Property")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Type")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Default")),
								html.Th(ui.CSSClass("px-4", "py-2", "text-left"), g.Text("Description")),
							),
						),
						html.TBody(
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("CurrentPage")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("int")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("-")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Current page number (1-based)")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("TotalPages")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("int")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("-")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Total number of pages")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("BaseURL")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("-")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Base URL for pagination links")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("MaxVisible")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("int")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("7")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Maximum number of visible page links")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("ShowEnds")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("bool")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("false")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Whether to show first/last page links")),
							),
						),
					),
				),
			),
		),
	)

	// Create page content with proper sidebar navigation
	pageContent := core.Page("Navigation Components",
		html.Div(
			ui.CSSClass("flex", "min-h-screen", "-mx-4", "-mt-4"),
			// Sidebar with full component navigation
			html.Nav(
				ui.CSSClass("w-64", "bg-white", "border-r", "border-gray-200", "p-6", "overflow-y-auto"),
				html.H1(
					ui.CSSClass("text-lg", "font-semibold", "text-gray-900", "mb-6"),
					g.Text("Components"),
				),
				g.Group(BuildSidebar()),
			),
			// Main content
			html.Main(
				ui.CSSClass("flex-1", "p-6", "overflow-y-auto"),
				html.Div(
					ui.CSSClass("max-w-4xl"),
					demoContent,
				),
			),
		),
	)

	// Use existing Layout function
	response := core.Layout(
		"Navigation Components - UI Gallery",
		pageContent,
	)

	w.Header().Set("Content-Type", "text/html")
	response.Render(w)
}

// createNavItemsWithIcons creates navigation items with geometric shape icons for demo purposes
func createNavItemsWithIcons(labels []string, activeHref string) []ui.NavItem {
	items := make([]ui.NavItem, len(labels))
	
	for i, label := range labels {
		href := "/" + strings.ToLower(label)
		items[i] = ui.NavItem{
			Label:    label,
			Href:     href,
			Active:   href == activeHref,
			Disabled: false,
			Icon:     createSideNavIcon(label),
		}
	}
	
	return items
}

// createSideNavIcon creates a unique geometric shape icon based on the label for gallery demos
func createSideNavIcon(label string) g.Node {
	// Use the first character of the label to determine which shape to use
	var shape g.Node
	
	firstChar := ""
	if len(label) > 0 {
		firstChar = string(label[0])
	}
	
	switch firstChar {
	case "D": // Dashboard - Square
		shape = g.El("rect",
			g.Attr("x", "4"),
			g.Attr("y", "4"),
			g.Attr("width", "12"),
			g.Attr("height", "12"),
		)
	case "P": // Projects - Triangle
		shape = g.El("polygon",
			g.Attr("points", "10,3 17,16 3,16"),
		)
	case "T": // Team/Tasks - Diamond
		shape = g.El("polygon",
			g.Attr("points", "10,2 18,10 10,18 2,10"),
		)
	case "S": // Settings - Hexagon
		shape = g.El("polygon",
			g.Attr("points", "10,1 16,4.5 16,11.5 10,15 4,11.5 4,4.5"),
		)
	case "R": // Reports - Circle
		shape = g.El("circle",
			g.Attr("cx", "10"),
			g.Attr("cy", "10"),
			g.Attr("r", "7"),
		)
	case "U": // Users - Octagon
		shape = g.El("polygon",
			g.Attr("points", "6,2 14,2 18,6 18,14 14,18 6,18 2,14 2,6"),
		)
	case "A": // Analytics/Admin - Star
		shape = g.El("polygon",
			g.Attr("points", "10,1 12,7 19,7 13.5,11 15.5,18 10,14 4.5,18 6.5,11 1,7 8,7"),
		)
	case "H": // Home - House shape
		shape = g.Group([]g.Node{
			g.El("polygon", g.Attr("points", "10,3 18,9 18,17 2,17 2,9")),
			g.El("rect", g.Attr("x", "7"), g.Attr("y", "12"), g.Attr("width", "6"), g.Attr("height", "5")),
		})
	case "N": // Notifications - Bell shape
		shape = g.El("path",
			g.Attr("d", "M6 8c0-3.3 2.7-6 6-6s6 2.7 6 6v4c0 1.1.9 2 2 2H4c1.1 0 2-.9 2-2V8z"),
		)
	default: // Default - Rounded square
		shape = g.El("rect",
			g.Attr("x", "4"),
			g.Attr("y", "4"),
			g.Attr("width", "12"),
			g.Attr("height", "12"),
			g.Attr("rx", "2"),
		)
	}
	
	return g.El("svg",
		ui.CSSClass("flex-shrink-0", "h-5", "w-5"),
		g.Attr("fill", "currentColor"),
		g.Attr("viewBox", "0 0 20 20"),
		g.Attr("aria-hidden", "true"),
		shape,
	)
}