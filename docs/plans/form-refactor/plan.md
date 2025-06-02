# Form Abstraction Refactor: Unified Self-Parsing System

## Overview

This plan outlines the refactoring of Petrock's form handling system from the current basic URL values approach to a unified, extensible self-parsing system. The new system eliminates the need for separate form objects by allowing commands and queries to parse themselves from multiple input sources while maintaining extensibility through composition.

## Current State Analysis

### Problems with Current Implementation

1. **Duplication**: Commands/queries and forms have similar structure and validation
2. **Limited Sources**: Only handles `url.Values`, no support for JSON, CLI args, etc.
3. **Hard-coded Validation**: Basic validation rules baked into `core/form.go`
4. **Type Conversion**: Manual conversion between form data and command/query types
5. **Error Handling**: Simple string-based errors, no structured validation responses

### Current Form Implementation
- Location: `core/form.go` 
- Features: Basic URL values parsing, simple validation methods
- Limitations: No extensibility, limited types, poor error reporting

## Proposed Solution: Self-Parsing Commands and Queries

### Core Principles

1. **Commands and queries parse themselves** - No intermediate form objects
2. **Multiple input sources** - HTTP forms, JSON, CLI args, maps
3. **Composition over inheritance** - Pluggable converters and validators
4. **Extensible by design** - Add new types and validation rules at runtime
5. **Rich error reporting** - Structured errors with metadata and codes

### Architecture Overview

```
Input Source → Parser → Command/Query (validated & typed)
     ↑            ↑           ↑
  FormSource  Registries  Domain Object
     │            │           │
  Multiple     Converters   Business Logic
  Formats      Validators
               TagParsers
```

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1)

#### 1.1 Create New Parsing System
- **File**: `core/parser.go`
- **Reference**: [extensible_parsing_system.go](./extensible_parsing_system.go)

**Components to implement:**
- `FormSource` interface for input abstraction
- `URLValuesSource`, `MapSource` implementations
- `ParseError` and `ParseErrors` for structured error handling
- `FieldContext` for rich validation context
- Core interfaces: `Converter`, `Validator`, `TagParser`

#### 1.2 Build Component Registries
- `ConverterRegistry` - manages type converters
- `ValidatorRegistry` - manages validation rules
- `Parser` - main coordination engine

#### 1.3 Implement Built-in Components
- `BasicConverter` - handles primitives (string, int, bool, etc.)
- `TimeConverter` - handles `time.Time` with multiple formats
- `RequiredValidator` - checks for required fields
- `LengthValidator` - string length constraints
- `RangeValidator` - numeric range validation
- `EmailValidator` - email format validation
- `StandardTagParser` - parses `validate:"required,minlen=2"` syntax

**Acceptance Criteria:**

**Core Infrastructure:**
- [ ] `FormSource` interface implemented with methods: `Get(key string) ([]string, bool)`, `Keys() []string`
- [ ] `URLValuesSource` converts `url.Values` to `FormSource` correctly for all HTTP form data
- [ ] `MapSource` converts `map[string]interface{}` to `FormSource` for JSON/API data
- [ ] `ParseError` struct includes: `Field string`, `Message string`, `Code string`, `Meta map[string]interface{}`
- [ ] `ParseErrors` slice implements `Error() string` with formatted multi-field error messages
- [ ] `FieldContext` provides: field name, struct tag access, parent context chain

**Type Conversion:**
- [ ] `BasicConverter` handles: `string`, `int`, `int32`, `int64`, `float32`, `float64`, `bool`, `[]string`, `[]int`
- [ ] `TimeConverter` parses: RFC3339, `2006-01-02`, `2006-01-02 15:04:05`, custom formats via tags
- [ ] Type conversion errors include actual value and expected type in error metadata
- [ ] Converter registry allows runtime registration: `parser.RegisterConverter(reflect.TypeOf(CustomType{}), customConverter)`

**Validation System:**
- [ ] `RequiredValidator` fails on: empty strings, zero values, nil pointers, empty slices
- [ ] `LengthValidator` validates string length with `minlen=N` and `maxlen=N` tags
- [ ] `RangeValidator` validates numeric ranges with `min=N` and `max=N` tags  
- [ ] `EmailValidator` accepts valid RFC 5322 email addresses, rejects malformed ones
- [ ] Validator registry supports custom validators: `parser.RegisterValidator("custom", customValidator)`

