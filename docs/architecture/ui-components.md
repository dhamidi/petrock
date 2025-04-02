# UI Components

Petrock uses the gomponents library to build UI components, combined with TailwindCSS for styling. These components provide a consistent way to build user interfaces.

## Base Components

```go
// Common components
func Page(title string, content ...gomponents.Node) gomponents.Node {
    return elem.HTML(
        elem.Head(
            elem.Title(text.Text(title)),
            elem.Link(attr.Rel("stylesheet"), attr.Href("/assets/tailwind.css")),
            elem.Script(attr.Src("https://unpkg.com/@hotwired/turbo@7.1.0/dist/turbo.es2017-umd.js")),
            elem.Script(attr.Src("https://unpkg.com/stimulus@3.0.1/dist/stimulus.umd.js")),
        ),
        elem.Body(
            attr.Class("bg-gray-100 min-h-screen"),
            content...,
        ),
    )
}

// Layout with navigation
func Layout(title string, content ...gomponents.Node) gomponents.Node {
    return Page(
        title,
        elem.Nav(
            attr.Class("bg-white shadow"),
            elem.Div(
                attr.Class("max-w-7xl mx-auto px-4"),
                elem.Div(
                    attr.Class("flex justify-between h-16"),
                    elem.Div(
                        attr.Class("flex"),
                        elem.Div(
                            attr.Class("flex-shrink-0 flex items-center"),
                            elem.A(
                                attr.Href("/"),
                                attr.Class("text-xl font-bold text-gray-800"),
                                text.Text("Petrock"),
                            ),
                        ),
                    ),
                ),
            ),
        ),
        elem.Main(
            attr.Class("max-w-7xl mx-auto py-6 sm:px-6 lg:px-8"),
            content...,
        ),
    )
}

// Container for content
func Container(children ...gomponents.Node) gomponents.Node {
    return elem.Div(
        attr.Class("bg-white overflow-hidden shadow rounded-lg"),
        elem.Div(
            attr.Class("px-4 py-5 sm:p-6"),
            children...,
        ),
    )
}
```

## Form Components

```go
// Form with Turbo integration
func Form(action string, method string, turboFrame string, children ...gomponents.Node) gomponents.Node {
    attrs := []attr.Attribute{
        attr.Action(action),
        attr.Method(method),
    }
    
    if turboFrame != "" {
        attrs = append(attrs, attr.Custom("data-turbo-frame", turboFrame))
    }
    
    return elem.Form(append(attrs, children...)...)
}

// Input field with label and error handling
func FormField(name, label, value string, hasError bool, errorMsg string) gomponents.Node {
    var errorNode gomponents.Node
    var inputClass string
    
    if hasError {
        errorNode = elem.P(
            attr.Class("mt-2 text-sm text-red-600"),
            text.Text(errorMsg),
        )
        inputClass = "mt-1 block w-full rounded-md border-red-300 shadow-sm focus:border-red-500 focus:ring-red-500"
    } else {
        errorNode = nil
        inputClass = "mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
    }
    
    return elem.Div(
        attr.Class("mb-4"),
        elem.Label(
            attr.For(name),
            attr.Class("block text-sm font-medium text-gray-700"),
            text.Text(label),
        ),
        elem.Input(
            attr.Type("text"),
            attr.ID(name),
            attr.Name(name),
            attr.Value(value),
            attr.Class(inputClass),
        ),
        errorNode,
    )
}

// Textarea field with label and error handling
func TextareaField(name, label, value string, hasError bool, errorMsg string) gomponents.Node {
    var errorNode gomponents.Node
    var inputClass string
    
    if hasError {
        errorNode = elem.P(
            attr.Class("mt-2 text-sm text-red-600"),
            text.Text(errorMsg),
        )
        inputClass = "mt-1 block w-full rounded-md border-red-300 shadow-sm focus:border-red-500 focus:ring-red-500"
    } else {
        errorNode = nil
        inputClass = "mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
    }
    
    return elem.Div(
        attr.Class("mb-4"),
        elem.Label(
            attr.For(name),
            attr.Class("block text-sm font-medium text-gray-700"),
            text.Text(label),
        ),
        elem.Textarea(
            attr.ID(name),
            attr.Name(name),
            attr.Class(inputClass),
            attr.Rows("4"),
            text.Text(value),
        ),
        errorNode,
    )
}

// Button component
func Button(text string, attrs ...attr.Attribute) gomponents.Node {
    defaultAttrs := []attr.Attribute{
        attr.Type("submit"),
        attr.Class("inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"),
    }
    
    allAttrs := append(defaultAttrs, attrs...)
    
    return elem.Button(
        append(allAttrs, text.Text(text))...,
    )
}

// Error message component
func ErrorMessage(message string) gomponents.Node {
    return elem.Div(
        attr.Class("bg-red-50 border-l-4 border-red-400 p-4"),
        elem.Div(
            attr.Class("flex"),
            elem.Div(
                attr.Class("flex-shrink-0"),
                // Error icon could be added here
            ),
            elem.Div(
                attr.Class("ml-3"),
                elem.P(
                    attr.Class("text-sm text-red-700"),
                    text.Text(message),
                ),
            ),
        ),
    )
}
```

## Turbo Components

