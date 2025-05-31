package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"strconv"
)

// AccordionItem represents a single item in an accordion
type AccordionItem struct {
	Title   string
	Content []g.Node
	ID      string // Unique identifier for the item
}

// AccordionProps defines the properties for the Accordion component
type AccordionProps struct {
	Items         []AccordionItem
	AllowMultiple bool // whether multiple items can be open at once
}

// Accordion creates an expandable content component using CSS-only implementation
func Accordion(props AccordionProps) g.Node {
	if len(props.Items) == 0 {
		return html.Div()
	}

	var accordionItems []g.Node

	for i, item := range props.Items {
		itemID := item.ID
		if itemID == "" {
			itemID = "accordion-item-" + strconv.Itoa(i)
		}

		// CSS classes for the accordion item
		itemClasses := []string{
			"border-b", "border-gray-200",
		}

		// Add bottom border except for last item
		if i == len(props.Items)-1 {
			itemClasses = append(itemClasses, "border-b-0")
		}

		var inputType string
		var inputName string
		if props.AllowMultiple {
			inputType = "checkbox"
			inputName = itemID + "-checkbox"
		} else {
			inputType = "radio"
			inputName = "accordion-group"
		}

		accordionItem := html.Div(
			CSSClass(itemClasses...),
			
			// Hidden input for CSS-only interaction
			html.Input(
				html.Type(inputType),
				html.Name(inputName),
				html.ID(itemID),
				CSSClass("sr-only", "peer"),
			),
			
			// Label acts as the clickable header
			html.Label(
				html.For(itemID),
				CSSClass(
					"flex", "justify-between", "items-center",
					"w-full", "py-4", "px-6", "cursor-pointer",
					"text-left", "font-medium", "text-gray-900",
					"hover:bg-gray-50", "transition-colors",
					"focus-within:outline-none", "focus-within:ring-2",
					"focus-within:ring-blue-500", "focus-within:ring-offset-2",
				),
				html.Span(
					g.Text(item.Title),
				),
				// Chevron icon
				html.Span(
					CSSClass(
						"ml-2", "transform", "transition-transform",
						"peer-checked:rotate-180",
					),
					g.Raw(`<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
					</svg>`),
				),
			),
			
			// Content panel
			html.Div(
				CSSClass(
					"overflow-hidden", "transition-all", "duration-300",
					"max-h-0", "peer-checked:max-h-screen",
				),
				html.Div(
					CSSClass("px-6", "pb-4", "text-gray-700"),
					g.Group(item.Content),
				),
			),
		)

		accordionItems = append(accordionItems, accordionItem)
	}

	return html.Div(
		CSSClass(
			"bg-white", "border", "border-gray-200", "rounded-lg",
			"shadow-sm", "overflow-hidden",
		),
		g.Group(accordionItems),
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