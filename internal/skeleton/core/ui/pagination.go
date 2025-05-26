package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"strconv"
)

// PaginationProps defines the properties for the Pagination component
type PaginationProps struct {
	CurrentPage int    // Current page number (1-based)
	TotalPages  int    // Total number of pages
	BaseURL     string // Base URL for pagination links (page number will be appended)
	ShowEnds    bool   // Whether to show first/last page links
	MaxVisible  int    // Maximum number of page links to show (default: 7)
}

// Pagination creates an accessible pagination component with page numbers and navigation
func Pagination(props PaginationProps) g.Node {
	if props.TotalPages <= 1 {
		return nil
	}

	// Set defaults
	maxVisible := props.MaxVisible
	if maxVisible == 0 {
		maxVisible = 7
	}

	return html.Nav(
		CSSClass("flex", "items-center", "justify-between"),
		html.Aria("label", "Pagination Navigation"),
		
		// Previous/Next info text
		html.Div(
			CSSClass("flex-1", "flex", "justify-between", "sm:hidden"),
			g.Group(buildMobilePagination(props)),
		),
		
		// Desktop pagination
		html.Div(
			CSSClass("hidden", "sm:flex-1", "sm:flex", "sm:items-center", "sm:justify-between"),
			html.Div(
				CSSClass("text-sm", "text-gray-700"),
				g.Text("Showing page "),
				html.Span(
					CSSClass("font-medium"),
					g.Text(strconv.Itoa(props.CurrentPage)),
				),
				g.Text(" of "),
				html.Span(
					CSSClass("font-medium"),
					g.Text(strconv.Itoa(props.TotalPages)),
				),
			),
			html.Div(
				html.Nav(
					CSSClass("relative", "z-0", "inline-flex", "rounded-md", "shadow-sm", "-space-x-px"),
					html.Aria("label", "Pagination"),
					g.Group(buildPaginationLinks(props, maxVisible)),
				),
			),
		),
	)
}

// buildMobilePagination creates mobile-friendly previous/next buttons
func buildMobilePagination(props PaginationProps) []g.Node {
	var buttons []g.Node

	// Previous button
	if props.CurrentPage > 1 {
		prevURL := props.BaseURL + strconv.Itoa(props.CurrentPage-1)
		buttons = append(buttons, html.A(
			html.Href(prevURL),
			CSSClass(
				"relative", "inline-flex", "items-center", "px-4", "py-2",
				"border", "border-gray-300", "text-sm", "font-medium",
				"rounded-md", "text-gray-700", "bg-white", "hover:bg-gray-50",
			),
			g.Text("Previous"),
		))
	} else {
		buttons = append(buttons, html.Span(
			CSSClass(
				"relative", "inline-flex", "items-center", "px-4", "py-2",
				"border", "border-gray-300", "text-sm", "font-medium",
				"rounded-md", "text-gray-400", "bg-gray-100", "cursor-not-allowed",
			),
			g.Text("Previous"),
		))
	}

	// Next button
	if props.CurrentPage < props.TotalPages {
		nextURL := props.BaseURL + strconv.Itoa(props.CurrentPage+1)
		buttons = append(buttons, html.A(
			html.Href(nextURL),
			CSSClass(
				"ml-3", "relative", "inline-flex", "items-center", "px-4", "py-2",
				"border", "border-gray-300", "text-sm", "font-medium",
				"rounded-md", "text-gray-700", "bg-white", "hover:bg-gray-50",
			),
			g.Text("Next"),
		))
	} else {
		buttons = append(buttons, html.Span(
			CSSClass(
				"ml-3", "relative", "inline-flex", "items-center", "px-4", "py-2",
				"border", "border-gray-300", "text-sm", "font-medium",
				"rounded-md", "text-gray-400", "bg-gray-100", "cursor-not-allowed",
			),
			g.Text("Next"),
		))
	}

	return buttons
}

