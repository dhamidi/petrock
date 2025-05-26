package ui

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// SectionProps defines the properties for the Section component
type SectionProps struct {
	Heading string // Section heading text
	Level   int    // Heading level (1-6), defaults to 2
}

// Section creates a semantic section component with proper heading hierarchy
func Section(props SectionProps, children ...g.Node) g.Node {
	// Set default heading level if not specified or invalid
	level := props.Level
	if level < 1 || level > 6 {
		level = 2
	}

	// Build the section content
	var sectionContent []g.Node

	// Add heading if provided
	if props.Heading != "" {
		var heading g.Node
		headingClasses := CSSClass("text-gray-900", "font-semibold", "mb-4")

		switch level {
		case 1:
			heading = html.H1(headingClasses, CSSClass("text-3xl"), g.Text(props.Heading))
		case 2:
			heading = html.H2(headingClasses, CSSClass("text-2xl"), g.Text(props.Heading))
		case 3:
			heading = html.H3(headingClasses, CSSClass("text-xl"), g.Text(props.Heading))
		case 4:
			heading = html.H4(headingClasses, CSSClass("text-lg"), g.Text(props.Heading))
		case 5:
			heading = html.H5(headingClasses, CSSClass("text-base"), g.Text(props.Heading))
		case 6:
			heading = html.H6(headingClasses, CSSClass("text-sm"), g.Text(props.Heading))
		}

		sectionContent = append(sectionContent, heading)
	}

	// Add children content
	if len(children) > 0 {
		contentDiv := html.Div(
			CSSClass("space-y-4"),
			g.Group(children),
		)
		sectionContent = append(sectionContent, contentDiv)
	}

	return html.Section(
		CSSClass("mb-8"),
		g.Group(sectionContent),
	)
}