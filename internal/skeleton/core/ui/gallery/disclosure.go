package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleDisclosureDetail handles the Disclosure component demo page
func HandleDisclosureDetail(w http.ResponseWriter, r *http.Request) {
	// Create demo content showing different disclosure components
	demoContent := html.Div(
		ui.CSSClass("space-y-8"),
		
		// Header section
		html.Div(
			ui.CSSClass("mb-8"),
			html.H1(
				ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
				g.Text("Disclosure Components"),
			),
			html.P(
				ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
				g.Text("Disclosure components allow users to show and hide content sections using CSS-only implementations for better performance and accessibility."),
			),
		),
		
		// Disclosure Widget Examples
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Disclosure Widget"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Simple expandable content sections that can be toggled individually:"),
			),
			
			// Basic disclosure widgets
			html.Div(
				ui.CSSClass("space-y-4", "p-4", "border", "rounded", "bg-gray-50"),
				ui.DisclosureWidget(
					ui.DisclosureWidgetProps{
						Title: "What is a disclosure widget?",
						ID:    "disclosure-1",
					},
					html.P(
						ui.CSSClass("text-gray-700"),
						g.Text("A disclosure widget is a user interface component that allows users to show or hide content by clicking on a trigger element. It's commonly used for FAQ sections, help text, and detailed information that users may not always need to see."),
					),
				),
				ui.DisclosureWidget(
					ui.DisclosureWidgetProps{
						Title:       "How does CSS-only implementation work?",
						DefaultOpen: true,
						ID:          "disclosure-2",
					},
					html.P(
						ui.CSSClass("text-gray-700"),
						g.Text("CSS-only disclosure uses hidden checkbox inputs and the :checked pseudo-class to control visibility. This approach provides smooth animations and accessibility without requiring JavaScript, making it lighter and more performant."),
					),
				),
				ui.DisclosureWidget(
					ui.DisclosureWidgetProps{
						Title: "When should I use disclosure widgets?",
						ID:    "disclosure-3",
					},
					html.Div(
						ui.CSSClass("text-gray-700", "space-y-2"),
						html.P(g.Text("Disclosure widgets are ideal for:")),
						html.Ul(
							ui.CSSClass("list-disc", "list-inside", "space-y-1", "ml-4"),
							html.Li(g.Text("FAQ sections")),
							html.Li(g.Text("Progressive disclosure of complex information")),
							html.Li(g.Text("Help documentation")),
							html.Li(g.Text("Settings panels")),
							html.Li(g.Text("Content that not all users need to see immediately")),
						),
					),
				),
			),
		),
		
		// Accordion Examples
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Accordion"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Accordion groups multiple disclosure sections together with consistent styling:"),
			),
			
			// Single-open accordion
			html.Div(
				ui.CSSClass("mb-6"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("Single Selection (Radio Behavior)"),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Accordion(ui.AccordionProps{
						AllowMultiple: false,
						Items: []ui.AccordionItem{
							{
								Title: "Getting Started",
								ID:    "accordion-single-1",
								Content: []g.Node{
									html.P(
										ui.CSSClass("mb-3"),
										g.Text("Welcome to our platform! This guide will help you get started with the basics."),
									),
									html.P(
										g.Text("First, make sure you have created an account and verified your email address. Then you can explore the dashboard and available features."),
									),
								},
							},
							{
								Title: "Account Settings",
								ID:    "accordion-single-2",
								Content: []g.Node{
									html.P(
										ui.CSSClass("mb-3"),
										g.Text("Manage your account preferences, security settings, and profile information."),
									),
									html.Ul(
										ui.CSSClass("list-disc", "list-inside", "space-y-1"),
										html.Li(g.Text("Update your profile picture")),
										html.Li(g.Text("Change your password")),
										html.Li(g.Text("Configure notification preferences")),
										html.Li(g.Text("Manage connected applications")),
									),
								},
							},
							{
								Title: "Billing & Subscriptions",
								ID:    "accordion-single-3",
								Content: []g.Node{
									html.P(
										ui.CSSClass("mb-3"),
										g.Text("View your current subscription, payment history, and billing information."),
									),
									html.P(
										g.Text("You can upgrade or downgrade your plan at any time. Changes will be prorated based on your billing cycle."),
									),
								},
							},
						},
					}),
				),
			),
			
			// Multi-open accordion
			html.Div(
				ui.CSSClass("mb-6"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("Multiple Selection (Checkbox Behavior)"),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Accordion(ui.AccordionProps{
						AllowMultiple: true,
						Items: []ui.AccordionItem{
							{
								Title: "Personal Information",
								ID:    "accordion-multi-1",
								Content: []g.Node{
									html.P(
										g.Text("Update your personal details including name, email, phone number, and address. This information is used for account verification and communication."),
									),
								},
							},
							{
								Title: "Privacy Settings",
								ID:    "accordion-multi-2",
								Content: []g.Node{
									html.P(
										ui.CSSClass("mb-3"),
										g.Text("Control who can see your information and how it's used:"),
									),
									html.Ul(
										ui.CSSClass("list-disc", "list-inside", "space-y-1"),
										html.Li(g.Text("Profile visibility settings")),
										html.Li(g.Text("Data sharing preferences")),
										html.Li(g.Text("Cookie and tracking options")),
									),
								},
							},
							{
								Title: "Notification Preferences",
								ID:    "accordion-multi-3",
								Content: []g.Node{
									html.P(
										g.Text("Choose how and when you receive notifications. You can customize settings for email, SMS, and in-app notifications separately."),
									),
								},
							},
						},
					}),
				),
			),
		),
		
		// Interactive Examples
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Interactive Examples"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Common use cases and patterns for disclosure components:"),
			),
			
			// FAQ section
			html.Div(
				ui.CSSClass("mb-6"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("FAQ Section"),
				),
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					ui.Accordion(ui.AccordionProps{
						AllowMultiple: true,
						Items: []ui.AccordionItem{
							{
								Title: "How do I reset my password?",
								ID:    "faq-1",
								Content: []g.Node{
									html.P(
										g.Text("Click the 'Forgot Password' link on the login page, enter your email address, and follow the instructions in the reset email we send you."),
									),
								},
							},
							{
								Title: "Can I cancel my subscription at any time?",
								ID:    "faq-2",
								Content: []g.Node{
									html.P(
										g.Text("Yes, you can cancel your subscription at any time from your account settings. Your access will continue until the end of your current billing period."),
									),
								},
							},
							{
								Title: "Is my data secure?",
								ID:    "faq-3",
								Content: []g.Node{
									html.P(
										ui.CSSClass("mb-2"),
										g.Text("We take security seriously and implement multiple layers of protection:"),
									),
									html.Ul(
										ui.CSSClass("list-disc", "list-inside", "space-y-1"),
										html.Li(g.Text("End-to-end encryption for all data")),
										html.Li(g.Text("Regular security audits and penetration testing")),
										html.Li(g.Text("SOC 2 Type II compliance")),
										html.Li(g.Text("24/7 monitoring and threat detection")),
									),
								},
							},
						},
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
				g.Text(`// Basic disclosure widget
ui.DisclosureWidget(
    ui.DisclosureWidgetProps{
        Title: "Click to expand",
        ID:    "unique-id",
    },
    html.P(g.Text("Hidden content goes here...")),
)

// Disclosure widget open by default
ui.DisclosureWidget(
    ui.DisclosureWidgetProps{
        Title:       "Already expanded",
        DefaultOpen: true,
        ID:          "open-widget",
    },
    html.Div(g.Text("Content visible by default")),
)

// Simple accordion (single selection)
ui.Accordion(ui.AccordionProps{
    AllowMultiple: false,
    Items: []ui.AccordionItem{
        {
            Title: "Section 1",
            ID:    "section-1",
            Content: []g.Node{
                html.P(g.Text("First section content")),
            },
        },
        {
            Title: "Section 2", 
            ID:    "section-2",
            Content: []g.Node{
                html.P(g.Text("Second section content")),
            },
        },
    },
})

// Multi-selection accordion (checkbox behavior)
ui.Accordion(ui.AccordionProps{
    AllowMultiple: true,
    Items: []ui.AccordionItem{
        // Multiple items can be open at once
        {Title: "Item 1", ID: "item-1", Content: []g.Node{...}},
        {Title: "Item 2", ID: "item-2", Content: []g.Node{...}},
    },
})`),
			),
		),
		
		// Properties documentation
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Properties"),
			),
			
			// DisclosureWidget properties
			html.Div(
				ui.CSSClass("mb-8"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("DisclosureWidget Properties"),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Title")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("The clickable title/trigger text")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("DefaultOpen")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("bool")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("false")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Whether the widget is expanded by default")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("ID")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"disclosure-widget\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Unique identifier for the widget")),
							),
						),
					),
				),
			),
			
			// Accordion properties
			html.Div(
				ui.CSSClass("mb-8"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("Accordion Properties"),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Items")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("[]AccordionItem")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("[]")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Array of accordion items to display")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("AllowMultiple")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("bool")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("false")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Whether multiple items can be open simultaneously")),
							),
						),
					),
				),
			),
			
			// AccordionItem properties
			html.Div(
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "text-gray-900", "mb-3"),
					g.Text("AccordionItem Properties"),
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
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Title")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("The title text for the accordion item")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("Content")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("[]g.Node")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("nil")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("The content nodes to display when expanded")),
							),
							html.Tr(
								ui.CSSClass("border-t"),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("ID")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("string")),
								html.Td(ui.CSSClass("px-4", "py-2", "font-mono", "text-sm"), g.Text("\"accordion-item-{index}\"")),
								html.Td(ui.CSSClass("px-4", "py-2"), g.Text("Unique identifier for the accordion item")),
							),
						),
					),
				),
			),
		),
	)

	// Create page content with proper sidebar navigation
	pageContent := ui.Page("Disclosure Components",
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
	response := ui.Layout(
		"Disclosure Components - UI Gallery",
		pageContent,
	)

	w.Header().Set("Content-Type", "text/html")
	response.Render(w)
}