package petrock_example_feature_name

import (
	"fmt" // Import fmt package
	// "strings" // Removed unused import

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
	g "maragu.dev/gomponents"                     // Alias for gomponents
	. "maragu.dev/gomponents/components"
	"maragu.dev/gomponents/html" // HTML specific components
)

// ItemView renders the HTML representation of a single item.
// Adapt the fields and structure based on the 'Result' type in messages.go.
func ItemView(item Result) g.Node {
	return html.Div(
		// Header with title and navigation
		html.Div(
			Classes{"flex": true, "justify-between": true, "items-center": true, "mb-4": true},
			html.H3(Classes{"text-xl": true, "font-bold": true}, g.Text(item.Name)),
			html.A(g.Attr("href", "/petrock_example_feature_name"), Classes{"text-blue-500": true}, g.Text("Back to List")),
		),
		
		// Item details with labels and values
		html.Div(Classes{"mb-6": true, "bg-gray-50": true, "p-4": true, "rounded": true},
			// ID field
			html.Div(Classes{"mb-2": true},
				html.Label(Classes{"font-semibold": true, "block": true}, g.Text("ID:")),
				html.Div(Classes{"pl-2": true}, g.Text(item.ID)),
			),
			// Description field
			html.Div(Classes{"mb-2": true},
				html.Label(Classes{"font-semibold": true, "block": true}, g.Text("Description:")),
				html.Div(Classes{"pl-2": true}, g.Text(item.Description)),
			),
			// Metadata
			html.Div(Classes{"mt-4": true, "text-sm": true, "text-gray-600": true},
				html.Div(g.Textf("Created: %s", item.CreatedAt.Format("2006-01-02 15:04"))),
				html.Div(g.Textf("Last Updated: %s", item.UpdatedAt.Format("2006-01-02 15:04"))),
				html.Div(g.Textf("Version: %d", item.Version)),
			),
		),
		
		// Action buttons
		html.Div(Classes{"mt-4": true, "flex": true, "space-x-2": true},
			// Edit link
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/edit"),
				Classes{"px-4": true, "py-2": true, "bg-blue-500": true, "text-white": true, "rounded": true, "hover:bg-blue-600": true},
				g.Text("Edit"),
			),
			// Delete link
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/delete"),
				Classes{"px-4": true, "py-2": true, "bg-red-500": true, "text-white": true, "rounded": true, "hover:bg-red-600": true},
				g.Text("Delete"),
			),
		),
	)
}

// ItemForm renders an HTML <form> for creating or editing an item.
// It uses core.Form for data and error handling.
// 'item' can be nil when creating a new item.
// 'csrfToken' should be provided by the handler.
func ItemForm(form *core.Form, item *Result, csrfToken string) g.Node {
	// Determine if we're creating or editing
	isEdit := item != nil
	var title, submitLabel string
	var actionURL string
	
	if isEdit {
		title = "Edit Item"
		submitLabel = "Update Item"
		actionURL = fmt.Sprintf("/petrock_example_feature_name/%s/edit", item.ID) 
	} else {
		title = "Create New Item"
		submitLabel = "Create Item"
		actionURL = "/petrock_example_feature_name/new"
	}

	// Get values from form (if validation failed) or from item (if editing)
	nameValue := form.Get("name")
	descriptionValue := form.Get("description")
	if !form.HasError("name") && isEdit {
		nameValue = item.Name
	}
	if !form.HasError("description") && isEdit {
		descriptionValue = item.Description
	}

	// Form error helpers
	formErrorMessage := func(field string) g.Node {
		if !form.HasError(field) {
			return nil
		}
		return html.Div(
			Classes{"text-red-500": true, "text-sm": true, "mt-1": true},
			g.Text(form.GetError(field)),
		)
	}

	return html.Div(
		// Form header with title and back button
		html.Div(
			Classes{"flex": true, "justify-between": true, "items-center": true, "mb-6": true},
			html.H2(Classes{"text-xl": true, "font-semibold": true}, g.Text(title)),
			html.A(
				g.Attr("href", backLink(isEdit, item)),
				Classes{"text-blue-500": true},
				g.Text("Back"),
			),
		),
		
		html.Form(
			// Form attributes
			html.Action(actionURL),
			html.Method("POST"),

			// CSRF Token (as hidden input)
			html.Input(
				g.Attr("type", "hidden"),
				g.Attr("name", "csrf_token"),
				g.Attr("value", csrfToken),
			),

			// Form Fields
			html.Div(
				Classes{"mb-4": true},
				html.Label(
					html.For("name"),
					Classes{"block": true, "text-sm": true, "font-medium": true, "text-gray-700": true, "mb-1": true},
					g.Text("Name"),
				),
				html.Input(
					g.Attr("type", "text"),
					g.Attr("name", "name"),
					g.Attr("id", "name"),
					g.Attr("value", nameValue),
					Classes{
						"w-full": true, "p-2": true, "border": true, "rounded": true,
						"border-red-500": form.HasError("name"),
						"border-gray-300": !form.HasError("name"),
					},
				),
				formErrorMessage("name"),
			),
			html.Div(
				Classes{"mb-4": true},
				html.Label(
					html.For("description"),
					Classes{"block": true, "text-sm": true, "font-medium": true, "text-gray-700": true, "mb-1": true},
					g.Text("Description"),
				),
				html.TextArea(
					g.Attr("name", "description"),
					g.Attr("id", "description"),
					Classes{
						"w-full": true, "p-2": true, "border": true, "rounded": true,
						"border-red-500": form.HasError("description"),
						"border-gray-300": !form.HasError("description"),
					},
					g.Text(descriptionValue),
				),
				formErrorMessage("description"),
			),

			// Submit Button
			html.Div(
				Classes{"mt-6": true},
				html.Button(
					g.Attr("type", "submit"),
					Classes{"px-4": true, "py-2": true, "bg-blue-500": true, "text-white": true, "rounded": true, "hover:bg-blue-600": true},
					g.Text(submitLabel),
				),
			),
		),
	)
}

