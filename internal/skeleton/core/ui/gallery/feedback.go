package gallery

import (
	"net/http"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// HandleFeedbackDetail handles the Feedback components demo page
func HandleFeedbackDetail(w http.ResponseWriter, r *http.Request) {
	// Create demo content showing different feedback components
	demoContent := html.Div(
		ui.CSSClass("space-y-8"),
		
		// Header section
		html.Div(
			ui.CSSClass("mb-8"),
			html.H1(
				ui.CSSClass("text-3xl", "font-bold", "text-gray-900", "mb-4"),
				g.Text("Feedback Components"),
			),
			html.P(
				ui.CSSClass("text-lg", "text-gray-600", "mb-4"),
				g.Text("Feedback components provide visual and accessible feedback to users about the status of operations, notifications, and progress."),
			),
		),
		
		// Alert Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Alert Component"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Alerts display important messages to users with different severity levels."),
			),
			
			// Alert variants
			html.Div(
				ui.CSSClass("space-y-4"),
				// Success alert
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Success Alert"),
					),
					ui.Alert(ui.AlertProps{
						Type:    "success",
						Title:   "Success!",
						Message: "Your changes have been saved successfully.",
					}),
				),
				
				// Warning alert
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Warning Alert"),
					),
					ui.Alert(ui.AlertProps{
						Type:    "warning",
						Title:   "Warning",
						Message: "Please review your settings before proceeding.",
					}),
				),
				
				// Error alert
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Error Alert"),
					),
					ui.Alert(ui.AlertProps{
						Type:    "error",
						Title:   "Error",
						Message: "There was an error processing your request.",
					}),
				),
				
				// Info alert
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Info Alert"),
					),
					ui.Alert(ui.AlertProps{
						Type:    "info",
						Message: "New features are available in the latest update.",
					}),
				),
				
				// Dismissible alert
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-2"),
						g.Text("Dismissible Alert"),
					),
					ui.Alert(ui.AlertProps{
						Type:        "info",
						Title:       "Update Available",
						Message:     "A new version is available for download.",
						Dismissible: true,
					}),
				),
			),
		),
		
		// Badge Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Badge Component"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Badges display small status indicators, counts, or labels."),
			),
			
			// Badge variants
			html.Div(
				ui.CSSClass("space-y-4"),
				// Variant badges
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Badge Variants"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-3", "items-center"),
						ui.Badge(ui.BadgeProps{Variant: "primary"}, g.Text("Primary")),
						ui.Badge(ui.BadgeProps{Variant: "secondary"}, g.Text("Secondary")),
						ui.Badge(ui.BadgeProps{Variant: "success"}, g.Text("Success")),
						ui.Badge(ui.BadgeProps{Variant: "warning"}, g.Text("Warning")),
						ui.Badge(ui.BadgeProps{Variant: "error"}, g.Text("Error")),
						ui.Badge(ui.BadgeProps{Variant: "info"}, g.Text("Info")),
					),
				),
				
				// Badge sizes
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Badge Sizes"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-3", "items-center"),
						ui.Badge(ui.BadgeProps{Size: "small"}, g.Text("Small")),
						ui.Badge(ui.BadgeProps{Size: "medium"}, g.Text("Medium")),
						ui.Badge(ui.BadgeProps{Size: "large"}, g.Text("Large")),
					),
				),
				
				// Count badges
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Count Badges"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-4", "items-center"),
						html.Div(
							ui.CSSClass("flex", "items-center", "gap-2"),
							g.Text("Messages"),
							ui.Badge(ui.BadgeProps{Count: 3}),
						),
						html.Div(
							ui.CSSClass("flex", "items-center", "gap-2"),
							g.Text("Notifications"),
							ui.Badge(ui.BadgeProps{Count: 12, Variant: "error"}),
						),
						html.Div(
							ui.CSSClass("flex", "items-center", "gap-2"),
							g.Text("Updates"),
							ui.Badge(ui.BadgeProps{Count: 142, Variant: "success"}),
						),
					),
				),
				
				// Dot badges
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Dot Badges (Status Indicators)"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-4", "items-center"),
						html.Div(
							ui.CSSClass("flex", "items-center", "gap-2"),
							ui.Badge(ui.BadgeProps{Variant: "success"}),
							g.Text("Online"),
						),
						html.Div(
							ui.CSSClass("flex", "items-center", "gap-2"),
							ui.Badge(ui.BadgeProps{Variant: "warning"}),
							g.Text("Away"),
						),
						html.Div(
							ui.CSSClass("flex", "items-center", "gap-2"),
							ui.Badge(ui.BadgeProps{Variant: "error"}),
							g.Text("Offline"),
						),
					),
				),
			),
		),
		
		// Progress Bar Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Progress Bar Component"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Progress bars show the completion status of tasks or operations."),
			),
			
			html.Div(
				ui.CSSClass("space-y-4"),
				// Basic progress bars
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Basic Progress Bars"),
					),
					html.Div(
						ui.CSSClass("space-y-3"),
						ui.ProgressBar(ui.ProgressBarProps{Value: 25, Max: 100}),
						ui.ProgressBar(ui.ProgressBarProps{Value: 50, Max: 100}),
						ui.ProgressBar(ui.ProgressBarProps{Value: 75, Max: 100}),
						ui.ProgressBar(ui.ProgressBarProps{Value: 100, Max: 100}),
					),
				),
				
				// Labeled progress bars
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Labeled Progress Bars"),
					),
					html.Div(
						ui.CSSClass("space-y-4"),
						ui.ProgressBar(ui.ProgressBarProps{
							Value: 35,
							Max:   100,
							Label: "Upload Progress",
						}),
						ui.ProgressBar(ui.ProgressBarProps{
							Value: 68,
							Max:   100,
							Label: "Installation",
							Color: "success",
						}),
						ui.ProgressBar(ui.ProgressBarProps{
							Value: 15,
							Max:   100,
							Label: "Loading",
							Color: "warning",
						}),
					),
				),
				
				// Different sizes
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Progress Bar Sizes"),
					),
					html.Div(
						ui.CSSClass("space-y-3"),
						html.Div(
							ui.CSSClass("space-y-1"),
							html.P(ui.CSSClass("text-sm", "text-gray-600"), g.Text("Small")),
							ui.ProgressBar(ui.ProgressBarProps{Value: 60, Max: 100, Size: "small"}),
						),
						html.Div(
							ui.CSSClass("space-y-1"),
							html.P(ui.CSSClass("text-sm", "text-gray-600"), g.Text("Medium")),
							ui.ProgressBar(ui.ProgressBarProps{Value: 60, Max: 100, Size: "medium"}),
						),
						html.Div(
							ui.CSSClass("space-y-1"),
							html.P(ui.CSSClass("text-sm", "text-gray-600"), g.Text("Large")),
							ui.ProgressBar(ui.ProgressBarProps{Value: 60, Max: 100, Size: "large"}),
						),
					),
				),
			),
		),
		
		// Loading Spinner Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Loading Spinner Component"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Loading spinners indicate ongoing processes and provide visual feedback during wait times."),
			),
			
			html.Div(
				ui.CSSClass("space-y-4"),
				// Basic spinners
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Basic Spinners"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-6", "items-center"),
						html.Div(
							ui.CSSClass("flex", "flex-col", "items-center", "gap-2"),
							ui.LoadingSpinner(ui.LoadingSpinnerProps{Size: "small"}),
							html.Span(ui.CSSClass("text-xs", "text-gray-500"), g.Text("Small")),
						),
						html.Div(
							ui.CSSClass("flex", "flex-col", "items-center", "gap-2"),
							ui.LoadingSpinner(ui.LoadingSpinnerProps{Size: "medium"}),
							html.Span(ui.CSSClass("text-xs", "text-gray-500"), g.Text("Medium")),
						),
						html.Div(
							ui.CSSClass("flex", "flex-col", "items-center", "gap-2"),
							ui.LoadingSpinner(ui.LoadingSpinnerProps{Size: "large"}),
							html.Span(ui.CSSClass("text-xs", "text-gray-500"), g.Text("Large")),
						),
					),
				),
				
				// Spinner colors
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Spinner Colors"),
					),
					html.Div(
						ui.CSSClass("flex", "flex-wrap", "gap-6", "items-center"),
						html.Div(
							ui.CSSClass("flex", "flex-col", "items-center", "gap-2"),
							ui.LoadingSpinner(ui.LoadingSpinnerProps{Color: "primary"}),
							html.Span(ui.CSSClass("text-xs", "text-gray-500"), g.Text("Primary")),
						),
						html.Div(
							ui.CSSClass("flex", "flex-col", "items-center", "gap-2"),
							ui.LoadingSpinner(ui.LoadingSpinnerProps{Color: "secondary"}),
							html.Span(ui.CSSClass("text-xs", "text-gray-500"), g.Text("Secondary")),
						),
						html.Div(
							ui.CSSClass("flex", "flex-col", "items-center", "gap-2", "bg-gray-800", "p-3", "rounded"),
							ui.LoadingSpinner(ui.LoadingSpinnerProps{Color: "white"}),
							html.Span(ui.CSSClass("text-xs", "text-white"), g.Text("White")),
						),
					),
				),
				
				// Spinner with text
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Spinner with Text"),
					),
					html.Div(
						ui.CSSClass("space-y-3"),
						ui.LoadingSpinnerWithText(ui.LoadingSpinnerProps{
							Label: "Loading data...",
						}),
						ui.LoadingSpinnerWithText(ui.LoadingSpinnerProps{
							Label: "Processing request...",
							Size:  "small",
						}),
						ui.LoadingSpinnerWithText(ui.LoadingSpinnerProps{
							Label: "Uploading file...",
							Color: "secondary",
						}),
					),
				),
			),
		),
		
		// Toast Component
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Toast Component"),
			),
			html.P(
				ui.CSSClass("text-gray-600", "mb-4"),
				g.Text("Toast notifications provide temporary, non-intrusive feedback to users. Note: These examples show the static appearance; in a real application, toasts would be positioned fixed and animate in/out."),
			),
			
			html.Div(
				ui.CSSClass("space-y-4"),
				// Toast variants
				html.Div(
					ui.CSSClass("p-4", "border", "rounded", "bg-gray-50"),
					html.H3(
						ui.CSSClass("text-base", "font-medium", "text-gray-900", "mb-3"),
						g.Text("Toast Variants (Static Examples)"),
					),
					html.Div(
						ui.CSSClass("space-y-3", "relative"),
						// Remove fixed positioning for demo
						html.Div(
							ui.CSSClass("relative", "max-w-sm"),
							ui.Toast(ui.ToastProps{
								Type:    "success",
								Title:   "Success!",
								Message: "Your changes have been saved.",
							}),
						),
						html.Div(
							ui.CSSClass("relative", "max-w-sm"),
							ui.Toast(ui.ToastProps{
								Type:    "warning",
								Title:   "Warning",
								Message: "Your session will expire soon.",
							}),
						),
						html.Div(
							ui.CSSClass("relative", "max-w-sm"),
							ui.Toast(ui.ToastProps{
								Type:    "error",
								Title:   "Error",
								Message: "Failed to save changes.",
							}),
						),
						html.Div(
							ui.CSSClass("relative", "max-w-sm"),
							ui.Toast(ui.ToastProps{
								Type:        "info",
								Title:       "Update Available",
								Message:     "A new version is available.",
								Dismissible: true,
							}),
						),
					),
				),
			),
		),
		
		// Usage Examples
		html.Section(
			ui.CSSClass("space-y-6"),
			html.H2(
				ui.CSSClass("text-2xl", "font-semibold", "text-gray-900", "mb-4"),
				g.Text("Usage Examples"),
			),
			html.Pre(
				ui.CSSClass("bg-gray-100", "p-4", "rounded", "text-sm", "overflow-x-auto"),
				g.Text(`// Alert examples
ui.Alert(ui.AlertProps{
    Type: "success",
    Title: "Success!",
    Message: "Operation completed successfully.",
})

// Badge examples
ui.Badge(ui.BadgeProps{Count: 5})
ui.Badge(ui.BadgeProps{Variant: "error"}, g.Text("Error"))

// Progress bar examples
ui.ProgressBar(ui.ProgressBarProps{
    Value: 75,
    Max: 100,
    Label: "Upload Progress",
    Color: "success",
})

// Loading spinner examples
ui.LoadingSpinner(ui.LoadingSpinnerProps{Size: "large"})
ui.LoadingSpinnerWithText(ui.LoadingSpinnerProps{
    Label: "Loading...",
})

// Toast examples
ui.Toast(ui.ToastProps{
    Type: "info",
    Title: "Notification",
    Message: "You have new messages.",
    Dismissible: true,
})`),
			),
		),
	)

	// Create page content with proper sidebar navigation
	pageContent := core.Page("Feedback Components",
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
		"Feedback Components - UI Gallery",
		pageContent,
	)

	w.Header().Set("Content-Type", "text/html")
	response.Render(w)
}