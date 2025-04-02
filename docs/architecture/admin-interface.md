# Admin Interface

The Petrock admin interface provides a web-based dashboard for managing and monitoring a Petrock application. It allows administrators to view and manipulate application state, inspect the message log, and perform administrative tasks.

## Admin Module

The admin interface is implemented as a module:

```go
// Admin module
type AdminModule struct {
    core      *Core
    sections  map[string]AdminSection
    actions   map[string]AdminAction
    router    *Router
    adminPort int
}

// AdminSection represents a section in the admin interface
type AdminSection struct {
    Name     string
    Title    string
    Handler  func(*Context) gomponents.Node
    Priority int
}

// AdminAction represents an action that can be performed in the admin interface
type AdminAction struct {
    Name        string
    Title       string
    Description string
    Handler     func(*Context) error
}

// NewAdminModule creates a new admin module
func NewAdminModule(core *Core, adminPort int) *AdminModule {
    return &AdminModule{
        core:      core,
        sections:  make(map[string]AdminSection),
        actions:   make(map[string]AdminAction),
        router:    NewRouter(),
        adminPort: adminPort,
    }
}
```

## Section Registration

Modules can register sections in the admin interface:

```go
// RegisterSection registers a section in the admin interface
func (m *AdminModule) RegisterSection(name, title string, handler func(*Context) gomponents.Node, priority int) {
    m.sections[name] = AdminSection{
        Name:     name,
        Title:    title,
        Handler:  handler,
        Priority: priority,
    }
}

// RegisterAction registers an action in the admin interface
func (m *AdminModule) RegisterAction(name, title, description string, handler func(*Context) error) {
    m.actions[name] = AdminAction{
        Name:        name,
        Title:       title,
        Description: description,
        Handler:     handler,
    }
}
```

## Admin Routes

The admin interface includes several built-in routes:

```go
// Initialize admin routes
func (m *AdminModule) Init() error {
    // Register built-in routes
    m.router.Get("/", m.HandleDashboard)
    m.router.Get("/log", m.HandleLog)
    m.router.Get("/modules", m.HandleModules)
    m.router.Get("/actions", m.HandleActions)
    m.router.Post("/actions/:name", m.HandleRunAction)
    
    // Register section routes
    for name := range m.sections {
        m.router.Get("/sections/"+name, m.HandleSection(name))
    }
    
    // Start admin server
    go m.startAdminServer()
    
    return nil
}

// Start the admin server
func (m *AdminModule) startAdminServer() {
    addr := fmt.Sprintf(":%d", m.adminPort)
    log.Printf("Admin server listening on %s", addr)
    if err := http.ListenAndServe(addr, m.router); err != nil {
        log.Printf("Admin server error: %v", err)
    }
}
```

## Dashboard Handler

The dashboard shows an overview of the application:

```go
// HandleDashboard renders the admin dashboard
func (m *AdminModule) HandleDashboard(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r)
    
    // Collect system information
    info := map[string]interface{}{
        "uptime":        time.Since(m.core.startTime).String(),
        "goroutines":    runtime.NumGoroutine(),
        "message_count": m.getMessageCount(),
        "modules":       len(m.core.modules),
    }
    
    // Render dashboard
    Render(w, m.adminLayout(
        "Dashboard",
        elem.Div(
            attr.Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8"),
            m.statCard("Uptime", info["uptime"].(string)),
            m.statCard("Goroutines", fmt.Sprintf("%d", info["goroutines"])),
            m.statCard("Messages", fmt.Sprintf("%d", info["message_count"])),
            m.statCard("Modules", fmt.Sprintf("%d", info["modules"])),
        ),
        elem.Div(
            attr.Class("grid grid-cols-1 md:grid-cols-2 gap-4"),
            m.recentMessagesCard(),
            m.actionsCard(),
        ),
    ))
}
```

## Log Handler

The log handler shows the message log:

```go
// HandleLog renders the message log
func (m *AdminModule) HandleLog(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r)
    
    // Get query parameters
    typeFilter := r.URL.Query().Get("type")
    entityFilter := r.URL.Query().Get("entity")
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit <= 0 {
        limit = 100
    }
    
    // Get messages from the log
    messages, err := m.getMessages(typeFilter, entityFilter, limit)
    if err != nil {
        RenderError(w, err)
        return
    }
    
    // Render log page
    Render(w, m.adminLayout(
        "Message Log",
        elem.Div(
            attr.Class("mb-4"),
            elem.Form(
                attr.Method("get"),
                attr.Class("flex space-x-2"),
                elem.Div(
                    attr.Class("flex-1"),
                    elem.Label(
                        attr.For("type"),
                        attr.Class("block text-sm font-medium text-gray-700"),
                        text.Text("Type Filter"),
                    ),
                    elem.Input(
                        attr.Type("text"),
                        attr.ID("type"),
                        attr.Name("type"),
                        attr.Value(typeFilter),
                        attr.Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"),
                    ),
                ),
                elem.Div(
                    attr.Class("flex-1"),
                    elem.Label(
                        attr.For("entity"),
                        attr.Class("block text-sm font-medium text-gray-700"),
                        text.Text("Entity Filter"),
                    ),
                    elem.Input(
                        attr.Type("text"),
                        attr.ID("entity"),
                        attr.Name("entity"),
                        attr.Value(entityFilter),
                        attr.Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"),
                    ),
                ),
                elem.Div(
                    attr.Class("flex-1"),
                    elem.Label(
                        attr.For("limit"),
                        attr.Class("block text-sm font-medium text-gray-700"),
                        text.Text("Limit"),
                    ),
                    elem.Input(
                        attr.Type("number"),
                        attr.ID("limit"),
                        attr.Name("limit"),
                        attr.Value(strconv.Itoa(limit)),
                        attr.Class("mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"),
                    ),
                ),
                elem.Div(
                    attr.Class("flex items-end"),
                    elem.Button(
                        attr.Type("submit"),
                        attr.Class("bg-indigo-600 px-4 py-2 text-white rounded-md"),
                        text.Text("Filter"),
                    ),
                ),
            ),
        ),
        m.messagesTable(messages),
    ))
}
```

