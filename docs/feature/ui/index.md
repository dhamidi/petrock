# UI

The `ui` directory contains UI components and templates for rendering the web interface of the feature.

## Structure

- `components/` - Reusable UI components
  - `forms.go` - Form input components
  - `tables.go` - Table and list display components
- `layouts/` - Page layouts
  - `main.go` - Standard page layout
  - `modal.go` - Modal dialog layouts
- `pages/` - Complete page views
  - `list.go` - List view for multiple items
  - `detail.go` - Detail view for single items
  - `forms.go` - Form views for creating and editing
  - `delete.go` - Delete confirmation view
- `helpers.go` - View helper functions and utilities

## Responsibilities

- Define the visual presentation of the feature
- Implement interactive UI components
- Render data from queries into HTML
- Handle user input through forms
- Maintain consistent UI patterns

## Component Pattern

UI components follow a compositional pattern where:

1. Small, reusable components are defined in the components directory
2. These components are composed into larger views in the pages directory
3. Layouts provide the overall structure for complete pages

## Usage

UI components are used by handlers to render HTTP responses, typically by passing query results to page templates.