**Tag Parsing:**
- [ ] `StandardTagParser` parses: `validate:"required,minlen=2,maxlen=50,email"`
- [ ] Tag parser handles escaped commas and equals signs in values
- [ ] Unknown validation tags generate warning logs but don't fail parsing
- [ ] Multiple validators on same field execute in tag order

**Integration Requirements:**
- [ ] `DefaultParser` instance works without configuration for basic use cases
- [ ] Parsing a struct with no tags succeeds with all converters applied
- [ ] Error messages are human-readable: "Field 'email' must be a valid email address"

### Phase 2: Feature Implementation (Week 2)

#### 2.1 Add Command Line Arguments Support
- **Reference**: Example in [extensible_usage_examples.go](./extensible_usage_examples.go)

**Implementation:**
- `ArgsSource` struct to parse CLI arguments
- Support for `--key=value`, `--key value`, and `--flag` formats
- Handle multiple values for slice fields

#### 2.2 Enhance Error Reporting
- Add metadata to `ParseError` for context
- Implement error codes for programmatic handling
- Create helper methods for common error patterns

#### 2.3 Create Extension Examples
- **Reference**: [extensible_usage_examples.go](./extensible_usage_examples.go)

**Custom Components:**
- `UserRoleConverter` - domain-specific enum converter
- `RegexValidator` - pattern matching validation
- `UniqueUsernameValidator` - business rule validation
- `ExtendedTagParser` - custom tag syntax parsing

**Acceptance Criteria:**

**CLI Arguments Support:**
- [ ] `ArgsSource` parses `--key=value` format correctly for all basic types
- [ ] `ArgsSource` parses `--key value` format with proper value assignment
- [ ] Boolean flags work with `--flag` (true) and `--no-flag` (false) syntax
- [ ] Multiple values for slices: `--tags=go --tags=web` creates `[]string{"go", "web"}`
- [ ] CLI parsing fails gracefully with helpful error for malformed arguments
- [ ] Mixed formats work together: `--name=john --age 25 --active`

**Enhanced Error Reporting:**
- [ ] `ParseError.Meta` includes: `actual_value`, `expected_type`, `constraint_value` for range errors
- [ ] Error codes are consistent: `required`, `invalid_type`, `min_length`, `max_length`, `invalid_format`
- [ ] Error messages include field path for nested structs: `user.address.city`
- [ ] `ParseErrors.Error()` formats errors as numbered list with field context
- [ ] Validation errors include the tag rule that failed: `field 'age' failed validation rule 'min=18'`

**Extension Examples:**
- [ ] `UserRoleConverter` converts strings `"admin"`, `"user"`, `"guest"` to custom `UserRole` enum
- [ ] `RegexValidator` validates fields against custom patterns: `validate:"pattern=^[A-Z]{2,3}$"`
- [ ] `UniqueUsernameValidator` demonstrates async validation pattern (mock implementation)
- [ ] `ExtendedTagParser` supports custom syntax: `validate:"businessrule:unique_username"`
- [ ] Custom components integrate without modifying core parser code

**Multi-Parser Support:**
- [ ] Create two parser instances with different validator sets, verify they work independently
- [ ] Custom components integrate without modifying core parser code

### Phase 3: Integration with Existing Codebase (Week 3)

#### 3.1 Update Core Command/Query Interfaces
- Add validation tags to existing command/query structs
- Remove old form-based parsing from handlers

#### 3.2 Refactor HTTP Handlers
**Target files:**
- `cmd/petrock_example_project_name/serve.go`
- Feature handlers in `petrock_example_feature_name/handlers/`

**Changes:**
- Replace all `core.NewForm()` calls with `ParseFromURLValues()`
- Update templates to display structured error messages
- Remove manual validation code from handlers

**Example transformation:**
```go
// Before
form := core.NewForm(r.PostForm)
form.ValidateRequired("name", "email")
if !form.IsValid() { /* handle errors */ }
cmd := commands.CreateCommand{
    Name: form.Get("name"),
    Email: form.Get("email"),
}

// After  
var cmd commands.CreateCommand
if err := ParseFromURLValues(r.PostForm, &cmd); err != nil {
    // Handle structured ParseErrors
}
```