// ItemsListView renders a list of items, potentially with pagination.
func ItemsListView(result ListResult) g.Node {
	// Header with title and New button
	header := html.Div(
		Classes{"flex": true, "justify-between": true, "items-center": true, "mb-6": true},
		html.H2(Classes{"text-xl": true, "font-semibold": true}, g.Text("Items")),
		html.A(
			g.Attr("href", "/petrock_example_feature_name/new"),
			Classes{"px-4": true, "py-2": true, "bg-green-500": true, "text-white": true, "rounded": true, "hover:bg-green-600": true},
			g.Text("New Item"),
		),
	)

	// Item list
	var itemList g.Node
	if len(result.Items) == 0 {
		itemList = html.Div(
			Classes{"p-6": true, "text-center": true, "bg-gray-50": true, "rounded": true, "text-gray-500": true},
			g.Text("No items found. Create one using the 'New Item' button."),
		)
	} else {
		items := make([]g.Node, 0, len(result.Items))
		for _, item := range result.Items {
			// Create a simplified item row for the list view
			items = append(items, html.Div(
				Classes{"p-4": true, "border": true, "border-gray-200": true, "rounded": true, "mb-2": true, "hover:bg-gray-50": true},
				html.Div(
					Classes{"flex": true, "justify-between": true, "items-center": true},
					html.Div(
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID),
							Classes{"font-semibold": true, "text-blue-600": true, "hover:underline": true},
							g.Text(item.Name),
						),
						html.Div(Classes{"text-sm": true, "text-gray-600": true, "mt-1": true}, g.Text(item.Description)),
					),
					html.Div(
						Classes{"flex": true, "space-x-2": true},
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/edit"),
							Classes{"text-blue-500": true, "hover:underline": true},
							g.Text("Edit"),
						),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/delete"),
							Classes{"text-red-500": true, "hover:underline": true},
							g.Text("Delete"),
						),
					),
				),
			))
		}
		itemList = html.Div(g.Group(items))
	}

	// Pagination
	pagination := html.Div(
		Classes{"mt-6": true, "flex": true, "justify-between": true, "items-center": true},
		html.Span(Classes{"text-sm": true, "text-gray-600": true}, g.Textf("Total: %d", result.TotalCount)),
		html.Div(
			Classes{"flex": true, "space-x-2": true},
			html.Span(Classes{"text-sm": true, "text-gray-600": true}, g.Textf("Page %d of %d", result.Page, (result.TotalCount+result.PageSize-1)/result.PageSize)),
		),
	)

	return html.Div(header, itemList, pagination)
}

// DeleteConfirmForm renders a form to confirm deletion of an item.
func DeleteConfirmForm(item *Result, csrfToken string) g.Node {
	return html.Div(
		// Header with title and back button
		html.Div(
			Classes{"flex": true, "justify-between": true, "items-center": true, "mb-6": true},
			html.H2(Classes{"text-xl": true, "font-semibold": true}, g.Text("Confirm Delete")),
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID),
				Classes{"text-blue-500": true},
				g.Text("Back"),
			),
		),
		
		// Item details
		html.Div(
			Classes{"mb-6": true, "p-4": true, "bg-red-50": true, "border": true, "border-red-200": true, "rounded": true},
			html.P(
				Classes{"mb-2": true},
				g.Text("Are you sure you want to delete this item?"),
			),
			html.Div(
				Classes{"font-semibold": true},
				g.Text(item.Name),
			),
			html.Div(
				Classes{"text-sm": true, "text-gray-600": true, "mt-1": true},
				g.Text(item.Description),
			),
		),
		
		// Confirmation form
		html.Form(
			html.Action("/petrock_example_feature_name/"+item.ID+"/delete"),
			html.Method("POST"),
			
			// CSRF Token
			html.Input(
				g.Attr("type", "hidden"),
				g.Attr("name", "csrf_token"),
				g.Attr("value", csrfToken),
			),
			
			// Action buttons
			html.Div(
				Classes{"flex": true, "space-x-4": true},
				html.Button(
					g.Attr("type", "submit"),
					Classes{"px-4": true, "py-2": true, "bg-red-500": true, "text-white": true, "rounded": true, "hover:bg-red-600": true},
					g.Text("Delete Item"),
				),
				html.A(
					g.Attr("href", "/petrock_example_feature_name/"+item.ID),
					Classes{"px-4": true, "py-2": true, "bg-gray-300": true, "text-gray-800": true, "rounded": true, "hover:bg-gray-400": true},
					g.Text("Cancel"),
				),
			),
		),
	)
}

// NewItemButton renders a button or link to navigate to the item creation page/view.
func NewItemButton() g.Node {
	return html.A(
		g.Attr("href", "/petrock_example_feature_name/new"),
		Classes{"px-4": true, "py-2": true, "bg-green-500": true, "text-white": true, "rounded": true, "hover:bg-green-600": true},
		g.Text("New Item"),
	)
}

// backLink returns the appropriate back link URL based on whether we're in edit mode
func backLink(isEdit bool, item *Result) string {
	if isEdit && item != nil {
		return "/petrock_example_feature_name/" + item.ID
	}
	return "/petrock_example_feature_name"
}
