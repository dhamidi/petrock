package core

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// FormSource represents different input sources (HTTP forms, JSON, CLI args, etc.)
type FormSource interface {
	Get(key string) string
	GetAll(key string) []string
	Keys() []string
}

// URLValuesSource adapts url.Values to FormSource interface
type URLValuesSource struct {
	Values url.Values
}

func (u URLValuesSource) Get(key string) string {
	return u.Values.Get(key)
}

func (u URLValuesSource) GetAll(key string) []string {
	return u.Values[key]
}

func (u URLValuesSource) Keys() []string {
	keys := make([]string, 0, len(u.Values))
	for k := range u.Values {
		keys = append(keys, k)
	}
	return keys
}

// MapSource adapts map[string]interface{} to FormSource interface for JSON/API data
type MapSource struct {
	Data map[string]interface{}
}

func (m MapSource) Get(key string) string {
	if val, exists := m.Data[key]; exists {
		switch v := val.(type) {
		case string:
			return v
		case []string:
			if len(v) > 0 {
				return v[0]
			}
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func (m MapSource) GetAll(key string) []string {
	if val, exists := m.Data[key]; exists {
		switch v := val.(type) {
		case string:
			return []string{v}
		case []string:
			return v
		case []interface{}:
			result := make([]string, len(v))
			for i, item := range v {
				result[i] = fmt.Sprintf("%v", item)
			}
			return result
		default:
			return []string{fmt.Sprintf("%v", v)}
		}
	}
	return nil
}

func (m MapSource) Keys() []string {
	keys := make([]string, 0, len(m.Data))
	for k := range m.Data {
		keys = append(keys, k)
	}
	return keys
}

// ParseError represents a single parsing or validation error
type ParseError struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Code    string                 `json:"code,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"` // Additional context
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ParseErrors holds multiple parsing errors
type ParseErrors struct {
	Errors []ParseError `json:"errors"`
}

func (e *ParseErrors) Add(field, message, code string) {
	e.Errors = append(e.Errors, ParseError{
		Field: field, Message: message, Code: code,
	})
}

func (e *ParseErrors) AddWithMeta(field, message, code string, meta map[string]interface{}) {
	e.Errors = append(e.Errors, ParseError{
		Field: field, Message: message, Code: code, Meta: meta,
	})
}

// Helper methods for common error patterns
func (e *ParseErrors) AddRequired(field string) {
	e.Add(field, "This field is required", "required")
}

func (e *ParseErrors) AddInvalidType(field, expectedType, actualValue string) {
	e.AddWithMeta(field, 
		fmt.Sprintf("Invalid type, expected %s", expectedType), 
		"invalid_type",
		map[string]interface{}{
			"expected_type": expectedType,
			"actual_value": actualValue,
		})
}

func (e *ParseErrors) AddOutOfRange(field string, min, max, actual int64) {
	var message string
	var code string
	var meta map[string]interface{}

	if actual < min {
		message = fmt.Sprintf("Must be at least %d", min)
		code = "min_value"
		meta = map[string]interface{}{"min_value": min, "actual_value": actual}
	} else {
		message = fmt.Sprintf("Must be no more than %d", max) 
		code = "max_value"
		meta = map[string]interface{}{"max_value": max, "actual_value": actual}
	}

	e.AddWithMeta(field, message, code, meta)
}

func (e *ParseErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e ParseErrors) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Validation failed with %d errors:\n", len(e.Errors)))
	for i, err := range e.Errors {
		builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}
	return builder.String()
}

// FieldContext provides context about a field being processed
type FieldContext struct {
	Name        string      // Field name
	Value       interface{} // Current value
	FieldType   reflect.Type
	StructField reflect.StructField
	Tags        map[string]string // Parsed struct tags
}

// GetTag gets a tag value with default
func (ctx *FieldContext) GetTag(key, defaultValue string) string {
	if val, exists := ctx.Tags[key]; exists {
		return val
	}
	return defaultValue
}

// GetTagBool gets a boolean tag value
func (ctx *FieldContext) GetTagBool(key string) bool {
	val := ctx.GetTag(key, "false")
	return val == "true" || val == "1"
}

// GetTagInt gets an integer tag value
func (ctx *FieldContext) GetTagInt(key string, defaultValue int) int {
	val := ctx.GetTag(key, "")
	if val == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(val); err == nil {
		return parsed
	}
	return defaultValue
}

// Converter interface for type conversion
type Converter interface {
	// CanConvert returns true if this converter can handle the target type
	CanConvert(targetType reflect.Type) bool
	// Convert converts a string value to the target type
	Convert(value string, targetType reflect.Type) (interface{}, error)
	// ConvertSlice converts multiple string values to a slice of the target type
	ConvertSlice(values []string, targetType reflect.Type) (interface{}, error)
}

// Validator interface for field validation
type Validator interface {
	// CanValidate returns true if this validator can handle the field context
	CanValidate(ctx *FieldContext) bool
	// Validate performs validation and returns any errors
	Validate(ctx *FieldContext) []ParseError
}

// TagParser interface for parsing struct tags into context
type TagParser interface {
	// ParseTags extracts relevant tags from a struct field
	ParseTags(field reflect.StructField) map[string]string
}

// ConverterRegistry manages type converters
type ConverterRegistry struct {
	converters []Converter
}

func NewConverterRegistry() *ConverterRegistry {
	return &ConverterRegistry{
		converters: make([]Converter, 0),
	}
}

func (r *ConverterRegistry) Register(converter Converter) {
	r.converters = append(r.converters, converter)
}

func (r *ConverterRegistry) Convert(value string, targetType reflect.Type) (interface{}, error) {
	for _, converter := range r.converters {
		if converter.CanConvert(targetType) {
			return converter.Convert(value, targetType)
		}
	}
	return nil, fmt.Errorf("no converter found for type %s", targetType)
}

func (r *ConverterRegistry) ConvertSlice(values []string, targetType reflect.Type) (interface{}, error) {
	for _, converter := range r.converters {
		if converter.CanConvert(targetType.Elem()) {
			return converter.ConvertSlice(values, targetType)
		}
	}
	return nil, fmt.Errorf("no converter found for slice type %s", targetType)
}

// ValidatorRegistry manages validators
type ValidatorRegistry struct {
	validators []Validator
}

func NewValidatorRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		validators: make([]Validator, 0),
	}
}