// buildPaginationLinks creates the full pagination link structure
func buildPaginationLinks(props PaginationProps, maxVisible int) []g.Node {
	var links []g.Node

	// Previous button
	if props.CurrentPage > 1 {
		prevURL := props.BaseURL + strconv.Itoa(props.CurrentPage-1)
		links = append(links, html.A(
			html.Href(prevURL),
			CSSClass(
				"relative", "inline-flex", "items-center", "px-2", "py-2",
				"rounded-l-md", "border", "border-gray-300", "bg-white",
				"text-sm", "font-medium", "text-gray-500", "hover:bg-gray-50",
			),
			html.Aria("label", "Previous page"),
			html.Span(CSSClass("sr-only"), g.Text("Previous")),
			createChevronLeft(),
		))
	} else {
		links = append(links, html.Span(
			CSSClass(
				"relative", "inline-flex", "items-center", "px-2", "py-2",
				"rounded-l-md", "border", "border-gray-300", "bg-gray-100",
				"text-sm", "font-medium", "text-gray-400", "cursor-not-allowed",
			),
			html.Aria("label", "Previous page"),
			createChevronLeft(),
		))
	}

	// Page numbers
	links = append(links, buildPageNumbers(props, maxVisible)...)

	// Next button
	if props.CurrentPage < props.TotalPages {
		nextURL := props.BaseURL + strconv.Itoa(props.CurrentPage+1)
		links = append(links, html.A(
			html.Href(nextURL),
			CSSClass(
				"relative", "inline-flex", "items-center", "px-2", "py-2",
				"rounded-r-md", "border", "border-gray-300", "bg-white",
				"text-sm", "font-medium", "text-gray-500", "hover:bg-gray-50",
			),
			html.Aria("label", "Next page"),
			html.Span(CSSClass("sr-only"), g.Text("Next")),
			createChevronRight(),
		))
	} else {
		links = append(links, html.Span(
			CSSClass(
				"relative", "inline-flex", "items-center", "px-2", "py-2",
				"rounded-r-md", "border", "border-gray-300", "bg-gray-100",
				"text-sm", "font-medium", "text-gray-400", "cursor-not-allowed",
			),
			createChevronRight(),
		))
	}

	return links
}

// buildPageNumbers creates the numbered page links with ellipsis
func buildPageNumbers(props PaginationProps, maxVisible int) []g.Node {
	var pages []g.Node
	
	// Calculate which pages to show
	start := 1
	end := props.TotalPages
	
	if props.TotalPages > maxVisible {
		half := maxVisible / 2
		start = props.CurrentPage - half
		end = props.CurrentPage + half
		
		if start < 1 {
			start = 1
			end = maxVisible
		}
		if end > props.TotalPages {
			end = props.TotalPages
			start = props.TotalPages - maxVisible + 1
		}
	}

	// Add first page and ellipsis if needed
	if props.ShowEnds && start > 1 {
		pages = append(pages, createPageLink(1, props.BaseURL, false))
		if start > 2 {
			pages = append(pages, createEllipsis())
		}
	}

	// Add page numbers
	for i := start; i <= end; i++ {
		isCurrent := i == props.CurrentPage
		pages = append(pages, createPageLink(i, props.BaseURL, isCurrent))
	}

	// Add last page and ellipsis if needed
	if props.ShowEnds && end < props.TotalPages {
		if end < props.TotalPages-1 {
			pages = append(pages, createEllipsis())
		}
		pages = append(pages, createPageLink(props.TotalPages, props.BaseURL, false))
	}

	return pages
}

// createPageLink creates a single page number link
func createPageLink(page int, baseURL string, isCurrent bool) g.Node {
	pageStr := strconv.Itoa(page)
	
	var classes []string
	classes = append(classes,
		"relative", "inline-flex", "items-center", "px-4", "py-2",
		"border", "text-sm", "font-medium",
	)

	if isCurrent {
		classes = append(classes,
			"z-10", "bg-blue-50", "border-blue-500", "text-blue-600",
		)
		return html.Span(
			CSSClass(classes...),
			html.Aria("current", "page"),
			g.Text(pageStr),
		)
	} else {
		classes = append(classes,
			"bg-white", "border-gray-300", "text-gray-500", "hover:bg-gray-50",
		)
		return html.A(
			html.Href(baseURL+pageStr),
			CSSClass(classes...),
			html.Aria("label", "Page "+pageStr),
			g.Text(pageStr),
		)
	}
}

// createEllipsis creates an ellipsis indicator
func createEllipsis() g.Node {
	return html.Span(
		CSSClass(
			"relative", "inline-flex", "items-center", "px-4", "py-2",
			"border", "border-gray-300", "bg-white", "text-sm",
			"font-medium", "text-gray-700",
		),
		g.Text("..."),
	)
}

// createChevronLeft creates a left-pointing chevron icon
func createChevronLeft() g.Node {
	return g.El("svg",
		CSSClass("h-5", "w-5"),
		g.Attr("viewBox", "0 0 20 20"),
		g.Attr("fill", "currentColor"),
		g.Attr("aria-hidden", "true"),
		g.El("path",
			g.Attr("fill-rule", "evenodd"),
			g.Attr("d", "M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"),
			g.Attr("clip-rule", "evenodd"),
		),
	)
}

// createChevronRight creates a right-pointing chevron icon
func createChevronRight() g.Node {
	return g.El("svg",
		CSSClass("h-5", "w-5"),
		g.Attr("viewBox", "0 0 20 20"),
		g.Attr("fill", "currentColor"),
		g.Attr("aria-hidden", "true"),
		g.El("path",
			g.Attr("fill-rule", "evenodd"),
			g.Attr("d", "M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"),
			g.Attr("clip-rule", "evenodd"),
		),
	)
}