package layouts

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
	"github.com/petrock/example_module_path/core/ui"
)

// ModalLayout renders a modal overlay with backdrop and centered content
func ModalLayout(title string, content g.Node, actions ...g.Node) g.Node {
	return html.Div(
		ui.CSSClass("fixed", "inset-0", "z-50", "overflow-y-auto"),
		html.Div(
			ui.CSSClass("flex", "items-center", "justify-center", "min-h-screen", "px-4", "pt-4", "pb-20", "text-center", "sm:block", "sm:p-0"),
			
			// Backdrop
			html.Div(
				ui.CSSClass("fixed", "inset-0", "transition-opacity", "bg-gray-500", "bg-opacity-75"),
			),
			
			// Modal content
			html.Div(
				ui.CSSClass("inline-block", "align-bottom", "bg-white", "rounded-lg", "text-left", "overflow-hidden", "shadow-xl", "transform", "transition-all", "sm:my-8", "sm:align-middle", "sm:max-w-lg", "sm:w-full"),
				
				ui.Card(ui.CardProps{Variant: "default", Padding: "large"},
					ui.CardHeader(
						html.H3(ui.CSSClass("text-lg", "font-medium", "text-gray-900"), g.Text(title)),
					),
					ui.CardBody(content),
					func() g.Node {
						if len(actions) > 0 {
							return ui.CardFooter(
								ui.ButtonGroup(ui.ButtonGroupProps{
									Orientation: "horizontal",
									Spacing:     "medium",
								}, actions...),
							)
						}
						return nil
					}(),
				),
			),
		),
	)
}

// ConfirmModalLayout renders a confirmation modal with standard confirm/cancel actions
func ConfirmModalLayout(title, message string, confirmAction g.Node) g.Node {
	return ModalLayout(title,
		html.P(ui.CSSClass("text-sm", "text-gray-500"), g.Text(message)),
		html.Button(
			ui.CSSClass("w-full", "inline-flex", "justify-center", "rounded-md", "border", "border-gray-300", "shadow-sm", "px-4", "py-2", "bg-white", "text-base", "font-medium", "text-gray-700", "hover:bg-gray-50", "focus:outline-none", "focus:ring-2", "focus:ring-offset-2", "focus:ring-indigo-500", "mt-3", "sm:mt-0", "sm:ml-3", "sm:w-auto", "sm:text-sm"),
			g.Text("Cancel"),
		),
		confirmAction,
	)
}
