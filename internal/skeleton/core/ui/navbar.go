package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// NavItem represents an item in a navigation bar
type NavItem struct {
	Label    string
	Href     string
	Active   bool
	Disabled bool
	Icon     g.Node // Optional icon for the navigation item
}

// NavBarProps defines the properties for the NavBar component
type NavBarProps struct {
	Brand string    // Brand text or logo
	Items []NavItem // Navigation items
}

// NavBar creates a responsive navigation bar with brand and navigation items
func NavBar(props NavBarProps) g.Node {
	return html.Nav(
		CSSClass("bg-white", "shadow-sm", "border-b", "border-gray-200"),
		html.Div(
			CSSClass("max-w-7xl", "mx-auto", "px-4", "sm:px-6", "lg:px-8"),
			html.Div(
				CSSClass("flex", "justify-between", "h-16"),
				html.Div(
					CSSClass("flex"),
					// Brand/Logo section
					html.Div(
						CSSClass("flex-shrink-0", "flex", "items-center"),
						html.A(
							html.Href("/"),
							CSSClass("text-xl", "font-bold", "text-gray-900", "hover:text-gray-700"),
							g.Text(props.Brand),
						),
					),
					// Navigation items
					html.Div(
						CSSClass("hidden", "sm:ml-6", "sm:flex", "sm:space-x-8"),
						g.Group(buildNavItems(props.Items)),
					),
				),
				// Mobile menu button (placeholder for future enhancement)
				html.Div(
					CSSClass("-mr-2", "flex", "items-center", "sm:hidden"),
					html.Button(
						html.Type("button"),
						CSSClass(
							"bg-white", "inline-flex", "items-center", "justify-center",
							"p-2", "rounded-md", "text-gray-400", "hover:text-gray-500",
							"hover:bg-gray-100", "focus:outline-none", "focus:ring-2",
							"focus:ring-inset", "focus:ring-blue-500",
						),
						html.Aria("expanded", "false"),
						html.Span(
							CSSClass("sr-only"),
							g.Text("Open main menu"),
						),
						// Hamburger icon
						g.El("svg",
						CSSClass("block", "h-6", "w-6"),
						g.Attr("fill", "none"),
						 g.Attr("viewBox", "0 0 24 24"),
						g.Attr("stroke", "currentColor"),
						g.Attr("aria-hidden", "true"),
						g.El("path",
							g.Attr("stroke-linecap", "round"),
							g.Attr("stroke-linejoin", "round"),
							g.Attr("stroke-width", "2"),
							g.Attr("d", "M4 6h16M4 12h16M4 18h16"),
						),
					),
					),
				),
			),
		),
	)
}

// buildNavItems creates navigation item elements
func buildNavItems(items []NavItem) []g.Node {
	var navItems []g.Node
	
	for _, item := range items {
		var classes []string
		classes = append(classes,
			"inline-flex", "items-center", "px-1", "pt-1", "text-sm", "font-medium",
			"transition-colors", "duration-200",
		)

		if item.Disabled {
			classes = append(classes, "text-gray-400", "cursor-not-allowed")
		} else if item.Active {
			classes = append(classes,
				"border-blue-500", "text-gray-900", "border-b-2",
			)
		} else {
			classes = append(classes,
				"border-transparent", "text-gray-500", "hover:text-gray-700",
				"hover:border-gray-300", "border-b-2",
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

		navItems = append(navItems, html.A(
			append(attrs, g.Text(item.Label))...,
		))
	}

	return navItems
}