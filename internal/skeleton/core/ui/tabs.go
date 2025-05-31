package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// TabItem represents an item in a tab interface
type TabItem struct {
	ID       string // Unique identifier for the tab
	Label    string // Display label for the tab
	Content  g.Node // Content to display when tab is active
	Disabled bool   // Whether the tab is disabled
}

// TabsProps defines the properties for the Tabs component
type TabsProps struct {
	Items     []TabItem // Tab items
	ActiveTab string    // ID of the currently active tab
}

// Tabs creates an accessible tab interface with ARIA roles and keyboard navigation
func Tabs(props TabsProps) g.Node {
	// Set default active tab if not specified
	activeTab := props.ActiveTab
	if activeTab == "" && len(props.Items) > 0 {
		// Find first non-disabled tab
		for _, item := range props.Items {
			if !item.Disabled {
				activeTab = item.ID
				break
			}
		}
	}

	return html.Div(
		CSSClass("w-full"),
		// Tab list
		html.Div(
			CSSClass("border-b", "border-gray-200"),
			html.Nav(
				CSSClass("-mb-px", "flex", "space-x-8"),
				html.Aria("label", "Tabs"),
				html.Role("tablist"),
				g.Group(buildTabHeaders(props.Items, activeTab)),
			),
		),
		// Tab panels
		html.Div(
			CSSClass("mt-4"),
			g.Group(buildTabPanels(props.Items, activeTab)),
		),
	)
}

// buildTabHeaders creates tab header buttons
func buildTabHeaders(items []TabItem, activeTab string) []g.Node {
	var tabs []g.Node

	for _, item := range items {
		isActive := item.ID == activeTab
		
		var classes []string
		classes = append(classes,
			"whitespace-nowrap", "py-2", "px-1", "border-b-2", "font-medium",
			"text-sm", "transition-colors", "duration-200", "focus:outline-none",
			"focus:ring-2", "focus:ring-blue-500", "focus:ring-offset-2",
		)

		if item.Disabled {
			classes = append(classes,
				"text-gray-400", "cursor-not-allowed", "border-transparent",
			)
		} else if isActive {
			classes = append(classes,
				"border-blue-500", "text-blue-600",
			)
		} else {
			classes = append(classes,
				"border-transparent", "text-gray-500", "hover:text-gray-700",
				"hover:border-gray-300", "cursor-pointer",
			)
		}

		var attrs []g.Node
		attrs = append(attrs,
			CSSClass(classes...),
			html.Role("tab"),
			html.Aria("controls", "panel-"+item.ID),
			html.ID("tab-"+item.ID),
		)

		if isActive {
			attrs = append(attrs, html.Aria("selected", "true"))
		} else {
			attrs = append(attrs, html.Aria("selected", "false"))
		}

		if item.Disabled {
			attrs = append(attrs, html.Aria("disabled", "true"))
		} else {
			if isActive {
				attrs = append(attrs, html.TabIndex("0"))
			} else {
				attrs = append(attrs, html.TabIndex("-1"))
			}
		}

		tabs = append(tabs, html.Button(
			append(attrs, g.Text(item.Label))...,
		))
	}

	return tabs
}

// buildTabPanels creates tab content panels
func buildTabPanels(items []TabItem, activeTab string) []g.Node {
	var panels []g.Node

	for _, item := range items {
		isActive := item.ID == activeTab

		var attrs []g.Node
		attrs = append(attrs,
			html.ID("panel-"+item.ID),
			html.Role("tabpanel"),
			html.Aria("labelledby", "tab-"+item.ID),
			html.TabIndex("0"),
		)

		// Hide inactive panels
		if !isActive {
			attrs = append(attrs, CSSClass("hidden"))
		}

		panels = append(panels, html.Div(
			append(attrs, item.Content)...,
		))
	}

	return panels
}

// TabsStyles provides CSS for enhanced tab functionality
const TabsStyles = `
/* Focus visible styles for better accessibility */
.tab-button:focus-visible {
	outline: 2px solid #3b82f6;
	outline-offset: 2px;
}

/* Smooth transitions for tab switching */
.tab-panel {
	animation: fadeIn 0.2s ease-in-out;
}

@keyframes fadeIn {
	from { opacity: 0; }
	to { opacity: 1; }
}

/* Enhanced hover states */
.tab-button:not([aria-disabled="true"]):hover {
	border-color: #d1d5db;
}

.tab-button[aria-selected="true"] {
	font-weight: 600;
}
`