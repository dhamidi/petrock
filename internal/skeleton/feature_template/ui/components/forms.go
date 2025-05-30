package components

import (
	"fmt"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
)

// formFieldClass returns the appropriate CSS class for a form field based on its error state
// DEPRECATED: Use ui.TextInputWithValidation, ui.TextAreaWithValidation instead
func FormFieldClass(form *core.Form, fieldName string) string {
	if form.HasError(fieldName) {
		return "block w-full rounded-md sm:text-sm border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 focus:border-red-500"
	}
	return "block w-full rounded-md border-slate-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
}

// csrfField returns a hidden input field for CSRF protection
// DEPRECATED: Use ui.CSRFInput instead
func CsrfField(token string) g.Node {
	return ui.CSRFInput(token)
}

// formErrorDisplay renders an error message for a form field
// DEPRECATED: Use ui.FormError instead
func FormErrorDisplay(form *core.Form, fieldName string) g.Node {
	return ui.FormError(form, fieldName)
}

// successAlert renders a success message
func SuccessAlert(message string) g.Node {
	if message == "" {
		return nil
	}
	return html.Div(
		g.Attr("class", "mb-4 p-4 rounded-md bg-green-50 border border-green-200"),
		html.Div(
			g.Attr("class", "flex"),
			html.Div(
				g.Attr("class", "ml-3"),
				html.P(
					g.Attr("class", "text-sm font-medium text-green-800"),
					html.Span(
						g.Attr("class", "inline-block mr-1"),
						g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`),
					),
					g.Text(message),
				),
			),
		),
	)
}

// NewItemButton renders a button or link to navigate to the item creation page/view.
func NewItemButton() g.Node {
	return html.A(
		html.Href("/petrock_example_feature_name/new"),
		ui.CSSClass("px-4", "py-2", "bg-green-500", "text-white", "rounded", "hover:bg-green-600"),
		g.Text("New Item"),
	)
}

// ItemForm renders an HTML <form> for creating or editing an item.
// It uses core.Form for data and error handling with new ui components.
// 'item' can be nil when creating a new item.
// 'csrfToken' should be provided by the handler.
func ItemForm(form *core.Form, item *state.Item, csrfToken string) g.Node {
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

	// Set values for form fields - use existing values from item if editing and no validation errors
	var nameValue, descriptionValue string
	if isEdit && item != nil {
		nameValue = item.Name
		descriptionValue = item.Description
	}
	// Override with form values if they exist (from failed validation)
	if form.Get("name") != "" {
		nameValue = form.Get("name")
	}
	if form.Get("description") != "" {
		descriptionValue = form.Get("description")
	}

	return html.Div(
		ui.CSSClass("space-y-8"),

		// Back navigation
		html.Div(
			ui.CSSClass("flex", "justify-end"),
			html.A(
				html.Href(backLink(isEdit, item)),
				ui.CSSClass("text-sm", "font-medium", "text-indigo-600", "hover:text-indigo-500"),
				html.Span(html.Aria("hidden", "true"), g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>`)),
				g.Text(" Back"),
			),
		),

		html.Form(
			ui.CSSClass("bg-white", "shadow-sm", "border", "border-slate-200", "rounded-lg", "overflow-hidden"),
			html.Action(actionURL),
			html.Method("POST"),

			// Form header
			html.Div(
				ui.CSSClass("border-b", "border-slate-200", "bg-slate-50", "px-4", "py-5", "sm:px-6"),
				html.H3(
					ui.CSSClass("text-lg", "font-medium", "leading-6", "text-slate-900"),
					g.Text(title),
				),
				html.P(
					ui.CSSClass("mt-1", "text-sm", "text-slate-500"),
					g.Text(func() string {
						if isEdit {
							return "Update the item details below."
						}
						return "Fill out the form below to create a new item."
					}()),
				),
			),

			// Form body
			html.Div(
				ui.CSSClass("px-4", "py-5", "sm:p-6", "space-y-6"),

				// CSRF Token
				ui.CSRFInput(csrfToken),

				// Name field using new ui components
				ui.FormGroupWithValidation(form, "name", "Name",
					ui.TextInputWithValidation(form, ui.TextInputProps{
						Name:        "name",
						Type:        "text",
						Value:       nameValue,
						Placeholder: "Enter item name",
						Required:    true,
					}),
					"A unique name for this item",
				),

				// Description field using new ui components
				ui.FormGroupWithValidation(form, "description", "Description",
					ui.TextAreaWithValidation(form, ui.TextAreaProps{
						Name:        "description",
						Value:       descriptionValue,
						Placeholder: "Enter item description",
						Rows:        4,
						Required:    true,
					}),
					"A detailed description of this item",
				),
			),

			// Form footer with submit button
			html.Div(
				ui.CSSClass("px-4", "py-3", "bg-slate-50", "text-right", "sm:px-6", "border-t", "border-slate-200"),
				ui.Button(ui.ButtonProps{
					Type:    "submit",
					Variant: "primary",
				}, g.Text(submitLabel)),
			),
		),
	)
}

// DeleteConfirmForm renders a form to confirm deletion of an item.
func DeleteConfirmForm(item *state.Item, csrfToken string) g.Node {
	return html.Div(
		// Form container with back link
		g.Attr("class", "space-y-8"),

		// Back navigation
		html.Div(
			g.Attr("class", "flex justify-end"),
			html.A(
				g.Attr("href", "/petrock_example_feature_name/"+item.ID),
				g.Attr("class", "text-sm font-medium text-indigo-600 hover:text-indigo-500"),
				html.Span(g.Attr("aria-hidden", "true"), g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>`)),
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
				CsrfField(csrfToken),

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

// backLink returns the appropriate back link URL based on whether we're in edit mode
func backLink(isEdit bool, item *state.Item) string {
	if isEdit && item != nil {
		return "/petrock_example_feature_name/" + item.ID
	}
	return "/petrock_example_feature_name"
}
