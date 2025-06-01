# Form Validation Guide

This guide covers the new unified form validation system in Petrock, which replaces the old `core.Form` approach with a more powerful, extensible tag-based validation system.

## Overview

The new system uses struct tags to define validation rules directly on command and query fields, providing:

- **Declarative validation** - Rules defined in struct tags
- **Multiple input sources** - HTTP forms, JSON APIs, CLI arguments 
- **Extensible validators** - Add custom validation logic
- **Structured errors** - Rich error information with codes and metadata
- **Type safety** - Automatic type conversion with validation

## Basic Usage

### 1. Define Validation Tags

Add validation tags to your command or query structs:

```go
type CreateUserCommand struct {
    Name     string `json:"name" validate:"required,minlen=2,maxlen=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"min=0,max=120"`
    Active   bool   `json:"active"`
    Tags     []string `json:"tags"`
}
```

### 2. Parse and Validate Input

Replace manual form validation with the new parsing system:

```go
// Before (old system)
form := core.NewForm(r.PostForm)
form.ValidateRequired("name", "email")
if !form.IsValid() {
    // handle errors
}
cmd := CreateUserCommand{
    Name:  form.Get("name"),
    Email: form.Get("email"),
}

// After (new system)
var cmd CreateUserCommand
if err := core.ParseFromURLValues(r.PostForm, &cmd); err != nil {
    if parseErrors, ok := err.(*core.ParseErrors); ok {
        // Handle structured validation errors
        for _, e := range parseErrors.Errors {
            fmt.Printf("Field %s: %s (code: %s)\n", e.Field, e.Message, e.Code)
        }
    }
    return
}
```

## Available Validation Rules

### Built-in Validators

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field cannot be empty | `validate:"required"` |
| `minlen=N` | Minimum string length | `validate:"minlen=3"` |
| `maxlen=N` | Maximum string length | `validate:"maxlen=100"` |
| `min=N` | Minimum numeric value | `validate:"min=0"` |
| `max=N` | Maximum numeric value | `validate:"max=120"` |
| `email` | Valid email format | `validate:"email"` |

### Advanced Validators

| Tag | Description | Example |
|-----|-------------|---------|
| `confirm_field=X` | Must match another field | `validate:"confirm_field=password"` |
| `required_if=field:value` | Required when condition met | `validate:"required_if=type:premium"` |
| `message=X` | Custom error message | `validate:"required,message=Username is required"` |

### Combining Rules

Multiple validation rules can be combined with commas:

```go
type User struct {
    Username string `validate:"required,minlen=3,maxlen=20"`
    Password string `validate:"required,minlen=8"`
    Confirm  string `validate:"required,confirm_field=password"`
}
```

## Input Sources

### HTTP Forms

Parse HTML form submissions:

```go
var cmd CreateUserCommand
err := core.ParseFromURLValues(r.PostForm, &cmd)
```

### JSON APIs

Parse JSON request bodies:

```go
var cmd CreateUserCommand  
err := core.ParseFromMap(jsonData, &cmd)
```

### CLI Arguments

Parse command-line arguments:

```go
var config ServerConfig
err := core.ParseFromArgs(os.Args[1:], &config)
```

## Error Handling

### Structured Errors

The system provides rich error information:

```go
if err := core.ParseFromURLValues(values, &cmd); err != nil {
    if parseErrors, ok := err.(*core.ParseErrors); ok {
        for _, e := range parseErrors.Errors {
            log.Printf("Field: %s", e.Field)
            log.Printf("Message: %s", e.Message) 
            log.Printf("Code: %s", e.Code)
            log.Printf("Meta: %+v", e.Meta)
        }
    }
}
```

### Error Codes

Standard error codes for programmatic handling:

- `required` - Required field is missing
- `min_length` / `max_length` - String length validation
- `min_value` / `max_value` - Numeric range validation  
- `invalid_email` - Email format validation
- `invalid_type` - Type conversion error
- `field_mismatch` - Cross-field validation error

### API Error Responses

For JSON APIs, errors are returned in a structured format:

```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "email",
      "message": "Must be a valid email address",
      "code": "invalid_email"
    },
    {
      "field": "age", 
      "message": "Must be at least 0",
      "code": "min_value",
      "meta": {
        "min_value": 0,
        "actual_value": -5
      }
    }
  ]
}
```

## Custom Validators

### Creating Custom Validators

Implement the `Validator` interface:

```go
type PasswordStrengthValidator struct{}

