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
		// Item information card
		html.Div(
			g.Attr("class", "space-y-6"),
			
			// Item metadata in a grid layout - responsive
			html.Dl(
				g.Attr("class", "grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-6"),
				
				// ID field
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("ID")),
					html.Dd(g.Attr("class", "mt-1 text-sm text-slate-900"), g.Text(item.ID)),
				),
				
				// Name field
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Name")),
					html.Dd(g.Attr("class", "mt-1 text-sm text-slate-900 font-medium"), g.Text(item.Name)),
				),
				
				// Description field - spans full width
				html.Div(
					g.Attr("class", "col-span-1 sm:col-span-2"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Description")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-900 whitespace-pre-wrap"), 
						g.Text(item.Description),
					),
				),
				
				// Created date 
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Created")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-500"),
						g.Text(item.CreatedAt.Format("Jan 2, 2006 at 15:04")),
					),
				),
				
				// Updated date
				html.Div(
					g.Attr("class", "col-span-1"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Last Updated")),
					html.Dd(
						g.Attr("class", "mt-1 text-sm text-slate-500"),
						g.Text(item.UpdatedAt.Format("Jan 2, 2006 at 15:04")),
					),
				),
				
				// Version number
				html.Div(
					g.Attr("class", "col-span-1 sm:col-span-2"),
					html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Version")),
					html.Dd(g.Attr("class", "mt-1 text-sm text-slate-500"), g.Textf("%d", item.Version)),
				),
			),
		),
		
		// Actions section - responsive
		html.Div(
			g.Attr("class", "mt-8 flex flex-col sm:flex-row gap-3"),
			// Edit button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/edit"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:w-auto w-full"),
				g.Text("Edit Item"),
			),
			// Delete button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/delete"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:w-auto w-full"),
				g.Text("Delete Item"),
			),
			// Back to list
			html.A(
				g.Attr("href", "/petrock_example_feature_name"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-slate-300 text-sm font-medium rounded-md shadow-sm text-slate-700 bg-white hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:w-auto w-full"),
				g.Text("Back to List"),
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

	// Pass-through to global form error display helper
	formErrorMessage := func(field string) g.Node {
		return formErrorDisplay(form, field)
	}

	return html.Div(
		// Form container with back link
		g.Attr("class", "space-y-8"),
		
		// Back navigation
		html.Div(
			g.Attr("class", "flex justify-end"),
			html.A(
				g.Attr("href", backLink(isEdit, item)),
				g.Attr("class", "text-sm font-medium text-indigo-600 hover:text-indigo-500"),
				html.Span(g.Attr("aria-hidden", "true"), g.Text("u2190")),
				g.Text(" Back"),
			),
		),
		
		html.Form(
			g.Attr("class", "bg-white shadow-sm border border-slate-200 rounded-lg overflow-hidden"),
			// Form attributes
			html.Action(actionURL),
			html.Method("POST"),

			// Form header
			html.Div(
				g.Attr("class", "border-b border-slate-200 bg-slate-50 px-4 py-5 sm:px-6"),
				html.H3(
					g.Attr("class", "text-lg font-medium leading-6 text-slate-900"),
					g.Text(title),
				),
				html.P(
					g.Attr("class", "mt-1 text-sm text-slate-500"),
					g.Text(isEdit ? "Update the item details below." : "Fill out the form below to create a new item."),
				),
			),

			// Form body
			html.Div(
				g.Attr("class", "px-4 py-5 sm:p-6 space-y-6"),
				
				// CSRF Token
				csrfField(csrfToken),

				// Name field
				html.Div(
					html.Label(
						g.Attr("for", "name"),
						g.Attr("class", "block text-sm font-medium text-slate-700"),
						g.Text("Name"),
					),
					html.Div(
						g.Attr("class", "mt-1"),
						html.Input(
						g.Attr("type", "text"),
						g.Attr("name", "name"),
						g.Attr("id", "name"),
						g.Attr("value", nameValue),
						g.Attr("class", formFieldClass(form, "name")),
					),
						formErrorMessage("name"),
					),
				),

				// Description field
				html.Div(
					html.Label(
						g.Attr("for", "description"),
						g.Attr("class", "block text-sm font-medium text-slate-700"),
						g.Text("Description"),
					),
					html.Div(
						g.Attr("class", "mt-1"),
						html.Textarea(
						g.Attr("name", "description"),
						g.Attr("id", "description"),
						g.Attr("rows", "4"),
						g.Attr("class", formFieldClass(form, "description")),
						g.Text(descriptionValue),
					),
						formErrorMessage("description"),
					),
				),
			),

			// Form footer with submit button
			html.Div(
				g.Attr("class", "px-4 py-3 bg-slate-50 text-right sm:px-6 border-t border-slate-200"),
				html.Button(
					g.Attr("type", "submit"),
					g.Attr("class", "inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
					g.Text(submitLabel),
				),
			),
		),
	)
}

// ItemsListView renders a list of items, potentially with pagination.
func ItemsListView(result ListResult) g.Node {
	// Determine total number of pages
	totalPages := 1
	if result.PageSize > 0 {
		totalPages = (result.TotalCount + result.PageSize - 1) / result.PageSize
	}

	return html.Div(
		g.Attr("class", "space-y-8"),
		
		// List header with actions
		html.Div(
			g.Attr("class", "flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"),
			
			// Stats section
			html.Div(
				g.Attr("class", "text-sm text-slate-500"),
				g.Textf(
					"%d item%s • Page %d of %d", 
					result.TotalCount, 
					pluralize(result.TotalCount), 
					result.Page, 
					totalPages,
				),
			),
			
			// Create new button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/new"),
				g.Attr("class", "inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
				g.Text("New Item"),
			),
		),

		// Item list with item cards
		func() g.Node {
			// Empty state
			if len(result.Items) == 0 {
				return html.Div(
					g.Attr("class", "text-center p-12 border border-dashed border-slate-300 rounded-lg"),
					html.Div(
						g.Attr("class", "mx-auto h-12 w-12 text-slate-400"),
						// Simple placeholder icon
						html.Svg(
							g.Attr("fill", "none"),
							g.Attr("viewBox", "0 0 24 24"),
							g.Attr("stroke", "currentColor"),
							g.Attr("aria-hidden", "true"),
							html.Path(
								g.Attr("stroke-linecap", "round"),
								g.Attr("stroke-linejoin", "round"),
								g.Attr("stroke-width", "2"),
								g.Attr("d", "M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"),
							),
						),
					),
					html.H3(
						g.Attr("class", "mt-2 text-sm font-medium text-slate-900"),
						g.Text("No items"),
					),
					html.P(
						g.Attr("class", "mt-1 text-sm text-slate-500"),
						g.Text("Get started by creating a new item."),
					),
					html.Div(
						g.Attr("class", "mt-6"),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/new"),
							g.Attr("class", "inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
							g.Text("Create New Item"),
						),
					),
				)
			}
			
			// Items list in a responsive grid
			items := make([]g.Node, 0, len(result.Items))
			for _, item := range result.Items {
				items = append(items, html.Div(
					g.Attr("class", "bg-white shadow-sm rounded-lg overflow-hidden border border-slate-200 hover:shadow-md transition-shadow duration-200"),
					
					// Card header with name and created date
					html.Div(
						g.Attr("class", "px-4 py-5 sm:px-6 border-b border-slate-200 bg-slate-50"),
						html.Div(
							g.Attr("class", "flex justify-between items-center flex-wrap gap-2"),
							html.H3(
								g.Attr("class", "text-lg font-medium leading-6 text-slate-900 break-all"),
								html.A(
									g.Attr("href", "/petrock_example_feature_name/"+item.ID),
									g.Attr("class", "hover:text-indigo-600"),
									g.Text(item.Name),
								),
							),
							html.Span(
								g.Attr("class", "inline-flex items-center rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-600"),
								g.Text(item.CreatedAt.Format("Jan 2, 2006")),
							),
						),
					),
					
					// Card body with description
					html.Div(
						g.Attr("class", "px-4 py-5 sm:p-6"),
						html.P(
							g.Attr("class", "text-sm text-slate-700 break-words line-clamp-3"),
							g.Text(item.Description),
						),
					),
					
					// Card footer with actions
					html.Div(
						g.Attr("class", "border-t border-slate-200 bg-slate-50 px-4 py-4 sm:px-6 flex justify-end space-x-3"),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/edit"),
							g.Attr("class", "text-sm font-medium text-indigo-600 hover:text-indigo-500"),
							g.Text("Edit"),
						),
						html.A(
							g.Attr("href", "/petrock_example_feature_name/"+item.ID+"/delete"),
							g.Attr("class", "text-sm font-medium text-red-600 hover:text-red-500"),
							g.Text("Delete"),
						),
					),
				))
			}
			
			return html.Div(
				g.Attr("class", "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"),
				g.Group(items),
			)
		}(),

		// Pagination controls (if more than one page)
		func() g.Node {
			if totalPages <= 1 {
				return nil
			}
			
			// Simple pagination with prev/next buttons
			return html.Nav(
				g.Attr("class", "border-t border-slate-200 px-4 flex items-center justify-between sm:px-0"),
				
				// Previous page button
				html.Div(func() g.Node {
					if result.Page <= 1 {
						return html.Span(
							g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-300 bg-white border border-slate-300 cursor-not-allowed rounded-md"),
							g.Text("Previous"),
						)
					}
					
					prevPage := result.Page - 1
					prevURL := fmt.Sprintf("/petrock_example_feature_name/?page=%d&pageSize=%d", prevPage, result.PageSize)
					return html.A(
						g.Attr("href", prevURL),
						g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50"),
						g.Text("Previous"),
					)
				}()),
				
				// Next page button
				html.Div(func() g.Node {
					if result.Page >= totalPages {
						return html.Span(
							g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-300 bg-white border border-slate-300 cursor-not-allowed rounded-md"),
							g.Text("Next"),
						)
					}
					
					nextPage := result.Page + 1
					nextURL := fmt.Sprintf("/petrock_example_feature_name/?page=%d&pageSize=%d", nextPage, result.PageSize)
					return html.A(
						g.Attr("href", nextURL),
						g.Attr("class", "inline-flex items-center px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-md hover:bg-slate-50"),
						g.Text("Next"),
					)
				}()),
			)
		}(),
	)
}

// pluralize returns "s" if count is not 1, otherwise returns empty string
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// formFieldClass returns the appropriate CSS class for a form field based on its error state
func formFieldClass(form *core.Form, fieldName string) string {
	if form.HasError(fieldName) {
		return "block w-full rounded-md sm:text-sm border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 focus:border-red-500"
	}
	return "block w-full rounded-md border-slate-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
}

// csrfField returns a hidden input field for CSRF protection
func csrfField(token string) g.Node {
	return html.Input(
		g.Attr("type", "hidden"),
		g.Attr("name", "csrf_token"),
		g.Attr("value", token),
	)
}

// formErrorDisplay renders an error message for a form field
func formErrorDisplay(form *core.Form, fieldName string) g.Node {
	if !form.HasError(fieldName) {
		return nil
	}
	return html.Div(
		g.Attr("class", "mt-2 text-sm text-red-600"),
		g.Text(form.GetError(fieldName)),
	)
}

// successAlert renders a success message
func successAlert(message string) g.Node {
	if message == "" {
		return nil
	}
	return html.Div(
		g.Attr("class", "mb-4 p-4 rounded-md bg-green-50 border border-green-200"),
		html.Div(
			g.Attr("class", "flex"),
			html.Div(
				g.Attr("class", "flex-shrink-0"),
				html.Svg(
					g.Attr("class", "h-5 w-5 text-green-400"),
					g.Attr("viewBox", "0 0 20 20"),
					g.Attr("fill", "currentColor"),
					html.Path(
						g.Attr("fill-rule", "evenodd"),
						g.Attr("d", "M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"),
						g.Attr("clip-rule", "evenodd"),
					),
				),
			),
			html.Div(
				g.Attr("class", "ml-3"),
				html.P(
					g.Attr("class", "text-sm font-medium text-green-800"),
					g.Text(message),
				),
			),
		),
	)
}

// DeleteConfirmForm renders a form to confirm deletion of an item.
func DeleteConfirmForm(item *Result, csrfToken string) g.Node {
	return html.Div(
		// Form container with back link
		g.Attr("class", "space-y-8"),
		
		// Back navigation
		html.Div(
			g.Attr("class", "flex justify-end"),
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID),
				g.Attr("class", "text-sm font-medium text-indigo-600 hover:text-indigo-500"),
				html.Span(g.Attr("aria-hidden", "true"), g.Text("←")),
				g.Text(" Back"),
			),
		),
		
		// Form with warning card
		html.Form(
			g.Attr("class", "bg-white shadow-sm border border-slate-200 rounded-lg overflow-hidden"),
			// Form attributes
			html.Action("/petrock_example_feature_name/"+item.ID+"/delete"),
			html.Method("POST"),

			// Form header
			html.Div(
				g.Attr("class", "border-b border-slate-200 bg-red-50 px-4 py-5 sm:px-6"),
				html.H3(
					g.Attr("class", "text-lg font-medium leading-6 text-red-800"),
					g.Text("Confirm Deletion"),
				),
				html.P(
					g.Attr("class", "mt-1 text-sm text-red-700"),
					g.Text("This action cannot be undone. Please confirm that you want to permanently delete this item."),
				),
			),

			// Item details
			html.Div(
				g.Attr("class", "px-4 py-5 sm:p-6 border-b border-slate-200"),
				
				// CSRF Token
				csrfField(csrfToken),
				
				// Item details
				html.Dl(
					g.Attr("class", "grid grid-cols-1 gap-x-4 gap-y-4"),
					
					// ID field
					html.Div(
						g.Attr("class", "col-span-1"),
						html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("ID")),
						html.Dd(g.Attr("class", "mt-1 text-sm text-slate-900"), g.Text(item.ID)),
					),
					
					// Name field
					html.Div(
						g.Attr("class", "col-span-1"),
						html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Name")),
						html.Dd(g.Attr("class", "mt-1 text-sm text-slate-900 font-medium"), g.Text(item.Name)),
					),
					
					// Description field
					html.Div(
						g.Attr("class", "col-span-1"),
						html.Dt(g.Attr("class", "text-sm font-medium text-slate-500"), g.Text("Description")),
						html.Dd(
							g.Attr("class", "mt-1 text-sm text-slate-900 whitespace-pre-wrap"), 
							g.Text(item.Description),
						),
					),
				),
			),
			
			// Form footer with action buttons
			html.Div(
				g.Attr("class", "px-4 py-3 bg-slate-50 sm:px-6 border-t border-slate-200 flex flex-col-reverse sm:flex-row sm:justify-between sm:space-x-4"),
				// Cancel button
				html.A(
					g.Attr("href", "/petrock_example_feature_name/"+item.ID),
					g.Attr("class", "w-full sm:w-auto mt-3 sm:mt-0 inline-flex justify-center py-2 px-4 border border-slate-300 shadow-sm text-sm font-medium rounded-md text-slate-700 bg-white hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
					g.Text("Cancel"),
				),
				// Delete button
				html.Button(
					g.Attr("type", "submit"),
					g.Attr("class", "w-full sm:w-auto inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"),
					g.Text("Delete Item"),
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
