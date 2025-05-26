package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleFormControlsDetail handles the form controls component detail page
func HandleFormControlsDetail(w http.ResponseWriter, r *http.Request) {
	// Create demo content showing different form control components
	demoContent := html.Div(
		ui.CSSClass("space-y-8"),
		
		// Header section
		html.Div(
			ui.CSSClass("mb-8"),
			html.H1(
				ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
				g.Text("Form Controls Component"),
			),
			html.P(
				ui.CSSClass("text-lg", "text-gray-600"),
				g.Text("Interactive form controls including checkboxes, radio buttons, and toggle switches with accessibility features."),
			),
		),

		// Checkbox Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Checkbox"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("Checkboxes allow users to select multiple options from a set of choices."),
			),
			
			// Basic checkboxes
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Basic Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-3"),
					ui.Checkbox(ui.CheckboxProps{
						Label: "Subscribe to newsletter",
						Value: "newsletter",
						Name:  "preferences",
					}),
					ui.Checkbox(ui.CheckboxProps{
						Label:   "Enable notifications",
						Value:   "notifications",
						Name:    "preferences",
						Checked: true,
					}),
					ui.Checkbox(ui.CheckboxProps{
						Label:    "Terms and conditions",
						Value:    "terms",
						Name:     "agreements",
						Required: true,
					}),
					ui.Checkbox(ui.CheckboxProps{
						Label:    "Disabled option",
						Value:    "disabled",
						Name:     "preferences",
						Disabled: true,
					}),
				),
			),
			
			// Checkbox group example
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Checkbox Groups"),
				),
				html.FieldSet(
					ui.CSSClass("border", "border-gray-300", "rounded-md", "p-4"),
					html.Legend(
						ui.CSSClass("text-sm", "font-medium", "text-gray-900", "px-2"),
						g.Text("Select your interests"),
					),
					html.Div(
						ui.CSSClass("space-y-2"),
						ui.Checkbox(ui.CheckboxProps{
							Label: "Web Development",
							Value: "web-dev",
							Name:  "interests",
						}),
						ui.Checkbox(ui.CheckboxProps{
							Label:   "Mobile Development",
							Value:   "mobile-dev",
							Name:    "interests",
							Checked: true,
						}),
						ui.Checkbox(ui.CheckboxProps{
							Label: "Data Science",
							Value: "data-science",
							Name:  "interests",
						}),
						ui.Checkbox(ui.CheckboxProps{
							Label:   "DevOps",
							Value:   "devops",
							Name:    "interests",
							Checked: true,
						}),
					),
				),
			),
		),

		// Radio Button Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Radio Buttons"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("Radio buttons allow users to select exactly one option from a set of choices."),
			),
			
			// Basic radio group
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Radio Group Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-6"),
					ui.RadioGroup(ui.RadioGroupProps{
						Name:    "payment-method",
						Value:   "credit-card",
						Label:   "Payment Method",
						Vertical: true,
						Options: []ui.RadioOption{
							{Value: "credit-card", Label: "Credit Card"},
							{Value: "paypal", Label: "PayPal"},
							{Value: "bank-transfer", Label: "Bank Transfer"},
							{Value: "crypto", Label: "Cryptocurrency", Disabled: true},
						},
					}),
					ui.RadioGroup(ui.RadioGroupProps{
						Name:    "shipping-speed",
						Value:   "standard",
						Label:   "Shipping Speed",
						Vertical: false,
						Options: []ui.RadioOption{
							{Value: "standard", Label: "Standard (5-7 days)"},
							{Value: "express", Label: "Express (2-3 days)"},
							{Value: "overnight", Label: "Overnight"},
						},
					}),
				),
			),
		),

		// Toggle Switch Examples
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Toggle Switches"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("Toggle switches are used for binary choices and provide immediate feedback."),
			),
			
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Toggle Examples"),
				),
				html.Div(
					ui.CSSClass("space-y-4"),
					ui.Toggle(ui.ToggleProps{
						Label: "Enable dark mode",
						Name:  "dark-mode",
						Value: "enabled",
					}),
					ui.Toggle(ui.ToggleProps{
						Label:   "Push notifications",
						Name:    "push-notifications",
						Value:   "enabled",
						Checked: true,
					}),
					ui.Toggle(ui.ToggleProps{
						Label: "Email marketing",
						Name:  "email-marketing",
						Value: "enabled",
					}),
					ui.Toggle(ui.ToggleProps{
						Label:    "Beta features (coming soon)",
						Name:     "beta-features",
						Value:    "enabled",
						Disabled: true,
					}),
				),
			),
			
			// Toggle group with descriptions
			html.Div(
				ui.CSSClass("space-y-4", "p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("Settings Panel"),
				),
				html.Div(
					ui.CSSClass("space-y-6"),
					html.Div(
						ui.Toggle(ui.ToggleProps{
							Label:   "Two-factor authentication",
							Name:    "2fa",
							Value:   "enabled",
							Checked: true,
						}),
						html.P(
							ui.CSSClass("text-sm", "text-gray-500", "mt-1", "ml-14"),
							g.Text("Add an extra layer of security to your account"),
						),
					),
					html.Div(
						ui.Toggle(ui.ToggleProps{
							Label: "Activity notifications",
							Name:  "activity-notifications",
							Value: "enabled",
						}),
						html.P(
							ui.CSSClass("text-sm", "text-gray-500", "mt-1", "ml-14"),
							g.Text("Get notified about important account activity"),
						),
					),
					html.Div(
						ui.Toggle(ui.ToggleProps{
							Label: "Marketing emails",
							Name:  "marketing-emails",
							Value: "enabled",
						}),
						html.P(
							ui.CSSClass("text-sm", "text-gray-500", "mt-1", "ml-14"),
							g.Text("Receive updates about new features and promotions"),
						),
					),
				),
			),
		),

		// Combined Form Example
		html.Div(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Complete Form Example"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-6"),
				g.Text("A comprehensive form using all form control types together."),
			),
			
			html.Div(
				ui.CSSClass("p-6", "bg-gray-50", "rounded-lg"),
				html.H3(
					ui.CSSClass("text-lg", "font-semibold", "mb-4"),
					g.Text("User Preferences Form"),
				),
				html.Form(
					ui.CSSClass("space-y-6"),
					
					// Account settings
					html.FieldSet(
						ui.CSSClass("border", "border-gray-300", "rounded-md", "p-4"),
						html.Legend(
							ui.CSSClass("text-sm", "font-medium", "text-gray-900", "px-2"),
							g.Text("Account Settings"),
						),
						html.Div(
							ui.CSSClass("space-y-3"),
							ui.Toggle(ui.ToggleProps{
								Label:   "Public profile",
								Name:    "public-profile",
								Checked: true,
							}),
							ui.Toggle(ui.ToggleProps{
								Label: "Allow search engines to index profile",
								Name:  "seo-indexing",
							}),
						),
					),
					
					// Communication preferences
					ui.RadioGroup(ui.RadioGroupProps{
						Name:    "email-frequency",
						Value:   "weekly",
						Label:   "Email Frequency",
						Vertical: true,
						Options: []ui.RadioOption{
							{Value: "daily", Label: "Daily digest"},
							{Value: "weekly", Label: "Weekly summary"},
							{Value: "monthly", Label: "Monthly newsletter"},
							{Value: "never", Label: "Never"},
						},
					}),
					
					// Content preferences
					html.FieldSet(
						ui.CSSClass("border", "border-gray-300", "rounded-md", "p-4"),
						html.Legend(
							ui.CSSClass("text-sm", "font-medium", "text-gray-900", "px-2"),
							g.Text("Content Interests"),
						),
						html.Div(
							ui.CSSClass("space-y-2"),
							ui.Checkbox(ui.CheckboxProps{
								Label:   "Technology news",
								Value:   "tech-news",
								Name:    "content-types",
								Checked: true,
							}),
							ui.Checkbox(ui.CheckboxProps{
								Label: "Product updates",
								Value: "product-updates",
								Name:  "content-types",
							}),
							ui.Checkbox(ui.CheckboxProps{
								Label: "Community highlights",
								Value: "community",
								Name:  "content-types",
							}),
							ui.Checkbox(ui.CheckboxProps{
								Label:   "Educational content",
								Value:   "education",
								Name:    "content-types",
								Checked: true,
							}),
						),
					),
				),
			),
		),
	)

	// Create page content with proper sidebar navigation
	pageContent := core.Page("Form Controls Component",
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
		"Form Controls Component - UI Gallery",
		pageContent,
	)

	w.Header().Set("Content-Type", "text/html")
	response.Render(w)
}