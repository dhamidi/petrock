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

// successAlert renders a success message using the new Alert component
func SuccessAlert(message string) g.Node {
	if message == "" {
		return nil
	}
	return ui.Alert(ui.AlertProps{
		Type:        "success",
		Title:       "Success",
		Message:     message,
		Dismissible: false,
	})
}

// NewItemButton renders a button or link to navigate to the item creation page/view.
func NewItemButton() g.Node {
	return html.A(
		html.Href("/petrock_example_feature_name/new"),
		ui.Button(ui.ButtonProps{
			Variant: "primary",
			Size:    "medium",
		}, g.Text("New Item")),
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

	return ui.Container(ui.ContainerProps{Variant: "default"},
		// Navigation breadcrumbs
		ui.Breadcrumbs(ui.BreadcrumbsProps{
			Items: []ui.BreadcrumbItem{
				{Label: "Items", Href: "/petrock_example_feature_name"},
				{Label: func() string {
					if isEdit {
						return "Edit " + item.Name
					}
					return "Create New Item"
				}(), Current: true},
			},
		}),

		ui.Section(ui.SectionProps{Heading: title, Level: 1},
			html.P(
				ui.CSSClass("text-lg", "text-gray-600", "mb-6"),
				g.Text(func() string {
					if isEdit {
						return "Update the item details below."
					}
					return "Fill out the form below to create a new item."
				}()),
			),

			ui.Card(ui.CardProps{Variant: "default", Padding: "large"},
				html.Form(
					html.Action(actionURL),
					html.Method("POST"),
					ui.CSSClass("space-y-6"),

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

					// Form actions
					ui.ButtonGroup(ui.ButtonGroupProps{
						Orientation: "horizontal",
						Spacing:     "medium",
					},
						html.A(
							html.Href(backLink(isEdit, item)),
							ui.Button(ui.ButtonProps{
								Variant: "secondary",
								Size:    "medium",
							}, g.Text("Cancel")),
						),
						ui.Button(ui.ButtonProps{
							Type:    "submit",
							Variant: "primary",
							Size:    "medium",
						}, g.Text(submitLabel)),
					),
				),
			),
		),
	)
}

// DeleteConfirmForm renders a form to confirm deletion of an item.
func DeleteConfirmForm(item *state.Item, csrfToken string) g.Node {
	return ui.Container(ui.ContainerProps{Variant: "default"},
		// Navigation breadcrumbs
		ui.Breadcrumbs(ui.BreadcrumbsProps{
			Items: []ui.BreadcrumbItem{
				{Label: "Items", Href: "/petrock_example_feature_name"},
				{Label: item.Name, Href: "/petrock_example_feature_name/" + item.ID},
				{Label: "Delete", Current: true},
			},
		}),

		ui.Section(ui.SectionProps{Heading: "Delete Item", Level: 1},
			// Warning alert
			ui.Alert(ui.AlertProps{
				Type:        "warning",
				Title:       "Confirm Deletion",
				Message:     "This action cannot be undone. Please confirm that you want to permanently delete this item.",
				Dismissible: false,
			}),

			// Item details card
			ui.Card(ui.CardProps{Variant: "default", Padding: "large"},
				ui.CardHeader(
					html.H3(ui.CSSClass("text-lg", "font-medium"), g.Text(item.Name)),
					html.P(ui.CSSClass("text-sm", "text-gray-500"), g.Text("Item details")),
				),
				ui.CardBody(
					html.Dl(ui.CSSClass("grid", "grid-cols-1", "gap-4"),
						html.Div(
							html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("ID")),
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900"), g.Text(item.ID)),
						),
						html.Div(
							html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Name")),
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900", "font-medium"), g.Text(item.Name)),
						),
						html.Div(
							html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Description")),
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900", "whitespace-pre-wrap"), g.Text(item.Description)),
						),
					),
				),
				ui.CardFooter(
					html.Form(
						html.Action("/petrock_example_feature_name/"+item.ID+"/delete"),
						html.Method("POST"),
						ui.CSSClass("w-full"),

						// CSRF Token
						ui.CSRFInput(csrfToken),

						ui.ButtonGroup(ui.ButtonGroupProps{
							Orientation: "horizontal",
							Spacing:     "medium",
						},
							html.A(
								html.Href("/petrock_example_feature_name/"+item.ID),
								ui.Button(ui.ButtonProps{
									Variant: "secondary",
									Size:    "medium",
								}, g.Text("Cancel")),
							),
							ui.Button(ui.ButtonProps{
								Type:    "submit",
								Variant: "danger",
								Size:    "medium",
							}, g.Text("Delete Item")),
						),
					),
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