func (v PasswordStrengthValidator) CanValidate(ctx *core.FieldContext) bool {
    return ctx.FieldType.Kind() == reflect.String && 
           ctx.GetTagBool("password_strength")
}

func (v PasswordStrengthValidator) Validate(ctx *core.FieldContext) []core.ParseError {
    str, ok := ctx.Value.(string)
    if !ok || str == "" {
        return nil
    }
    
    if len(str) < 8 {
        return []core.ParseError{{
            Field: ctx.Name,
            Message: "Password must be at least 8 characters",
            Code: "weak_password",
        }}
    }
    
    // Add more strength checks...
    return nil
}
```

### Registering Custom Validators

Add to the default parser or create a custom parser:

```go
// Global registration
core.RegisterValidator(PasswordStrengthValidator{})

// Or custom parser
parser := core.NewParser()
parser.RegisterValidator(PasswordStrengthValidator{})
```

## Migration Guide

### From Old Form System

Replace form-based validation:

```go
// Old approach
form := core.NewForm(r.PostForm)
form.ValidateRequired("name", "email")
form.ValidateEmail("email")
form.ValidateMinLength("name", 2)

if !form.IsValid() {
    // Handle errors
}

// New approach  
type FormData struct {
    Name  string `validate:"required,minlen=2"`
    Email string `validate:"required,email"`
}

var data FormData
if err := core.ParseFromURLValues(r.PostForm, &data); err != nil {
    // Handle structured errors
}
```

### Updating Templates

Convert form error display:

```html
<!-- Old template -->
{{if .Form.HasError "name"}}
    <div class="error">{{.Form.GetError "name"}}</div>
{{end}}

<!-- New template (using converted errors) -->
{{if .Form.HasError "name"}}
    <div class="error">{{.Form.GetError "name"}}</div>
{{end}}
```

Note: The new system maintains compatibility with existing form templates by converting `ParseErrors` to form errors in handlers.

## Best Practices

### Validation Tag Organization

Group related validation rules logically:

```go
type User struct {
    // Identity fields
    Username string `validate:"required,minlen=3,maxlen=20"`
    Email    string `validate:"required,email"`
    
    // Profile fields  
    FirstName string `validate:"required,minlen=1,maxlen=50"`
    LastName  string `validate:"required,minlen=1,maxlen=50"`
    
    // Settings
    Age     int  `validate:"min=13,max=120"`
    Active  bool `validate:""`
}
```

### Error Message Customization

Provide user-friendly error messages:

```go
type Registration struct {
    Username string `validate:"required,minlen=3,message=Username must be at least 3 characters"`
    Password string `validate:"required,minlen=8,message=Password must be at least 8 characters"`
}
```

### Performance Considerations

- Validation runs on every parse operation
- Custom validators should be efficient
- Consider caching compiled regex patterns
- Use appropriate validation complexity for your use case

### Testing Validation

Test validation rules comprehensively:

```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        data    map[string]interface{}
        wantErr bool
        errCode string
    }{
        {
            name: "valid user",
            data: map[string]interface{}{
                "username": "john",
                "email": "john@example.com",
            },
            wantErr: false,
        },
        {
            name: "missing username", 
            data: map[string]interface{}{
                "email": "john@example.com",
            },
            wantErr: true,
            errCode: "required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var user User
            err := core.ParseFromMap(tt.data, &user)
            
            if tt.wantErr {
                require.Error(t, err)
                parseErrors := err.(*core.ParseErrors)
                assert.Equal(t, tt.errCode, parseErrors.Errors[0].Code)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

## Troubleshooting

### Common Issues

**Validation not running:**
- Check that struct fields are exported (capitalized)
- Verify validation tags are properly formatted
- Ensure validator is registered with the parser

**Type conversion errors:**
- Check input data types match expected struct field types
- Use appropriate converters for custom types
- Validate input format before parsing

**Custom validator not working:**
- Implement `CanValidate()` method correctly
- Register validator with parser
- Check tag parsing logic

### Debugging Tips

Enable debug logging to see validation flow:

```go
parser := core.NewParser()
// Add debug logging to see which validators run
```

Inspect parsed tags:

```go
field := reflect.TypeOf(MyStruct{}).Field(0)
tags := core.StandardTagParser{}.ParseTags(field)
fmt.Printf("Parsed tags: %+v\n", tags)
```

## Examples

See `docs/plans/form-refactor/extensible_usage_examples.go` for complete examples of:

- Custom converters for enum types
- Business rule validators
- Multiple parser configurations
- CLI argument parsing
- Extended tag syntax
