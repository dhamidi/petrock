package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// GridProps defines the properties for the Grid component
type GridProps struct {
	Columns string // CSS grid-template-columns value (e.g., "1fr 1fr", "repeat(3, 1fr)", "200px 1fr")
	Gap     string // CSS gap value (e.g., "1rem", "16px", "2rem 1rem")
	Areas   string // CSS grid-template-areas value for named grid areas
}

// Grid creates a CSS Grid container component
func Grid(props GridProps, children ...g.Node) g.Node {
	// Build styles map
	styles := map[string]string{
		"display": "grid",
	}

	// Set grid-template-columns
	if props.Columns != "" {
		styles["grid-template-columns"] = props.Columns
	} else {
		styles["grid-template-columns"] = "1fr" // Default to single column
	}

	// Set gap
	if props.Gap != "" {
		styles["gap"] = props.Gap
	} else {
		styles["gap"] = "1rem" // Default gap
	}

	// Set grid-template-areas if provided
	if props.Areas != "" {
		styles["grid-template-areas"] = props.Areas
	}

	var attributes []g.Node
	
	// Add base grid class
	attributes = append(attributes, CSSClass("grid"))
	
	// Add styles
	attributes = append(attributes, Style(styles))

	// Combine attributes and children
	var allNodes []g.Node
	allNodes = append(allNodes, attributes...)
	allNodes = append(allNodes, children...)

	return html.Div(allNodes...)
}

// GridItem creates a grid item with optional area name
func GridItem(area string, children ...g.Node) g.Node {
	var styleStr string
	if area != "" {
		styleStr = "grid-area:" + area + ";"
	}

	return html.Div(
		html.Style(styleStr),
		g.Group(children),
	)
}