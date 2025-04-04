package petrock_example_feature_name

import (
	"fmt" // Import fmt package
	"strings" // Import strings package

	g "maragu.dev/gomponents"                 // Alias for gomponents
	"maragu.dev/gomponents/html"              // HTML specific components
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// ItemView renders the HTML representation of a single item.
// Adapt the fields and structure based on the 'Result' type in messages.go.
func ItemView(item Result) g.Node {
	return html.Div(
		// Example structure - customize as needed
		html.H3(g.Text(item.Name)),
		html.P(g.Textf("ID: %s", item.ID)),
		html.P(g.Textf("Description: %s", item.Description)),
		html.Small(
			g.Textf("Created: %s, Updated: %s, Version: %d",
				item.CreatedAt.Format("2006-01-02 15:04"),
				item.UpdatedAt.Format("2006-01-02 15:04"),
				item.Version,
			),
		),
		// Add Edit/Delete buttons if applicable
		html.Button(
			g.Classes{"ml-2": true, "text-blue-500": true, "hover:underline": true}, // Example styling
			g.Attr("hx-get", fmt.Sprintf("/feature-path/%s/edit", item.ID)),           // Placeholder URL
			g.Attr("hx-target", "#feature-content"),                                   // Placeholder target
			g.Text("Edit"),
		),
		html.Button(
			g.Classes{"ml-2": true, "text-red-500": true, "hover:underline": true}, // Example styling
			g.Attr("hx-delete", fmt.Sprintf("/feature-path/%s", item.ID)),             // Placeholder URL
			g.Attr("hx-target", "#feature-content"),                                   // Placeholder target
			g.Attr("hx-confirm", "Are you sure you want to delete this item?"),
			g.Text("Delete"),
		),
	)
}

// ItemForm renders an HTML <form> for creating or editing an item.
// It uses core.Form for data and error handling.
// 'item' can be nil when creating a new item.
// 'csrfToken' should be provided by the handler.
func ItemForm(form *core.Form, item *Result, csrfToken string) g.Node {
	actionURL := "/feature-path" // Placeholder URL for creating
	method := "POST"
	if item != nil {
		actionURL = fmt.Sprintf("/feature-path/%s", item.ID) // Placeholder URL for updating
		method = "PUT"                                       // Or PATCH
	}

	// Get values from form (if validation failed) or from item (if editing)
	nameValue := form.Get("name")
	descriptionValue := form.Get("description")
	if !form.HasError("name") && item != nil {
		nameValue = item.Name
	}
	if !form.HasError("description") && item != nil {
		descriptionValue = item.Description
	}

	return html.Form(
		// HTMX attributes for form submission
		g.Attr("hx-"+strings.ToLower(method), actionURL),
		g.Attr("hx-target", "#feature-content"), // Placeholder target
		g.Attr("hx-swap", "outerHTML"),          // Example swap strategy

		// CSRF Token
		core.CSRFTokenInput(csrfToken), // Assumes core.CSRFTokenInput exists

		// Form Fields using core components
		html.Div(
			g.Classes{"mb-4": true},
			html.Label(html.For("name"), g.Classes{"block": true, "text-sm": true, "font-medium": true, "text-gray-700": true}, g.Text("Name")),
			core.Input("text", "name", nameValue, html.ID("name"), g.Classes{"border-red-500": form.HasError("name")}), // Add error class conditionally
			core.FormError(form, "name"),
		),
		html.Div(
			g.Classes{"mb-4": true},
			html.Label(html.For("description"), g.Classes{"block": true, "text-sm": true, "font-medium": true, "text-gray-700": true}, g.Text("Description")),
			core.TextArea("description", descriptionValue, html.ID("description"), g.Classes{"border-red-500": form.HasError("description")}), // Add error class conditionally
			core.FormError(form, "description"),
		),

		// Submit Button
		core.Button("Save Item"), // Uses core.Button component
	)
}

// ItemsListView renders a list of items, potentially with pagination.
func ItemsListView(result ListResult) g.Node {
	itemNodes := make([]g.Node, 0, len(result.Items))
	if len(result.Items) == 0 {
		itemNodes = append(itemNodes, html.P(g.Text("No items found.")))
	} else {
		for _, item := range result.Items {
			// Render each item - could be a full ItemView or a summary row
			itemNodes = append(itemNodes, html.Li(ItemView(item))) // Example: using ItemView within a list
		}
	}

	return html.Div(
		html.H2(g.Classes{"text-xl": true, "font-semibold": true, "mb-4": true}, g.Text("Items")),
		NewItemButton(), // Add button to create new item
		html.Ul(
			g.Classes{"space-y-4": true}, // Example list styling
			g.Group(itemNodes),
		),
		// Basic Pagination Example (implement proper logic based on result)
		html.Div(
			g.Classes{"mt-4": true, "flex": true, "justify-between": true, "items-center": true},
			html.Span(g.Textf("Total: %d", result.TotalCount)),
			html.Span(g.Textf("Page %d/%d", result.Page, (result.TotalCount+result.PageSize-1)/result.PageSize)), // Calculate total pages
			// Add Previous/Next buttons with HTMX if needed
		),
	)
}

// NewItemButton renders a button or link to trigger loading the ItemForm for creation.
func NewItemButton() g.Node {
	return core.Button(
		"Create New Item",
		g.Attr("hx-get", "/feature-path/new"), // Placeholder URL for the form
		g.Attr("hx-target", "#feature-content"), // Placeholder target where the form should load
		// Add other attributes as needed
	)
}
