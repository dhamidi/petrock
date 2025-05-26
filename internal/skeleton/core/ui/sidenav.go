package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// SideNavProps defines the properties for the SideNav component
type SideNavProps struct {
	Items     []NavItem // Navigation items (reusing NavItem from navbar.go)
	Collapsed bool      // Whether the sidebar is collapsed
}

// SideNav creates a collapsible sidebar navigation component
func SideNav(props SideNavProps) g.Node {
	// Base container classes
	var containerClasses []string
	if props.Collapsed {
		containerClasses = append(containerClasses, "w-16")
	} else {
		containerClasses = append(containerClasses, "w-64")
	}
	
	containerClasses = append(containerClasses,
		"bg-white", "border-r", "border-gray-200", "flex", "flex-col",
		"transition-all", "duration-300", "ease-in-out",
	)

	return html.Nav(
		CSSClass(containerClasses...),
		html.Div(
			CSSClass("flex-1", "flex", "flex-col", "min-h-0"),
			html.Div(
				CSSClass("flex-1", "flex", "flex-col", "pt-5", "pb-4", "overflow-y-auto"),
				html.Nav(
					CSSClass("mt-5", "flex-1", "px-2", "space-y-1"),
					g.Group(buildSideNavItems(props.Items, props.Collapsed)),
				),
			),
		),
	)
}

// buildSideNavItems creates sidebar navigation item elements
func buildSideNavItems(items []NavItem, collapsed bool) []g.Node {
	var navItems []g.Node
	
	for _, item := range items {
		var classes []string
		classes = append(classes,
			"group", "flex", "items-center", "px-2", "py-2", "text-sm",
			"font-medium", "rounded-md", "transition-colors", "duration-200",
		)

		if item.Disabled {
			classes = append(classes, "text-gray-400", "cursor-not-allowed")
		} else if item.Active {
			classes = append(classes,
				"bg-blue-100", "text-blue-900", "border-r-2", "border-blue-500",
			)
		} else {
			classes = append(classes,
				"text-gray-600", "hover:bg-gray-50", "hover:text-gray-900",
			)
		}

		var attrs []g.Node
		attrs = append(attrs, CSSClass(classes...))
		
		if !item.Disabled {
			attrs = append(attrs, html.Href(item.Href))
		}

		if item.Disabled {
			attrs = append(attrs, html.Aria("disabled", "true"))
		}

		// Add title attribute for collapsed state
		if collapsed {
			attrs = append(attrs, html.Title(item.Label))
		}

		// Create icon placeholder (in a real implementation, you'd pass icons)
		var iconClasses []string
		iconClasses = append(iconClasses, "flex-shrink-0", "h-5", "w-5")
		if !collapsed {
			iconClasses = append(iconClasses, "mr-3")
		}

		var children []g.Node
		
		// Add icon from NavItem if provided, otherwise use a default placeholder
		if item.Icon != nil {
			// Wrap the provided icon in an SVG container with proper classes
			children = append(children, g.El("div",
				CSSClass(iconClasses...),
				item.Icon,
			))
		} else {
			// Default placeholder icon
			children = append(children, g.El("svg",
				CSSClass(iconClasses...),
				g.Attr("fill", "currentColor"),
				g.Attr("viewBox", "0 0 20 20"),
				g.Attr("aria-hidden", "true"),
				g.El("rect",
					g.Attr("x", "4"),
					g.Attr("y", "4"),
					g.Attr("width", "12"),
					g.Attr("height", "12"),
					g.Attr("rx", "2"),
				),
			))
		}

		// Add label text (hidden when collapsed)
		if !collapsed {
			children = append(children, g.Text(item.Label))
		}

		navItems = append(navItems, html.A(
			append(attrs, children...)...,
		))
	}

	return navItems
}