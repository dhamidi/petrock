# Deficiencies

## feature_template/worker.go

- Ensure workers can properly track their position in the event log
- Define clear patterns for error handling and retry strategies
- Consider standardizing common worker operations (like making HTTP calls to external services)

## Design system

* Part of the core 
* probably start with DaisyUI

## Tools and MCP support 

* need to be added to introspection
* subcommand `serve` should accept protocol (http, MCP)

## Command/Query generators

## Asset pipeline

* plain JS with importmap support
* compression + hashing


## Rules

- working in a petrock generated project is still rough, as the AI is lacking context about important rules (e.g. no direct mutation of the state)
- we should automatically generate rules 

## Design System Implementation

* Define and implement a comprehensive design system for visual coherence across petrock generated applications
* Extend DaisyUI integration with custom components that follow petrock's design principles
* Include theming support with sensible defaults

### Implementation Approach

* Create a `core/design` package with component library built on Gomponents
* Define a theme configuration system in `core/design/theme.go`
* Implement core UI patterns (cards, forms, tables, navigation) as reusable components
* Add design tokens for colors, spacing, typography in `core/design/tokens.go`

#### Core Components

* **Layout Components**
  * Container (with variants: default, narrow, wide, full)
  * Grid (flexible CSS Grid wrapper)
  * Card (with header, body, footer sections)
  * Section (semantic section with optional heading)
  * Divider (horizontal rule with styling options)

* **Navigation Components**
  * NavBar (responsive top navigation)
  * SideNav (collapsible sidebar navigation)
  * Tabs (accessible tab interface with progressive enhancement)
  * Breadcrumbs (navigation path indicator)
  * Pagination (for multiple page navigation)

* **Form Components**
  * TextInput (with validation states)
  * TextArea (multi-line input)
  * Select (dropdown select)
  * Checkbox & Radio (with proper styling)
  * Toggle (accessible switch component)
  * FormGroup (label + input wrapper with validation messages)
  * FieldSet (grouping related form elements)

* **Interactive Elements**
  * Button (with variants: primary, secondary, danger, link)
  * ButtonGroup (horizontally grouped buttons)
  * Accordion (expandable sections without JS, using CSS)
  * DisclosureWidget (show/hide content with CSS :target or :checked hacks)

* **Feedback Components**
  * Alert (contextual feedback messages)
  * Badge (status indicators)
  * Toast (could work with CSS animations for entry/exit)
  * ProgressBar (visual completion indicator)
  * LoadingSpinner (CSS-only animation)

#### Design Tokens

* **Colors**
  * Brand colors (primary, secondary, tertiary)
  * UI colors (background, surface, border)
  * Semantic colors (success, warning, error, info)
  * Neutral palette (10 shades from white to black)
  * Accessibility considerations (contrast ratios documented)

* **Typography**
  * Font families (heading, body, monospace)
  * Font sizes (scale from xs to 3xl)
  * Line heights (tight, normal, loose)
  * Font weights (light, regular, medium, bold)
  * Text styles (heading1-6, body, caption, code)

* **Spacing**
  * Spacing scale (4px, 8px, 16px, 24px, 32px, 48px, 64px, 96px, 128px)
  * Layout spacing (page margins, gutters, section spacing)
  * Component internal spacing (padding inside cards, buttons, etc.)

* **Borders & Shadows**
  * Border widths (thin, regular, thick)
  * Border radii (none, sm, md, lg, pill, circle)
  * Shadows (none, sm, md, lg, inner)
  * Focus states (accessible focus rings)

## Subpackage Structure

* Introduce subpackages to support many different files of each kind (every command in a separate file, every command handler in a separate file, etc)
* Define conventions for file organization within features
* Provide utilities for discovery and auto-registration of components

### Implementation Approach

* Restructure feature package layout:
  ```
  posts/
    commands/
      create.go    # Contains CreatePostCommand struct and validation
      update.go
    handlers/
      create.go    # Contains handler for CreatePostCommand
      update.go
    queries/
      list.go      # Contains ListPostsQuery and handler
      get.go
  ```
* Add reflection-based auto-registration in `feature/register.go`
* Implement code generation tooling to scaffold new commands/queries in this structure

## Component Generators

* Add generators for each individual petrock component (commands, queries, workers, views, etc.)
* Implement self-inspection mechanism to determine what pieces need to be added
* Create template unpacking system to selectively generate only required components

### Implementation Approach

* Add granular CLI commands:
  ```
  petrock generate command posts create
  petrock generate query posts list
  petrock generate worker posts emailNotifier
  petrock generate view posts show
  ```
* Leverage existing `core.GetInspectResult()` (from `internal/skeleton/core/inspect.go`) to analyze project structure
* Enhance inspectors to extract result types for queries and validation rules for commands
* Build component dependency resolution system that works with valid Go template files:
  ```go
  // Example pseudocode
  func generateCommand(feature, name string) {
    inspect := app.GetInspectResult()
    // Check if command already exists
    for _, cmd := range inspect.Commands {
      if cmd.Name == feature+"/"+name {
        return // Already exists
      }
    }
    // Check if feature exists, create if not
    if !contains(inspect.Features, feature) {
      generateFeature(feature)
    }
    // Copy and transform the template from internal/skeleton/feature_template/commands.go
    // using the established placeholder replacement pattern
    destPath := filepath.Join(feature, "commands", name+".go")
    templateContent := readEmbeddedFile("internal/skeleton/feature_template/commands.go")
    
    // Replace placeholders with actual values
    result := strings.ReplaceAll(templateContent, "petrock_example_feature_name", feature)
    result = strings.ReplaceAll(result, "ExampleCommand", strings.Title(name)+"Command")
    
    writeFile(destPath, result)
  }
  ```