#### 3.3 Update CLI Commands
**Target files:**
- Command handlers in `cmd/petrock_example_project_name/`

**Changes:**
- Replace cobra flag definitions with struct-based parsing
- Update help generation to reflect validation tags

#### 3.4 Update API Endpoints
**Target files:**
- API handlers that accept JSON input

**Changes:**
- Replace manual JSON unmarshaling with `ParseFromMap`
- Use consistent error response format across all endpoints

**Acceptance Criteria:**

**Handler Integration:**
- [ ] All POST/PUT handlers use `ParseFromURLValues(r.PostForm, &cmd)` pattern
- [ ] Form validation errors render with field-specific error messages in templates
- [ ] HTTP error responses include structured JSON for API endpoints
- [ ] File upload handlers preserve existing multipart/form-data behavior

**CLI Command Integration:**
- [ ] All CLI commands use struct-based parsing instead of cobra flags
- [ ] Help text generation includes validation rules from struct tags
- [ ] CLI error messages show validation failures clearly

**API Endpoint Integration:**
- [ ] JSON API handlers use `ParseFromMap` after `json.Unmarshal`
- [ ] API error responses follow consistent format across all endpoints

**Testing:**
- [ ] All existing functionality works with new parsing system
- [ ] Web UI forms show proper validation errors
- [ ] CLI commands work with new argument parsing

### Phase 4: Feature Enhancement (Week 4)

#### 4.1 Advanced Validation Features
- Cross-field validation (password confirmation)
- Conditional validation rules
- Async validation (database checks)
- Custom error messages per field

#### 4.2 Documentation and Examples
- Update feature generation templates
- Create validation best practices guide
- Document custom converter/validator patterns
- Add performance guidelines

**Acceptance Criteria:**

**Advanced Validation Features:**
- [ ] Cross-field validation: `validate:"confirm_field=password"` compares two fields
- [ ] Conditional validation: `validate:"required_if=type:premium"` validates only when another field has specific value
- [ ] Async validation: Interface supports validators that return `chan ValidationResult` for database checks
- [ ] Custom error messages: `validate:"required,message=Username is required"` overrides default messages
- [ ] Group validation: Validate related fields together with shared context and error aggregation



**Documentation and Examples:**
- [ ] README.md includes 5 common validation patterns with code examples
- [ ] API documentation generated from struct tags (if OpenAPI integration exists)
- [ ] Best practices guide covers: common mistakes, extension patterns
- [ ] Migration guide with before/after code examples for each handler type
- [ ] Troubleshooting guide for common validation errors and solutions

**Feature Template Updates:**
- [ ] Generated command structs include appropriate validation tags for common fields
- [ ] Generated handlers use new parsing system by default
- [ ] Template comments explain validation tag usage and examples
- [ ] Generated tests include validation test cases
- [ ] Feature templates demonstrate both simple and complex validation scenarios

### Phase 5: Migration and Cleanup (Week 5)

#### 5.1 Remove Old Form System
- Delete `core/form.go` 
- Remove form-related code from handlers
- Update imports throughout codebase

#### 5.2 Update Feature Templates
**Target files:**
- Feature generation templates used by `petrock feature` command

**Changes:**
- Commands/queries include validation tags
- Handlers use new parsing system
- Remove form object creation

#### 5.3 Update Documentation
- README examples
- API documentation 
- Migration guide for existing projects

**Acceptance Criteria:**

**Old System Removal:**
- [ ] `core/form.go` file deleted and removed from git history
- [ ] All imports of `core.NewForm`, `core.Form` removed from codebase
- [ ] Grep search for "form\." in codebase returns zero matches in handler files
- [ ] All references to `form.ValidateRequired`, `form.Get`, `form.IsValid` eliminated
- [ ] Build system compiles successfully with no undefined references

**Feature Template Migration:**
- [ ] `petrock feature` command generates structs with validation tags for all field types
- [ ] Generated handlers use `ParseFromURLValues` and `ParseFromMap` patterns
- [ ] Generated CLI commands use struct-based argument parsing
- [ ] Template tests include validation test cases for generated commands
- [ ] No generated code references old form system methods

