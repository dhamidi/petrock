package pages

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"github.com/petrock/example_module_path/core/ui"
)

// EditForm renders an HTML <form> for creating or editing an item.
func EditForm(form interface{}, item *Result, csrfToken string) g.Node {
	// Cast the form to FormData
	var formData *ui.FormData
	if fd, ok := form.(*ui.FormData); ok {
		formData = fd
	} else {
		// Create empty form data if not provided
		formData = ui.NewFormData(nil, nil)
	}

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
	
	return ui.Container(ui.ContainerProps{Variant: "default"},
		// Navigation breadcrumbs
		ui.Breadcrumbs(ui.BreadcrumbsProps{
			Items: []ui.BreadcrumbItem{
				{Label: "Items", Href: "/petrock_example_feature_name"},
				{Label: func() string {
					if isNewItem {
						return "Create New Item"
					}
					return "Edit " + item.Name
				}(), Current: true},
			},
		}),

		ui.Section(ui.SectionProps{Heading: formTitle, Level: 1},
			html.P(
				ui.CSSClass("text-lg", "text-gray-600", "mb-6"),
				g.Text(func() string {
					if isNewItem {
						return "Fill out the form below to create a new item."
					}
					return "Update the item details below."
				}()),
			),

			ui.Card(ui.CardProps{Variant: "default", Padding: "large"},
				html.Form(
					html.Method("POST"),
					html.Action(formAction),
					ui.CSSClass("space-y-6"),
					
					// CSRF protection
					ui.CSRFInput(csrfToken),
					
					// Form fields with validation
					ui.FormGroupWithValidation(formData, "name", "Name",
						ui.TextInputWithValidation(formData, ui.TextInputProps{
							Name:        "name",
							Type:        "text",
							Value:       func() string {
								if formData.HasValues() {
									return ""  // Let validation function use form data
								}
								return itemName
							}(),
							Placeholder: "Enter item name",
							Required:    true,
						}),
						"A unique name for this item",
					),

					ui.FormGroupWithValidation(formData, "description", "Description",
						ui.TextAreaWithValidation(formData, ui.TextAreaProps{
							Name:        "description",
							Value:       func() string {
								if formData.HasValues() {
									return ""  // Let validation function use form data
								}
								return itemDescription
							}(),
							Placeholder: "Enter item description",
							Rows:        3,
							Required:    true,
						}),
						"A brief description of this item",
					),

					ui.FormGroupWithValidation(formData, "content", "Content",
						ui.TextAreaWithValidation(formData, ui.TextAreaProps{
							Name:        "content",
							Value:       func() string {
								if formData.HasValues() {
									return ""  // Let validation function use form data
								}
								return itemContent
							}(),
							Placeholder: "Enter item content",
							Rows:        6,
							Required:    true,
						}),
						"The main content for this item. A summary will be automatically generated.",
					),
					
					// Form actions
					ui.ButtonGroup(ui.ButtonGroupProps{
						Orientation: "horizontal",
						Spacing:     "medium",
					},
						html.A(
							html.Href(func() string {
								if isNewItem {
									return "/petrock_example_feature_name"
								}
								return "/petrock_example_feature_name/" + item.ID
							}()),
							ui.Button(ui.ButtonProps{
								Variant: "secondary",
								Size:    "medium",
							}, g.Text("Cancel")),
						),
						ui.Button(ui.ButtonProps{
							Type:    "submit",
							Variant: "primary",
							Size:    "medium",
						}, g.Text(submitText)),
					),
				),
			),
		),
	)
}

// DeleteForm renders a form to confirm deletion of an item.
func DeleteForm(item *Result, csrfToken string) g.Node {
	if item == nil {
		return ui.Container(ui.ContainerProps{Variant: "default"},
			ui.Alert(ui.AlertProps{
				Type:        "error",
				Title:       "Item Not Found",
				Message:     "The item you're trying to delete could not be found.",
				Dismissible: false,
			}),
			html.A(
				html.Href("/petrock_example_feature_name"),
				ui.Button(ui.ButtonProps{
					Variant: "primary",
					Size:    "medium",
				}, g.Text("Return to Items List")),
			),
		)
	}
	
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
				Message:     "Are you sure you want to delete this item? This action cannot be undone.",
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
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900", "font-mono"), g.Text(item.ID)),
						),
						html.Div(
							html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Name")),
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900", "font-medium"), g.Text(item.Name)),
						),
						html.Div(
							html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Description")),
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-900"), g.Text(item.Description)),
						),
						html.Div(
							html.Dt(ui.CSSClass("text-sm", "font-medium", "text-gray-500"), g.Text("Created At")),
							html.Dd(ui.CSSClass("mt-1", "text-sm", "text-gray-500"), g.Text(item.CreatedAt.Format("Jan 2, 2006 at 15:04"))),
						),
					),
				),
				ui.CardFooter(
					html.Form(
						html.Method("POST"),
						html.Action("/petrock_example_feature_name/"+item.ID+"/delete"),
						ui.CSSClass("w-full"),
						
						// CSRF protection
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