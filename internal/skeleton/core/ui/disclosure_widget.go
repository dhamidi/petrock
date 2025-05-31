package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// DisclosureWidgetProps defines the properties for the DisclosureWidget component
type DisclosureWidgetProps struct {
	Title       string // The title/trigger text for the disclosure
	DefaultOpen bool   // Whether the widget is open by default
	ID          string // Unique identifier for the widget
}

// DisclosureWidget creates an expandable content widget using CSS-only implementation
func DisclosureWidget(props DisclosureWidgetProps, children ...g.Node) g.Node {
	// Generate unique ID if not provided
	widgetID := props.ID
	if widgetID == "" {
		widgetID = "disclosure-widget"
	}
	
	checkboxID := widgetID + "-checkbox"

	// Build CSS classes for the container
	containerClasses := []string{
		"border", "border-gray-200", "rounded-lg", "bg-white", "shadow-sm",
	}

	// Build CSS classes for the trigger button
	triggerClasses := []string{
		"w-full", "flex", "justify-between", "items-center",
		"px-4", "py-3", "text-left", "font-medium", "text-gray-900",
		"hover:bg-gray-50", "focus:outline-none", "focus:ring-2",
		"focus:ring-blue-500", "focus:ring-offset-2",
		"cursor-pointer", "transition-colors",
	}

	// Build CSS classes for content panel
	contentClasses := []string{
		"overflow-hidden", "transition-all", "duration-300", "ease-in-out",
		"max-h-0", "peer-checked:max-h-screen",
	}

	return html.Div(
		CSSClass(containerClasses...),
		
		// Hidden checkbox for CSS-only interaction
		html.Input(
			html.Type("checkbox"),
			html.ID(checkboxID),
			CSSClass("sr-only", "peer"),
			g.If(props.DefaultOpen, html.Checked()),
		),
		
		// Label acts as the clickable trigger
		html.Label(
			html.For(checkboxID),
			CSSClass(triggerClasses...),
			html.Span(
				g.Text(props.Title),
			),
			// Chevron icon that rotates when open
			html.Span(
				CSSClass(
					"ml-2", "transform", "transition-transform", "duration-200",
					"peer-checked:rotate-180",
				),
				g.Raw(`<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
				</svg>`),
			),
		),
		
		// Content panel that expands/collapses
		html.Div(
			CSSClass(contentClasses...),
			html.Div(
				CSSClass("px-4", "pb-4", "pt-2", "text-gray-700"),
				g.Group(children),
			),
		),
		
		// Add CSS for screen reader accessibility
		g.Raw(`<style>
			.sr-only {
				position: absolute;
				width: 1px;
				height: 1px;
				padding: 0;
				margin: -1px;
				overflow: hidden;
				clip: rect(0, 0, 0, 0);
				white-space: nowrap;
				border: 0;
			}
		</style>`),
	)
}