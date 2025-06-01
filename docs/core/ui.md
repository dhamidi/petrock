# core/ui/

## Overview

The `core/ui/` directory contains a comprehensive UI component library built with Gomponents and TailwindCSS. It provides reusable, styled components for building consistent web interfaces in Petrock applications.

## Architecture

### Component Structure

Each UI component follows a consistent pattern:
- **Props struct**: Defines the component's properties and configuration options
- **Component function**: Returns a `g.Node` (Gomponent) representing the rendered HTML
- **Styling**: Uses TailwindCSS classes via the `ui.CSSClass()` helper function
- **Variants**: Many components support multiple visual variants (primary, secondary, etc.)

### Key Files

#### ui.go - Core Utilities

Contains base interfaces, constants, and utilities used throughout the UI system:

```go
// Props is the base interface for all component properties
type Props interface{}

// Common CSS class constants
const (
    ClassContainer    = "container mx-auto px-4"
    ClassCard         = "bg-white rounded-lg shadow-md"
    ClassButton       = "px-4 py-2 rounded font-medium transition-colors"
    // ... more constants
)

// Spacing, color, and size variants
const (
    VariantPrimary   = "primary"
    VariantSecondary = "secondary"
    // ...
)
```

#### base.go - Base Components

Provides fundamental building blocks like `CSSClass()` helper function for consistent class application.

#### layout.go - Layout Components

Contains page layout components including headers, footers, and the main `Layout()` function for consistent page structure.

## Component Categories

### Navigation Components
- **navbar.go**: Top navigation bars with branding and menu items
- **breadcrumbs.go**: Hierarchical navigation breadcrumbs
- **sidenav.go**: Sidebar navigation for admin interfaces
- **tabs.go**: Tabbed navigation within pages
- **pagination.go**: Page navigation for data lists

### Form Components
- **text_input.go**: Text input fields with validation states
- **textarea.go**: Multi-line text input areas
- **select.go**: Dropdown selection components
- **radio.go**: Radio button groups
- **checkbox.go**: Checkbox inputs
- **button.go**: Action buttons with multiple variants
- **button_group.go**: Grouped button collections
- **form_group.go**: Form field grouping with labels and help text
- **fieldset.go**: Form section grouping

### Layout & Structure
- **container.go**: Content containers with responsive widths
- **grid.go**: CSS Grid layout system
- **section.go**: Content sections with headers
- **card.go**: Content cards with headers, bodies, and actions
- **divider.go**: Content separation elements

### Feedback Components
- **alert.go**: Status messages and notifications
- **toast.go**: Temporary notification overlays
- **badge.go**: Small status indicators
- **progress_bar.go**: Progress indication
- **loading_spinner.go**: Loading state indicators

### Interactive Components
- **accordion.go**: Collapsible content sections
- **disclosure_widget.go**: Show/hide content toggles
- **toggle.go**: Boolean switch controls

## Gallery System

### gallery/ Directory

The `gallery/` subdirectory contains a complete component documentation and demonstration system:

#### gallery.go - Main Gallery

- Provides the main gallery page at `/_/ui`
- Lists all available components with descriptions
- Includes navigation sidebar for easy component browsing

#### component.go - Individual Component Pages

- Handles routes for individual component demonstrations (e.g., `/_/ui/button`)
- Shows live examples of each component with different variants
- Includes usage code examples and property documentation

#### Gallery Categories

The gallery organizes components into logical categories:

- **Form Controls**: Input, textarea, select, radio, checkbox, button
- **Form Layout**: Form groups, fieldsets, button groups
- **Form Inputs**: Specialized form input demonstrations
- **Navigation**: Navbar, breadcrumbs, tabs, pagination
- **Feedback**: Alerts, toasts, badges, progress bars
- **Layout**: Containers, grids, sections, cards
- **Interactive**: Accordions, disclosure widgets, toggles

## Usage Patterns

### Basic Component Usage

```go
import (
    "github.com/petrock/example_module_path/core/ui"
    g "maragu.dev/gomponents"
    "maragu.dev/gomponents/html"
)

// Using a button component
button := ui.Button(ui.ButtonProps{
    Text:    "Click Me",
    Variant: ui.VariantPrimary,
    Size:    ui.SizeMedium,
})

// Using a container
content := ui.Container(ui.ContainerProps{
    Variant: "default",
}, 
    html.H1(ui.CSSClass("text-2xl", "font-bold"), g.Text("Page Title")),
    html.P(g.Text("Page content goes here")),
)
```

### Form Building

```go
form := ui.Form(ui.FormProps{},
    ui.FormGroup(ui.FormGroupProps{
        Label: "Email Address",
    },
        ui.TextInput(ui.TextInputProps{
            Type:        "email",
            Placeholder: "Enter your email",
            Required:    true,
        }),
    ),
    ui.ButtonGroup(ui.ButtonGroupProps{},
        ui.Button(ui.ButtonProps{
            Text:    "Submit",
            Variant: ui.VariantPrimary,
            Type:    "submit",
        }),
        ui.Button(ui.ButtonProps{
            Text:    "Cancel",
            Variant: ui.VariantSecondary,
        }),
    ),
)
```

### Layout Composition

```go
page := ui.Layout(ui.LayoutProps{
    Title: "My Page",
},
    ui.Container(ui.ContainerProps{},
        ui.Section(ui.SectionProps{
            Title: "Main Content",
        },
            ui.Card(ui.CardProps{
                Title: "Card Title",
            },
                html.P(g.Text("Card content")),
            ),
        ),
    ),
)
```

## Styling Guidelines

### TailwindCSS Integration

- All components use TailwindCSS utility classes
- CSS classes are applied via the `ui.CSSClass()` helper function
- No custom CSS files - all styling is embedded in components
- Responsive design using Tailwind's responsive prefixes

### Color Variants

Components support consistent color variants:
- **Primary**: Main brand colors (typically blue)
- **Secondary**: Secondary brand colors (typically gray)
- **Success**: Success states (typically green)
- **Warning**: Warning states (typically yellow)
- **Danger**: Error/danger states (typically red)
- **Info**: Informational states (typically blue/cyan)

### Size Variants

Many components support size variants:
- **Small (sm)**: Compact sizing for dense interfaces
- **Medium (md)**: Default sizing for most use cases
- **Large (lg)**: Prominent sizing for emphasis

## Design Principles

- **Composability**: Components can be easily combined to build complex interfaces
- **Consistency**: All components follow the same patterns and conventions
- **Accessibility**: Components include appropriate ARIA attributes and semantic HTML
- **Responsiveness**: Components are designed to work across different screen sizes
- **Type Safety**: Strong typing through Go structs for component properties
- **Documentation**: Comprehensive gallery with live examples and code samples
