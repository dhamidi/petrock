# Form Processing

Form handling is a key part of the Petrock framework, providing a bridge between HTTP requests and the command system. Forms handle validation, error display, and conversion of user input into commands.

## Form Interface

```go
// Form interface
type Form interface {
    Validate() []ValidationError
    ToCommand(ctx *Context) (interface{}, error)
}

// Validation utilities
type ValidationError struct {
    Field string
    Message string
}
```

## Form Middleware

Form middleware handles the processing of form submissions:

```go
// Form middleware
func FormMiddleware(form Form, handler func(form Form, ctx *Context)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := NewContext(w, r)
        
        // Only process on POST
        if r.Method != "POST" {
            // For GET, render the form
            RenderForm(w, form, nil)
            return
        }
        
        // Parse form data
        if err := r.ParseForm(); err != nil {
            http.Error(w, "Error parsing form", http.StatusBadRequest)
            return
        }
        
        // Populate form struct from request
        if err := schema.NewDecoder().Decode(form, r.PostForm); err != nil {
            http.Error(w, "Error decoding form", http.StatusBadRequest)
            return
        }
        
        // Validate form
        errors := form.Validate()
        if len(errors) > 0 {
            // Re-render form with errors
            RenderForm(w, form, errors)
            return
        }
        
        // Process valid form
        handler(form, ctx)
    }
}
```

## Form Rendering

Forms are rendered using gomponents:

```go
// Form rendering
func RenderForm(w http.ResponseWriter, form interface{}, errors []ValidationError) {
    // Convert errors to map for easier lookup
    errorMap := make(map[string]string)
    for _, err := range errors {
        errorMap[err.Field] = err.Message
    }
    
    // Reflect on form struct to extract field information
    formValue := reflect.ValueOf(form)
    formType := formValue.Type()
    
    // Build form components based on struct fields
    var fields []gomponents.Node
    for i := 0; i < formType.NumField(); i++ {
        field := formType.Field(i)
        tag := field.Tag.Get("form")
        if tag == "" {
            continue
        }
        
        value := formValue.Field(i).String()
        fieldName := field.Name
        
        // Check for error
        errorMsg, hasError := errorMap[fieldName]
        
        // Add appropriate input field
        fields = append(fields, ui.FormField(
            tag,
            field.Name,
            value,
            hasError,
            errorMsg,
        ))
    }
    
    // Render complete form
    formComponent := ui.Form(
        "/submit/path", // This would be dynamic in real implementation
        "POST",
        "form-frame",
        fields...,
    )
    
    Render(w, formComponent)
}
```

## Turbo Integration

Forms are integrated with Turbo to enable seamless updates:

```go
// Turbo-enabled form component
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
```

## Example Form Usage

Here's an example of defining and using a form:

```go
// Define the form
type CreatePostForm struct {
    Title   string `form:"title"`
    Content string `form:"content"`
}

// Validate the form
func (f CreatePostForm) Validate() []ValidationError {
    var errors []ValidationError
    
    if f.Title == "" {
        errors = append(errors, ValidationError{Field: "Title", Message: "Title is required"})
    }
    
    if len(f.Content) < 10 {
        errors = append(errors, ValidationError{Field: "Content", Message: "Content must be at least 10 characters"})
    }
    
    return errors
}

// Convert to command
func (f CreatePostForm) ToCommand(ctx *Context) (interface{}, error) {
    return &CreatePostCommand{
        ID:        core.NewID(),
        Title:     f.Title,
        Content:   f.Content,
        AuthorID:  ctx.CurrentUser.ID,
        PostedAt:  time.Now(),
    }, nil
}

// Register the form handler
func RegisterRoutes(router *Router) {
    router.Form("/posts/new", &CreatePostForm{}, HandleCreatePost)
}

// Handle the form submission
func HandleCreatePost(form Form, ctx *Context) {
    cmd, err := form.ToCommand(ctx)
    if err != nil {
        RenderError(ctx.ResponseWriter, err)
        return
    }
    
    _, err = core.Execute(cmd)
    if err != nil {
        RenderError(ctx.ResponseWriter, err)
        return
    }
    
    Redirect(ctx.ResponseWriter, ctx.Request, "/posts")
}
```

## Form Error Handling

When validation fails, the form is re-rendered with error messages, following Hotwired patterns:

1. User submits a form
2. Server validates the form
3. If validation fails:
   - The form is re-rendered with error messages
   - Turbo replaces the form in the DOM without a full page reload
4. If validation succeeds:
   - Command is created and processed
   - User is redirected or shown a success message
