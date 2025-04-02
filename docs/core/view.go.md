# Plan for core/view.go

This file provides shared, reusable Gomponents components for building the UI.

## Types

- None specific to this file. Component logic is encapsulated in functions.

## Functions

*These functions return `gomponents.Node` and typically accept attributes (`gomponents.Attr`) for customization.*

- `Button(text string, attrs ...gomponents.Attribute) gomponents.Node`: Renders a `<button>` element with standard styling. Can accept attributes for type, event handlers (hx-*, etc.).
- `Input(inputType, name, value string, attrs ...gomponents.Attribute) gomponents.Node`: Renders an `<input>` element. Takes type (text, password, email, hidden), name, and current value. Can accept attributes for placeholders, validation attributes, etc.
- `TextArea(name, value string, attrs ...gomponents.Attribute) gomponents.Node`: Renders a `<textarea>` element.
- `Select(name string, options map[string]string, selectedValue string, attrs ...gomponents.Attribute) gomponents.Node`: Renders a `<select>` element with `<option>`s.
- `FormError(form *Form, field string) gomponents.Node`: Renders error messages for a specific field from a `core.Form` instance. Typically renders a `<span>` or `<div>` with an error class if `form.HasError(field)` is true.
- `CSRFTokenInput(token string) gomponents.Node`: Renders a hidden input field for CSRF token protection. `g.Input(g.Type("hidden"), g.Name("csrf_token"), g.Value(token))`
