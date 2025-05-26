package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"strconv"
)

// ProgressBarProps defines the properties for the ProgressBar component
type ProgressBarProps struct {
	Value int    // current progress value (0-Max)
	Max   int    // maximum progress value
	Label string // optional label text
	Size  string // small, medium, large
	Color string // primary, success, warning, error
}

// ProgressBar creates a visual progress indicator component
func ProgressBar(props ProgressBarProps) g.Node {
	// Set default max if not specified
	max := props.Max
	if max <= 0 {
		max = 100
	}

	// Set default size if not specified
	size := props.Size
	if size == "" {
		size = "medium"
	}

	// Set default color if not specified
	color := props.Color
	if color == "" {
		color = "primary"
	}

	// Ensure value is within bounds
	value := props.Value
	if value < 0 {
		value = 0
	}
	if value > max {
		value = max
	}

	// Calculate percentage
	percentage := float64(value) / float64(max) * 100

	// Build container CSS classes
	var containerClasses []string
	containerClasses = append(containerClasses, "w-full")

	// Add size-specific height
	switch size {
	case "small":
		containerClasses = append(containerClasses, "h-1")
	case "large":
		containerClasses = append(containerClasses, "h-4")
	default: // "medium"
		containerClasses = append(containerClasses, "h-2")
	}

	containerClasses = append(containerClasses, 
		"bg-gray-200", "rounded-full", "overflow-hidden")

	// Build progress bar CSS classes
	var barClasses []string
	barClasses = append(barClasses, 
		"h-full", "transition-all", "duration-300", "ease-out")

	// Add color-specific styles
	switch color {
	case "success":
		barClasses = append(barClasses, "bg-green-500")
	case "warning":
		barClasses = append(barClasses, "bg-yellow-500")
	case "error":
		barClasses = append(barClasses, "bg-red-500")
	default: // "primary"
		barClasses = append(barClasses, "bg-blue-500")
	}

	// Create the progress bar element
	progressBar := html.Div(
		CSSClass(containerClasses...),
		html.Role("progressbar"),
		html.Aria("valuenow", strconv.Itoa(value)),
		html.Aria("valuemin", "0"),
		html.Aria("valuemax", strconv.Itoa(max)),
		html.Div(
			CSSClass(barClasses...),
			html.Style("width: "+strconv.FormatFloat(percentage, 'f', 1, 64)+"%"),
		),
	)

	// If no label is provided, return just the progress bar
	if props.Label == "" {
		return progressBar
	}

	// Create labeled progress bar with text above
	percentageText := strconv.FormatFloat(percentage, 'f', 0, 64) + "%"
	
	return html.Div(
		CSSClass("w-full"),
		// Label and percentage row
		html.Div(
			CSSClass("flex", "justify-between", "items-center", "mb-2"),
			html.Span(
				CSSClass("text-sm", "font-medium", "text-gray-700"),
				g.Text(props.Label),
			),
			html.Span(
				CSSClass("text-sm", "text-gray-500"),
				g.Text(percentageText),
			),
		),
		// Progress bar
		progressBar,
	)
}