func (r *ValidatorRegistry) Register(validator Validator) {
	r.validators = append(r.validators, validator)
}

func (r *ValidatorRegistry) Validate(ctx *FieldContext) []ParseError {
	var errors []ParseError
	for _, validator := range r.validators {
		if validator.CanValidate(ctx) {
			errors = append(errors, validator.Validate(ctx)...)
		}
	}
	return errors
}

// Parser is the main parsing engine that coordinates all components
type Parser struct {
	converters *ConverterRegistry
	validators *ValidatorRegistry
	tagParsers []TagParser
}

func NewParser() *Parser {
	p := &Parser{
		converters: NewConverterRegistry(),
		validators: NewValidatorRegistry(),
		tagParsers: make([]TagParser, 0),
	}

	// Register default components
	p.RegisterConverter(BasicConverter{})
	p.RegisterConverter(NewTimeConverter())

	p.RegisterValidator(RequiredValidator{})
	p.RegisterValidator(LengthValidator{})
	p.RegisterValidator(RangeValidator{})
	p.RegisterValidator(EmailValidator{})
	p.RegisterValidator(CrossFieldValidator{})
	p.RegisterValidator(ConditionalValidator{})
	p.RegisterValidator(CustomMessageValidator{})

	p.RegisterTagParser(StandardTagParser{})

	return p
}

func (p *Parser) RegisterConverter(converter Converter) {
	p.converters.Register(converter)
}

func (p *Parser) RegisterValidator(validator Validator) {
	p.validators.Register(validator)
}

func (p *Parser) RegisterTagParser(parser TagParser) {
	p.tagParsers = append(p.tagParsers, parser)
}

func (p *Parser) ParseFrom(source FormSource, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	targetValue = targetValue.Elem()
	targetType := targetValue.Type()

	var errors ParseErrors

	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Field(i)
		fieldType := targetType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get field name from JSON tag or struct field name
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "-" && parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Parse tags
		tags := make(map[string]string)
		for _, parser := range p.tagParsers {
			for k, v := range parser.ParseTags(fieldType) {
				tags[k] = v
			}
		}

		// Handle conversion
		if fieldType.Type.Kind() == reflect.Slice {
			values := source.GetAll(fieldName)
			if len(values) > 0 {
				converted, err := p.converters.ConvertSlice(values, fieldType.Type)
				if err != nil {
					errors.Add(fieldName, err.Error(), "conversion_error")
					continue
				}
				field.Set(reflect.ValueOf(converted))
			}
		} else {
			value := source.Get(fieldName)
			converted, err := p.converters.Convert(value, fieldType.Type)
			if err != nil {
				errors.Add(fieldName, err.Error(), "conversion_error")
				continue
			}
			field.Set(reflect.ValueOf(converted))
		}

		// Create context for validation
		ctx := &FieldContext{
			Name:        fieldName,
			Value:       field.Interface(),
			FieldType:   fieldType.Type,
			StructField: fieldType,
			Tags:        tags,
		}

		// Run validation
		validationErrors := p.validators.Validate(ctx)
		errors.Errors = append(errors.Errors, validationErrors...)
	}

	if errors.HasErrors() {
		return &errors
	}

	return nil
}

// Global default parser instance
var DefaultParser = NewParser()

// Convenience functions using the default parser
func ParseFromSource(source FormSource, target interface{}) error {
	return DefaultParser.ParseFrom(source, target)
}

func ParseFromURLValues(values url.Values, target interface{}) error {
	return DefaultParser.ParseFrom(URLValuesSource{Values: values}, target)
}

func ParseFromMap(data map[string]interface{}, target interface{}) error {
	return DefaultParser.ParseFrom(MapSource{Data: data}, target)
}

// RegisterConverter adds a converter to the default parser
func RegisterConverter(converter Converter) {
	DefaultParser.RegisterConverter(converter)
}

// RegisterValidator adds a validator to the default parser
func RegisterValidator(validator Validator) {
	DefaultParser.RegisterValidator(validator)
}

// RegisterTagParser adds a tag parser to the default parser
func RegisterTagParser(parser TagParser) {
	DefaultParser.RegisterTagParser(parser)
}
