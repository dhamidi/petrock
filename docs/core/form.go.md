# core/form.go

This file provides form handling functionality for web applications, supporting both the new tag-based validation system and legacy form templates.

## Current System (Tag-Based Validation)

The current form system uses struct tags for validation and supports multiple input sources. See [Form Validation Guide](../form-validation-guide.md) for detailed documentation.

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

## Legacy Form System (Template Compatibility)

The legacy `Form` struct is maintained for template compatibility:

### Types

- `Form`: A struct holding form data and validation errors for template rendering.
    - `Values url.Values`: The raw form values, typically parsed from an HTTP request.
    - `Errors map[string][]string`: A map where keys are field names and values are slices of error messages for that field.

### Functions

- `NewForm(data url.Values) *Form`: Creates a new `Form` instance initialized with the provided `url.Values`.
- `(f *Form) HasError(field string) bool`: Checks if a specific field has any associated validation errors.
- `(f *Form) GetError(field string) string`: Returns the *first* error message for a given field. Returns an empty string if there are no errors for that field.
- `(f *Form) AddError(field, message string)`: Adds a new error message to a specific field.
- `(f *Form) IsValid() bool`: Returns `true` if the `Errors` map is empty, `false` otherwise.

### Template Integration

ParseErrors from the new system can be converted to Form errors for template compatibility:

```go
form := core.NewForm(r.PostForm)
// Convert ParseErrors to form errors
for _, parseErr := range parseErrors.Errors {
    form.AddError(parseErr.Field, parseErr.Message)
}
```

## Migration Path

New code should use the tag-based validation system. The legacy Form system remains available for template compatibility and will be maintained until templates are updated to work directly with ParseErrors.
