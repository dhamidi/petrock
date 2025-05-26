package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// BreadcrumbItem represents an item in a breadcrumb trail
type BreadcrumbItem struct {
	Label    string // Display text for the breadcrumb
	Href     string // URL for the breadcrumb (empty for current page)
	Current  bool   // Whether this is the current page
	Disabled bool   // Whether the breadcrumb is disabled
}

// BreadcrumbsProps defines the properties for the Breadcrumbs component
type BreadcrumbsProps struct {
	Items []BreadcrumbItem // Breadcrumb items in order from root to current
}

// Breadcrumbs creates an accessible breadcrumb navigation component
func Breadcrumbs(props BreadcrumbsProps) g.Node {
	if len(props.Items) == 0 {
		return nil
	}

	return html.Nav(
		CSSClass("flex"),
		html.Aria("label", "Breadcrumb"),
		html.Ol(
			html.Role("list"),
			CSSClass("inline-flex", "items-center", "space-x-1", "md:space-x-3"),
			g.Group(buildBreadcrumbItems(props.Items)),
		),
	)
}

// buildBreadcrumbItems creates breadcrumb item elements
func buildBreadcrumbItems(items []BreadcrumbItem) []g.Node {
	var breadcrumbs []g.Node

	for i, item := range items {
		// Create list item
		var listItemAttrs []g.Node
		listItemAttrs = append(listItemAttrs, CSSClass("inline-flex", "items-center"))

		var listItemChildren []g.Node

		// Add separator for non-first items
		if i > 0 {
			listItemChildren = append(listItemChildren, g.El("svg",
				CSSClass("flex-shrink-0", "w-5", "h-5", "text-gray-400", "mx-1"),
				g.Attr("fill", "currentColor"),
				g.Attr("viewBox", "0 0 20 20"),
				g.Attr("aria-hidden", "true"),
				g.El("path",
					g.Attr("fill-rule", "evenodd"),
					g.Attr("d", "M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"),
					g.Attr("clip-rule", "evenodd"),
				),
			))
		}

		// Create breadcrumb link or text
		if item.Current || item.Href == "" {
			// Current page - just text, no link
			var textClasses []string
			if item.Current {
				textClasses = append(textClasses, "ml-1", "text-sm", "font-medium", "text-gray-500", "md:ml-2")
			} else {
				textClasses = append(textClasses, "ml-1", "text-sm", "font-medium", "text-gray-400", "md:ml-2")
			}
			
			if item.Disabled {
				textClasses = append(textClasses, "opacity-50")
			}

			var textAttrs []g.Node
			textAttrs = append(textAttrs, CSSClass(textClasses...))
			
			if item.Current {
				textAttrs = append(textAttrs, html.Aria("current", "page"))
			}

			listItemChildren = append(listItemChildren, html.Span(
				append(textAttrs, g.Text(item.Label))...,
			))
		} else {
			// Link to previous page
			var linkClasses []string
			linkClasses = append(linkClasses, 
				"ml-1", "text-sm", "font-medium", "transition-colors", "duration-200", "md:ml-2",
			)
			
			if item.Disabled {
				linkClasses = append(linkClasses, "text-gray-400", "cursor-not-allowed")
			} else {
				linkClasses = append(linkClasses, "text-gray-500", "hover:text-gray-700")
			}

			var linkAttrs []g.Node
			linkAttrs = append(linkAttrs, CSSClass(linkClasses...))
			
			if !item.Disabled {
				linkAttrs = append(linkAttrs, html.Href(item.Href))
			}

			if item.Disabled {
				linkAttrs = append(linkAttrs, html.Aria("disabled", "true"))
			}

			listItemChildren = append(listItemChildren, html.A(
				append(linkAttrs, g.Text(item.Label))...,
			))
		}

		breadcrumbs = append(breadcrumbs, html.Li(
			append(listItemAttrs, listItemChildren...)...,
		))
	}

	return breadcrumbs
}