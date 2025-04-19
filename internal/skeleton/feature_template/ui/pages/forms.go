package pages

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// EditForm renders an HTML <form> for creating or editing an item.
func EditForm(form interface{}, item *Result, csrfToken string) g.Node {
	// Determine if we're editing (item != nil) or creating (item == nil)
	isNewItem := item == nil
	formTitle := "Create New Item"
	submitText := "Create"
	formAction := "/petrock_example_feature_name/new"
	
	// Set values for editing existing item
	var itemName, itemDescription, itemContent string
	if !isNewItem {
		formTitle = "Edit Item"
		submitText = "Update"
		formAction = "/petrock_example_feature_name/" + item.ID + "/edit"
		itemName = item.Name
		itemDescription = item.Description
		itemContent = item.Content
	}
	
	return html.Form(
		g.Attr("method", "POST"),
		g.Attr("action", formAction),
		g.Attr("class", "space-y-6"),
		
		// CSRF protection
		html.Input(
			g.Attr("type", "hidden"),
			g.Attr("name", "csrf_token"),
			g.Attr("value", csrfToken),
		),
		
		// Form title
		html.H2(
			g.Attr("class", "text-2xl font-bold leading-7 text-slate-900"),
			g.Text(formTitle),
		),
		
		// Form fields container
		html.Div(
			g.Attr("class", "space-y-6"),

			// Item name field
			html.Div(
				g.Attr("class", "space-y-2"),
				html.Label(
					g.Attr("for", "name"),
					g.Attr("class", "block text-sm font-medium text-slate-700"),
					g.Text("Name"),
				),
				html.Input(
					g.Attr("type", "text"),
					g.Attr("id", "name"),
					g.Attr("name", "name"),
					g.Attr("value", itemName),
					g.Attr("required", "required"),
					g.Attr("class", "block w-full rounded-md border-slate-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"),
				),
			),

			// Item description field
			html.Div(
				g.Attr("class", "space-y-2"),
				html.Label(
					g.Attr("for", "description"),
					g.Attr("class", "block text-sm font-medium text-slate-700"),
					g.Text("Description"),
				),
				html.Textarea(
					g.Attr("id", "description"),
					g.Attr("name", "description"),
					g.Attr("rows", "3"),
					g.Attr("class", "block w-full rounded-md border-slate-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"),
					g.Text(itemDescription),
				),
			),

			// Item content field
			html.Div(
				g.Attr("class", "space-y-2"),
				html.Label(
					g.Attr("for", "content"),
					g.Attr("class", "block text-sm font-medium text-slate-700"),
					g.Text("Content"),
				),
				html.Textarea(
					g.Attr("id", "content"),
					g.Attr("name", "content"),
					g.Attr("rows", "6"),
					g.Attr("class", "block w-full rounded-md border-slate-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"),
					g.Text(itemContent),
				),
				html.P(
					g.Attr("class", "mt-1 text-xs text-slate-500"),
					g.Text("A summary will be automatically generated for this content."),
				),
			),
		),
		
		// Form actions
		html.Div(
			g.Attr("class", "flex justify-between items-center"),
			
			// Cancel button
			html.A(
				g.Attr("href", func() string {
					if isNewItem {
						return "/petrock_example_feature_name"
					}
					return "/petrock_example_feature_name/" + item.ID
				}()),
				g.Attr("class", "rounded-md border border-slate-300 bg-white py-2 px-4 text-sm font-medium text-slate-700 shadow-sm hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"),
				g.Text("Cancel"),
			),
			
			// Submit button
			html.Button(
				g.Attr("type", "submit"),
				g.Attr("class", "rounded-md bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"),
				g.Text(submitText),
			),
		),
	)
}