## Modules Handler

The modules handler shows information about registered modules:

```go
// HandleModules renders the modules page
func (m *AdminModule) HandleModules(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r)
    
    // Collect module information
    var modules []map[string]interface{}
    for _, mod := range m.core.modules {
        modules = append(modules, map[string]interface{}{
            "name":        mod.Name(),
            "commands":    mod.Commands(),
            "has_routes":  hasMethod(mod, "RegisterRoutes"),
            "has_admin":   hasMethod(mod, "RegisterAdminSections"),
            "has_execute": hasMethod(mod, "Execute"),
            "has_query":   hasMethod(mod, "Query"),
        })
    }
    
    // Render modules page
    Render(w, m.adminLayout(
        "Modules",
        m.modulesTable(modules),
    ))
}
```

## Actions Handler

The actions handler shows available admin actions:

```go
// HandleActions renders the actions page
func (m *AdminModule) HandleActions(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r)
    
    // Render actions page
    Render(w, m.adminLayout(
        "Actions",
        m.actionsTable(),
    ))
}

// HandleRunAction executes an admin action
func (m *AdminModule) HandleRunAction(w http.ResponseWriter, r *http.Request) {
    ctx := NewContext(w, r)
    
    // Get action name from URL
    name := chi.URLParam(r, "name")
    
    // Find the action
    action, ok := m.actions[name]
    if !ok {
        http.NotFound(w, r)
        return
    }
    
    // Execute the action
    err := action.Handler(ctx)
    if err != nil {
        RenderError(w, err)
        return
    }
    
    // Redirect back to actions page
    http.Redirect(w, r, "/admin/actions", http.StatusSeeOther)
}
```

## Section Handler

The section handler renders a specific admin section:

```go
// HandleSection returns a handler for a specific section
func (m *AdminModule) HandleSection(name string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := NewContext(w, r)
        
        // Find the section
        section, ok := m.sections[name]
        if !ok {
            http.NotFound(w, r)
            return
        }
        
        // Render the section
        Render(w, m.adminLayout(
            section.Title,
            section.Handler(ctx),
        ))
    }
}
```

## Admin Layout

The admin interface uses a consistent layout:

```go
// adminLayout creates the admin interface layout
func (m *AdminModule) adminLayout(title string, content ...gomponents.Node) gomponents.Node {
    return ui.Page(
        "Petrock Admin - "+title,
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
                    m.sidebarLink("/admin", "Dashboard", title == "Dashboard"),
                    m.sidebarLink("/admin/log", "Message Log", title == "Message Log"),
                    m.sidebarLink("/admin/modules", "Modules", title == "Modules"),
                    m.sidebarLink("/admin/actions", "Actions", title == "Actions"),
                    // Render section links
                    elem.Hr(attr.Class("my-4 border-gray-200")),
                    m.sectionLinks(title),
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
                            text.Text(title),
                        ),
                    ),
                ),
                elem.Main(
                    attr.Class("max-w-7xl mx-auto py-6 sm:px-6 lg:px-8"),
                    content...,
                ),
            ),
        ),
    )
}
```

## UI Components

The admin interface includes several UI components:

```go
// sidebarLink creates a sidebar link
func (m *AdminModule) sidebarLink(href, label string, active bool) gomponents.Node {
    class := "block px-4 py-2 text-sm font-medium rounded-md "
    if active {
        class += "text-gray-900 bg-gray-100"
    } else {
        class += "text-gray-600 hover:bg-gray-50 hover:text-gray-900"
    }
    
    return elem.A(
        attr.Href(href),
        attr.Class(class),
        text.Text(label),
    )
}

// sectionLinks renders links to all registered sections
func (m *AdminModule) sectionLinks(currentTitle string) gomponents.Node {
    // Sort sections by priority
    var sections []AdminSection
    for _, section := range m.sections {
        sections = append(sections, section)
    }
    sort.Slice(sections, func(i, j int) bool {
        return sections[i].Priority < sections[j].Priority
    })
    
    // Create links
    var links []gomponents.Node
    for _, section := range sections {
        links = append(links, m.sidebarLink(
            "/admin/sections/"+section.Name,
            section.Title,
            currentTitle == section.Title,
        ))
    }
    
    return elem.Div(links...)
}

// statCard creates a stat card
func (m *AdminModule) statCard(label, value string) gomponents.Node {
    return elem.Div(
        attr.Class("bg-white overflow-hidden shadow rounded-lg"),
        elem.Div(
            attr.Class("px-4 py-5 sm:p-6"),
            elem.Div(
                attr.Class("flex items-center"),
                elem.Div(
                    attr.Class("flex-shrink-0 text-3xl font-semibold text-gray-900"),
                    text.Text(value),
                ),
            ),
            elem.Div(
                attr.Class("mt-2"),
                elem.P(
                    attr.Class("text-sm font-medium text-gray-500 truncate"),
                    text.Text(label),
                ),
            ),
        ),
    )
}
```
