# Migrating from Legacy Form System to New FormData System

The templates have been updated to use the new `ui.FormData` system instead of the legacy `core.Form` system. This provides better integration with the tag-based validation system and eliminates the conversion step.

## What Changed

### Template Functions

**Old (Legacy):**
```go
ui.FormGroupWithValidation(form *core.Form, ...)
ui.TextInputWithValidation(form *core.Form, ...)
ui.TextAreaWithValidation(form *core.Form, ...)
ui.SelectWithValidation(form *core.Form, ...)
ui.FormError(form *core.Form, ...)
```

**New:**
```go
ui.FormGroupWithValidation(formData *ui.FormData, ...)
ui.TextInputWithValidation(formData *ui.FormData, ...)
ui.TextAreaWithValidation(formData *ui.FormData, ...)
ui.SelectWithValidation(formData *ui.FormData, ...)
ui.FormError(formData *ui.FormData, ...)
```

### Handler Pattern

**Old:**
```go
// Convert ParseErrors to legacy Form
form := core.NewForm(r.PostForm)
for _, parseErr := range parseErrors.Errors {
    form.AddError(parseErr.Field, parseErr.Message)
}
return RenderPage(w, title, ItemForm(form, item, token))
```

**New:**
```go
// Convert ParseErrors to ui.ParseError format
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
return RenderPage(w, title, ItemForm(formData, item, token))
```

## Benefits

1. **Direct Integration**: No conversion between `ParseErrors` and legacy `Form` errors
2. **Richer Error Information**: Access to error codes and metadata in templates
3. **Type Safety**: Direct use of structured error types
4. **Performance**: Eliminates unnecessary object creation and conversion
5. **Consistency**: Same error types throughout the validation pipeline

## Backward Compatibility

Legacy functions are still available with the `Legacy` suffix:

- `ui.FormGroupWithValidationLegacy()`
- `ui.TextInputWithValidationLegacy()`
- `ui.TextAreaWithValidationLegacy()`
- `ui.SelectWithValidationLegacy()`
- `ui.FormErrorLegacy()`
- `components.ItemFormLegacy()`

These can be used during migration or for compatibility with existing `core.Form`-based code.

## Next Steps

To fully remove the legacy form system:

1. ✅ Update all template functions to use `ui.FormData`
2. ✅ Update all handlers to create `ui.FormData` instead of `core.Form`
3. ✅ Update form component signatures to accept `ui.FormData`
4. 🔄 Remove legacy `core.Form` validation methods (optional)
5. 🔄 Remove legacy template compatibility functions (optional)

The system now works entirely with the new `ui.FormData` approach, with legacy functions available for compatibility during migration.