**Documentation Updates:**
- [ ] README.md examples use new parsing system exclusively
- [ ] API documentation reflects new error response format
- [ ] All code examples in docs compile and run successfully
- [ ] Migration guide helps existing projects upgrade from old system
- [ ] Changelog documents breaking changes and migration path

**Quality Assurance:**
- [ ] Full integration test suite passes with new system
- [ ] Manual testing of generated features confirms functionality
- [ ] No TODO comments or placeholder code remains in core parsing system
- [ ] Code coverage for new parsing system exceeds 90%

## Implementation Details

### File Organization

```
core/
├── parser.go           # Main parsing engine (Phase 1)
├── converters.go       # Built-in type converters (Phase 1)  
├── validators.go       # Built-in validators (Phase 1)
├── sources.go          # Input source implementations (Phase 1)
└── form.go            # DELETE in Phase 5

examples/
├── extensible_parsing_system.go    # Reference implementation
└── extensible_usage_examples.go    # Usage patterns
```

### Migration Strategy for Existing Commands/Queries

1. **Add validation tags** to existing struct fields
2. **Test parsing** with current input patterns
3. **Update handlers** one at a time
4. **Verify error handling** maintains user experience

### Validation Tag Syntax

Standard format: `validate:"rule1,rule2=value,rule3"`

**Built-in rules:**
- `required` - field cannot be empty
- `minlen=N` - minimum string length
- `maxlen=N` - maximum string length  
- `min=N` - minimum numeric value
- `max=N` - maximum numeric value
- `email` - valid email format
- `pattern=regex` - regex pattern matching

**Custom rules:** Easily added via validator registration

### Error Response Format

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

## Testing Strategy

### Unit Tests
- [ ] Test each converter with edge cases: empty strings, overflow values, invalid formats
- [ ] Validate all built-in validators with boundary conditions and error cases
- [ ] Test error message generation includes correct field names and constraint values
- [ ] Verify tag parsing logic handles malformed tags gracefully
- [ ] Test converter and validator registration/unregistration scenarios
- [ ] Verify thread safety of registry operations under concurrent access

### Integration Tests  
- [ ] Test full parsing pipeline with nested structs and complex validation chains
- [ ] Verify HTTP handler integration preserves existing error handling patterns
- [ ] Test CLI argument parsing with complex flag combinations and edge cases
- [ ] Validate JSON API integration with various Content-Type headers
- [ ] Test cross-platform CLI behavior (different shell environments)
- [ ] Verify parser instances with different configurations work independently



## Rollback Plan

If issues arise during implementation:

1. **Phase 1-2**: Easy rollback, new system is additive
2. **Phase 3-5**: Restore `core/form.go` from git and revert handler changes

**Risk Mitigation:**
- Comprehensive test coverage before integrating with handlers
- Test new system with existing command/query structs before migration

## Success Metrics

### Functionality
- [ ] All existing form validation works identically
- [ ] New input sources (CLI, JSON) work correctly
- [ ] Custom validation rules integrate seamlessly
- [ ] Error messages are more helpful than before

### Developer Experience
- [ ] Less boilerplate code in handlers
- [ ] Easier to add new field types
- [ ] Validation rules are more declarative
- [ ] Better error debugging information

## Timeline Summary

| Phase | Duration | Key Deliverables |
|-------|----------|------------------|
| 1 | Week 1 | Core parsing infrastructure |
| 2 | Week 2 | CLI support, extensibility examples |
| 3 | Week 3 | HTTP/API handler integration |
| 4 | Week 4 | Advanced features, documentation |
| 5 | Week 5 | Migration completion, cleanup |

**Total Estimated Duration: 5 weeks**

## Future Enhancements

After successful implementation, consider:

1. **GraphQL Integration** - Parse GraphQL inputs using same system
2. **Schema Generation** - Auto-generate OpenAPI schemas from validation tags
3. **IDE Support** - Language server integration for validation tag completion
4. **Async Validation** - Database/API validation rules
5. **Localization** - Multi-language error messages