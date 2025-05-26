# Design System Implementation

- Define and implement a comprehensive design system for visual coherence across petrock generated applications
- Include theming support with sensible defaults

## Implementation Approach

- Create a `core/ui` package with component library built on Gomponents (this goes into `internal/skeleton/core`)
  - replaces `internal/skeleton/core/view.go` and `internal/skeleton/core/layout.go`
- Define a theme configuration system in `core/ui/theme.go`
- Implement core UI patterns (cards, forms, tables, navigation) as reusable components
- Add design tokens for colors, spacing, typography in `core/ui/tokens.go`

## Code structure

The package lives in `internal/skeleton/core/ui`.

Each component lives in a separate `.go` file.

A component file is roughly structured like this:

```go
package ui

import (
  g "maragu.dev/gomponents"
  . "maragu.dev/gomponents/html"
)

type ContainerProps struct {
  Variant string // "default", "narrow", "wide", "full"
}

func Container(props ContainerProps, children ...g.Node) g.Node {
  return Div()
}
```

## Core Components

- **Layout Components**

  - Container (with variants: default, narrow, wide, full)
  - Grid (flexible CSS Grid wrapper)
  - Card (with header, body, footer sections)
  - Section (semantic section with optional heading)
  - Divider (horizontal rule with styling options)

- **Navigation Components**

  - NavBar (responsive top navigation)
  - SideNav (collapsible sidebar navigation)
  - Tabs (accessible tab interface with progressive enhancement)
  - Breadcrumbs (navigation path indicator)
  - Pagination (for multiple page navigation)

- **Form Components**

  - TextInput (with validation states)
  - TextArea (multi-line input)
  - Select (dropdown select)
  - Checkbox & Radio (with proper styling)
  - Toggle (accessible switch component)
  - FormGroup (label + input wrapper with validation messages)
  - FieldSet (grouping related form elements)

- **Interactive Elements**

  - Button (with variants: primary, secondary, danger, link)
  - ButtonGroup (horizontally grouped buttons)
  - Accordion (expandable sections without JS, using CSS)
  - DisclosureWidget (show/hide content with CSS :target or :checked hacks)

- **Feedback Components**
  - Alert (contextual feedback messages)
  - Badge (status indicators)
  - Toast (could work with CSS animations for entry/exit)
  - ProgressBar (visual completion indicator)
  - LoadingSpinner (CSS-only animation)

## Design Tokens

- **Colors**

  - Brand colors (primary, secondary, tertiary)
  - UI colors (background, surface, border)
  - Semantic colors (success, warning, error, info)
  - Neutral palette (10 shades from white to black)
  - Accessibility considerations (contrast ratios documented)

- **Typography**

  - Font families (heading, body, monospace)
  - Font sizes (scale from xs to 3xl)
  - Line heights (tight, normal, loose)
  - Font weights (light, regular, medium, bold)
  - Text styles (heading1-6, body, caption, code)

- **Spacing**

  - Spacing scale (4px, 8px, 16px, 24px, 32px, 48px, 64px, 96px, 128px)
  - Layout spacing (page margins, gutters, section spacing)
  - Component internal spacing (padding inside cards, buttons, etc.)

- **Borders & Shadows**
  - Border widths (thin, regular, thick)
  - Border radii (none, sm, md, lg, pill, circle)
  - Shadows (none, sm, md, lg, inner)
  - Focus states (accessible focus rings)

## Integration with the rest of the system

- a gallery page showcasing all components is exposed at `/_/ui`
- this pages shows a list of all components in a sidebar on the left, and the currently selected component in a panel on the right
- individual component pages can be found at `/_/ui/{component}`
- package `core/ui/gallery` defines pages and request handlers for the UI gallery

## Outline of the implementation approach

We first need to ensure that we have a way to inspect progress:

- add the `core/ui/gallery` package first, defining the gallery index page with the components list
- register the gallery page in `internal/skeleton/petrock_example_project_name/cmd/serve.go`

Then we proceed with adding a single layout component, the container.

- define core/ui/container.go
- add the component detail page in core/ui/gallery/
- register the detail page in serve.go

Then we proceed with adding the remaining components one by one.

Finally, we add design tokens:

- design tokens need to be visible in the sidebar on the `/_/ui` gallery page under an entry "Design Tokens"
- when selecting that entry, all tokens are listed in the main panel to the right
- the URL is going to be `/_/ui/design-tokens`
