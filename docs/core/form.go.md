# core/form.go

This file previously contained the legacy Form system. All form functionality has been moved to the modern tag-based validation system.

## Current System (Tag-Based Validation)

The form system uses struct tags for validation and supports multiple input sources. See [Form Validation Guide](../form-validation-guide.md) for detailed documentation.

### Key Functions

- `ParseFromURLValues(values url.Values, dest interface{}) error`: Parses HTML form data into a struct with validation
- `ParseFromMap(data map[string]interface{}, dest interface{}) error`: Parses JSON/map data into a struct with validation  
- `ParseFromArgs(args []string, dest interface{}) error`: Parses CLI arguments into a struct with validation

### Validation Tags

Validation rules are defined using struct tags:

```go
type CreateCommand struct {
    Name        string `validate:"required,minlen=2,maxlen=100"`
    Description string `validate:"required,minlen=5,maxlen=500"`
    Email       string `validate:"email"`
}
```

### Error Handling

Returns structured errors with field-level details:

```go
if err := core.ParseFromURLValues(r.PostForm, &cmd); err != nil {
    if parseErrors, ok := err.(*core.ParseErrors); ok {
        // Handle structured validation errors
        for _, e := range parseErrors.Errors {
            fmt.Printf("Field %s: %s\n", e.Field, e.Message)
        }
    }
}
```

## Template Integration

For template rendering, use `ui.FormData` which works directly with validation errors:

```go
// Convert ParseErrors to ui.ParseError format for templates
var uiErrors []ui.ParseError
for _, parseErr := range parseErrors.Errors {
    uiErrors = append(uiErrors, ui.ParseError{
        Field:   parseErr.Field,
        Message: parseErr.Message,
        Code:    parseErr.Code,
        Meta:    parseErr.Meta,
    })
}
formData := ui.NewFormData(r.PostForm, uiErrors)

// Use in templates
ui.FormGroupWithValidation(formData, "name", "Name", ...)
```

## Migration Complete

✅ **Legacy system removed:**
- `core.Form` struct and methods
- `core.NewForm()` function  
- Legacy validation methods (`ValidateRequired`, `ValidateEmail`, etc.)

✅ **Modern system in place:**
- Tag-based validation with struct tags
- `ui.FormData` for template rendering
- Direct integration with validation pipeline
- Rich error information with codes and metadata

For complete documentation, see [Form Validation Guide](../form-validation-guide.md).