// DeleteForm renders a form to confirm deletion of an item.
func DeleteForm(item *Result, csrfToken string) g.Node {
	if item == nil {
		return html.Div(
			g.Attr("class", "text-center p-8 bg-red-50 rounded-lg border border-red-100"),
			html.H2(
				g.Attr("class", "text-xl font-medium text-red-800"),
				g.Text("Item Not Found"),
			),
			html.P(
				g.Attr("class", "mt-2 text-sm text-red-700"),
				g.Text("The item you're trying to delete could not be found."),
			),
			html.Div(
				g.Attr("class", "mt-6"),
				html.A(
					g.Attr("href", "/petrock_example_feature_name"),
					g.Attr("class", "text-sm font-medium text-indigo-600 hover:text-indigo-500"),
					g.Text("Return to Items List"),
				),
			),
		)
	}
	
	return html.Form(
		g.Attr("method", "POST"),
		g.Attr("action", "/petrock_example_feature_name/"+item.ID+"/delete"),
		g.Attr("class", "space-y-6"),
		
		// CSRF protection
		html.Input(
			g.Attr("type", "hidden"),
			g.Attr("name", "csrf_token"),
			g.Attr("value", csrfToken),
		),
		
		// Warning message
		html.Div(
			g.Attr("class", "bg-red-50 border border-red-200 rounded-md p-4 mb-6"),
			html.Div(
				g.Attr("class", "flex items-start"),
				html.Div(
					g.Attr("class", "flex-shrink-0"),
					// Placeholder for warning icon
					html.Div(
						g.Attr("class", "h-5 w-5 text-red-400"),
						g.Text("⚠️"),
					),
				),
				html.Div(
					g.Attr("class", "ml-3"),
					html.H3(
						g.Attr("class", "text-sm font-medium text-red-800"),
						g.Text("Confirm Deletion"),
					),
					html.Div(
						g.Attr("class", "mt-2 text-sm text-red-700"),
						html.P(
							g.Text("Are you sure you want to delete this item? This action cannot be undone."),
						),
					),
				),
			),
		),
		
		// Item details card
		html.Div(
			g.Attr("class", "bg-white shadow overflow-hidden sm:rounded-lg"),
			html.Div(
				g.Attr("class", "px-4 py-5 sm:px-6"),
				html.H3(
					g.Attr("class", "text-lg leading-6 font-medium text-slate-900"),
					g.Text(item.Name),
				),
				html.P(
					g.Attr("class", "mt-1 max-w-2xl text-sm text-slate-500"),
					g.Text(item.Description),
				),
			),
			html.Div(
				g.Attr("class", "border-t border-slate-200 px-4 py-5 sm:p-0"),
				html.Dl(
					g.Attr("class", "sm:divide-y sm:divide-slate-200"),
					html.Div(
						g.Attr("class", "py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"),
						html.Dt(
							g.Attr("class", "text-sm font-medium text-slate-500"),
							g.Text("ID"),
						),
						html.Dd(
							g.Attr("class", "mt-1 text-sm text-slate-900 sm:mt-0 sm:col-span-2"),
							g.Text(item.ID),
						),
					),
					html.Div(
						g.Attr("class", "py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"),
						html.Dt(
							g.Attr("class", "text-sm font-medium text-slate-500"),
							g.Text("Created At"),
						),
						html.Dd(
							g.Attr("class", "mt-1 text-sm text-slate-900 sm:mt-0 sm:col-span-2"),
							g.Text(item.CreatedAt.Format("Jan 2, 2006 at 15:04")),
						),
					),
				),
			),
		),
		
		// Form actions
		html.Div(
			g.Attr("class", "flex justify-between mt-8"),
			
			// Cancel button
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID),
				g.Attr("class", "py-2 px-4 border border-slate-300 rounded-md shadow-sm text-sm font-medium text-slate-700 bg-white hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
				g.Text("Cancel"),
			),
			
			// Delete button
			html.Button(
				g.Attr("type", "submit"),
				g.Attr("class", "py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"),
				g.Text("Delete Item"),
			),
		),
	)
}