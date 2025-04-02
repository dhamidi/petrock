# Web Server

The web server is the entry point for HTTP requests in a Petrock application. It handles routing, middleware application, and integrates with the form system.

## Router

The router matches HTTP requests to handlers:

```go
type Router struct {
    mux *http.ServeMux
    middleware []func(http.Handler) http.Handler
}

func NewRouter() *Router {
    return &Router{
        mux: http.NewServeMux(),
    }
}

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
    r.middleware = append(r.middleware, middleware)
}

// Register a standard handler
func (r *Router) Handle(pattern string, handler http.Handler) {
    r.mux.Handle(pattern, handler)
}

// Register a function handler
func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
    r.mux.HandleFunc(pattern, handler)
}

// Register a GET handler
func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) {
    r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handler(w, r)
    })
}

// Register a form handler
func (r *Router) Form(pattern string, form Form, handler func(form Form, ctx *Context)) {
    r.mux.Handle(pattern, FormMiddleware(form, handler))
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Apply middleware in reverse order
    var handler http.Handler = r.mux
    for i := len(r.middleware) - 1; i >= 0; i-- {
        handler = r.middleware[i](handler)
    }
    handler.ServeHTTP(w, req)
}
```

## Context

The `Context` object provides access to request and response data:

```go
type Context struct {
    Request        *http.Request
    ResponseWriter http.ResponseWriter
    CurrentUser    *User
    Flash          map[string]string
    Params         map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
    return &Context{
        Request:        r,
        ResponseWriter: w,
        Flash:          make(map[string]string),
        Params:         make(map[string]string),
    }
}
```

## Response Helpers

Helper functions make it easier to render responses:

```go
// Render a gomponents Node as HTML
func Render(w http.ResponseWriter, node gomponents.Node) {
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    node.Render(w)
}

// Render an error
func RenderError(w http.ResponseWriter, err error) {
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusInternalServerError)
    
    errorComponent := ui.Error(err.Error())
    errorComponent.Render(w)
}

// Redirect to another page
func Redirect(w http.ResponseWriter, r *http.Request, path string) {
    http.Redirect(w, r, path, http.StatusSeeOther)
}

// Stream a Turbo response
func TurboStream(w http.ResponseWriter, action string, target string, content gomponents.Node) {
    w.Header().Set("Content-Type", "text/vnd.turbo-stream.html")
    w.WriteHeader(http.StatusOK)
    
    streamComponent := ui.TurboStream(action, target, content)
    streamComponent.Render(w)
}
```

## Application Server

The main application server ties everything together:

```go
type Application struct {
    router    *Router
    logStore  LogStore
    modules   []Module
    adminPort int
}

func NewApplication() *Application {
    return &Application{
        router:    NewRouter(),
        adminPort: 3001,
    }
}

func (a *Application) SetLogStore(store LogStore) {
    a.logStore = store
}

func (a *Application) RegisterModule(module Module) {
    a.modules = append(a.modules, module)
}

func (a *Application) Start(addr string) error {
    // Initialize log store if not set
    if a.logStore == nil {
        store, err := NewSQLiteLogStore("petrock.db")
        if err != nil {
            return err
        }
        a.logStore = store
    }
    
    // Initialize modules
    for _, module := range a.modules {
        if err := module.Init(a); err != nil {
            return err
        }
        module.RegisterRoutes(a.router)
    }
    
    // Set up admin server
    go a.startAdminServer()
    
    // Start main server
    return http.ListenAndServe(addr, a.router)
}

func (a *Application) startAdminServer() {
    adminRouter := NewRouter()
    
    // Register admin routes
    adminRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
        Render(w, ui.AdminDashboard())
    })
    
    // Register module-specific admin routes
    for _, module := range a.modules {
        module.RegisterAdminRoutes(adminRouter)
    }
    
    // Start admin server
    http.ListenAndServe(fmt.Sprintf(":%d", a.adminPort), adminRouter)
}
```

## Static File Handling

Petrock also handles static files:

```go
func (a *Application) ServeStatic(urlPrefix, dirPath string) {
    fileServer := http.FileServer(http.Dir(dirPath))
    a.router.Handle(urlPrefix, http.StripPrefix(urlPrefix, fileServer))
}
```

## Middleware

Common middleware functions can be registered with the router:

```go
// Authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check for session cookie
        cookie, err := r.Cookie("session")
        if err != nil {
            // Redirect to login
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        
        // Validate session
        user, err := auth.GetUserBySession(cookie.Value)
        if err != nil {
            // Clear invalid cookie
            http.SetCookie(w, &http.Cookie{
                Name:   "session",
                Value:  "",
                Path:   "/",
                MaxAge: -1,
            })
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        
        // Store user in context
        ctx := context.WithValue(r.Context(), "currentUser", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Request logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        rw := NewResponseWriter(w)
        
        // Call next handler
        next.ServeHTTP(rw, r)
        
        // Log request
        duration := time.Since(start)
        log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.StatusCode, duration)
    })
}
```

## Hotwired Integration

The web server integrates with the Hotwired stack (Turbo and Stimulus):

```go
// Add Hotwired assets to the page
func (a *Application) UseHotwired() {
    // Include Turbo and Stimulus scripts in all pages
    a.router.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only intercept HTML responses
            ww := NewResponseWriterWrapper(w)
            next.ServeHTTP(ww, r)
            
            // If this isn't an HTML response, do nothing
            if !strings.HasPrefix(ww.Header().Get("Content-Type"), "text/html") {
                return
            }
            
            // Inject scripts before </body>
            body := ww.Body()
            body = strings.Replace(body, "</body>", `
                <script src="https://unpkg.com/@hotwired/turbo@7.1.0/dist/turbo.es2017-umd.js"></script>
                <script src="https://unpkg.com/stimulus@3.0.1/dist/stimulus.umd.js"></script>
                <script src="/assets/application.js"></script>
                </body>
            `, 1)
            
            w.Header().Set("Content-Type", "text/html")
            w.Header().Set("Content-Length", strconv.Itoa(len(body)))
            w.WriteHeader(ww.StatusCode)
            w.Write([]byte(body))
        })
    })
}
```