```go
// Turbo Frame component
func TurboFrame(id string, children ...gomponents.Node) gomponents.Node {
    return elem.Custom("turbo-frame", 
        attr.ID(id),
        children...,
    )
}

// Turbo Stream component
func TurboStream(action string, target string, content gomponents.Node) gomponents.Node {
    return elem.Custom("turbo-stream",
        attr.Custom("action", action),
        attr.Custom("target", target),
        elem.Template(
            content,
        ),
    )
}
```

## Data Display Components

```go
// Card component
func Card(title string, content ...gomponents.Node) gomponents.Node {
    return elem.Div(
        attr.Class("bg-white overflow-hidden shadow rounded-lg divide-y divide-gray-200"),
        elem.Div(
            attr.Class("px-4 py-5 sm:px-6"),
            elem.H3(
                attr.Class("text-lg leading-6 font-medium text-gray-900"),
                text.Text(title),
            ),
        ),
        elem.Div(
            attr.Class("px-4 py-5 sm:p-6"),
            content...,
        ),
    )
}

// Table component
func Table(headers []string, rows [][]gomponents.Node) gomponents.Node {
    // Create header cells
    var headerCells []gomponents.Node
    for _, header := range headers {
        headerCells = append(headerCells, elem.Th(
            attr.Scope("col"),
            attr.Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"),
            text.Text(header),
        ))
    }
    
    // Create table rows
    var tableRows []gomponents.Node
    for _, row := range rows {
        var cells []gomponents.Node
        for _, cell := range row {
            cells = append(cells, elem.Td(
                attr.Class("px-6 py-4 whitespace-nowrap text-sm text-gray-500"),
                cell,
            ))
        }
        tableRows = append(tableRows, elem.Tr(cells...))
    }
    
    return elem.Div(
        attr.Class("flex flex-col"),
        elem.Div(
            attr.Class("-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8"),
            elem.Div(
                attr.Class("py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"),
                elem.Div(
                    attr.Class("shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"),
                    elem.Table(
                        attr.Class("min-w-full divide-y divide-gray-200"),
                        elem.Thead(
                            attr.Class("bg-gray-50"),
                            elem.Tr(headerCells...),
                        ),
                        elem.Tbody(
                            attr.Class("bg-white divide-y divide-gray-200"),
                            tableRows...,
                        ),
                    ),
                ),
            ),
        ),
    )
}

// Alert component
func Alert(message string, type_ string) gomponents.Node {
    var classes string
    switch type_ {
    case "success":
        classes = "bg-green-50 border-l-4 border-green-400 p-4"
    case "error":
        classes = "bg-red-50 border-l-4 border-red-400 p-4"
    case "warning":
        classes = "bg-yellow-50 border-l-4 border-yellow-400 p-4"
    default:
        classes = "bg-blue-50 border-l-4 border-blue-400 p-4"
    }
    
    return elem.Div(
        attr.Class(classes),
        elem.Div(
            attr.Class("flex"),
            elem.Div(
                attr.Class("ml-3"),
                elem.P(
                    attr.Class("text-sm"),
                    text.Text(message),
                ),
            ),
        ),
    )
}

// Empty state component
func EmptyState(title, description string, action gomponents.Node) gomponents.Node {
    return elem.Div(
        attr.Class("text-center py-10"),
        elem.H3(
            attr.Class("mt-2 text-sm font-medium text-gray-900"),
            text.Text(title),
        ),
        elem.P(
            attr.Class("mt-1 text-sm text-gray-500"),
            text.Text(description),
        ),
        elem.Div(
            attr.Class("mt-6"),
            action,
        ),
    )
}
```

## Admin Interface Components

```go
// Admin dashboard
func AdminDashboard() gomponents.Node {
    return Page(
        "Admin Dashboard",
        elem.Div(
            attr.Class("flex h-screen bg-gray-100"),
            // Sidebar
            elem.Div(
                attr.Class("w-64 bg-white shadow"),
                elem.Div(
                    attr.Class("h-16 flex items-center justify-center"),
                    elem.H1(
                        attr.Class("text-xl font-bold text-gray-800"),
                        text.Text("Petrock Admin"),
                    ),
                ),
                elem.Nav(
                    attr.Class("mt-5 px-2"),
                    elem.A(
                        attr.Href("/admin"),
                        attr.Class("block px-4 py-2 text-sm font-medium text-gray-900 bg-gray-100 rounded-md"),
                        text.Text("Dashboard"),
                    ),
                    elem.A(
                        attr.Href("/admin/log"),
                        attr.Class("block px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-md"),
                        text.Text("Message Log"),
                    ),
                    elem.A(
                        attr.Href("/admin/modules"),
                        attr.Class("block px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-md"),
                        text.Text("Modules"),
                    ),
                ),
            ),
            // Main content
            elem.Div(
                attr.Class("flex-1 overflow-auto"),
                elem.Header(
                    attr.Class("bg-white shadow"),
                    elem.Div(
                        attr.Class("max-w-7xl mx-auto py-6 px-4"),
                        elem.H1(
                            attr.Class("text-3xl font-bold text-gray-900"),
                            text.Text("Dashboard"),
                        ),
                    ),
                ),
                elem.Main(
                    attr.Class("max-w-7xl mx-auto py-6 sm:px-6 lg:px-8"),
                    // Dashboard content goes here
                ),
            ),
        ),
    )
}
```
