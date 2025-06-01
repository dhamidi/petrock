package core

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// FormSource represents different input sources (unchanged)
type FormSource interface {
	Get(key string) string
	GetAll(key string) []string
	Keys() []string
}

// URLValuesSource, MapSource implementations (same as before)
type URLValuesSource struct{ Values url.Values }
func (u URLValuesSource) Get(key string) string { return u.Values.Get(key) }
func (u URLValuesSource) GetAll(key string) []string { return u.Values[key] }
func (u URLValuesSource) Keys() []string {
	keys := make([]string, 0, len(u.Values))
	for k := range u.Values { keys = append(keys, k) }
	return keys
}

type MapSource struct{ Data map[string]interface{} }
func (m MapSource) Get(key string) string {
	if val, exists := m.Data[key]; exists {
		switch v := val.(type) {
		case string: return v
		case []string: if len(v) > 0 { return v[0] }
		default: return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
func (m MapSource) GetAll(key string) []string {
	if val, exists := m.Data[key]; exists {
		switch v := val.(type) {
		case string: return []string{v}
		case []string: return v
		case []interface{}:
			result := make([]string, len(v))
			for i, item := range v { result[i] = fmt.Sprintf("%v", item) }
			return result
		default: return []string{fmt.Sprintf("%v", v)}
		}
	}
	return nil
}
func (m MapSource) Keys() []string {
	keys := make([]string, 0, len(m.Data))
	for k := range m.Data { keys = append(keys, k) }
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
	return fmt.Sprintf("%d validation errors", len(e.Errors))
}

// FieldContext provides context about a field being processed
type FieldContext struct {
	Name       string      // Field name
	Value      interface{} // Current value
	FieldType  reflect.Type
	StructField reflect.StructField
	Tags       map[string]string // Parsed struct tags
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

// Built-in converters

// BasicConverter handles primitive types
type BasicConverter struct{}

func (c BasicConverter) CanConvert(targetType reflect.Type) bool {
	switch targetType.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		 reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	}
	return false
}

func (c BasicConverter) Convert(value string, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.String:
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			return reflect.Zero(targetType).Interface(), nil
		}
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer: %s", value)
		}
		switch targetType.Kind() {
		case reflect.Int: return int(intVal), nil
		case reflect.Int8: return int8(intVal), nil
		case reflect.Int16: return int16(intVal), nil
		case reflect.Int32: return int32(intVal), nil
		case reflect.Int64: return intVal, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			return reflect.Zero(targetType).Interface(), nil
		}
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid unsigned integer: %s", value)
		}
		switch targetType.Kind() {
		case reflect.Uint: return uint(uintVal), nil
		case reflect.Uint8: return uint8(uintVal), nil
		case reflect.Uint16: return uint16(uintVal), nil
		case reflect.Uint32: return uint32(uintVal), nil
		case reflect.Uint64: return uintVal, nil
		}
	case reflect.Float32, reflect.Float64:
		if value == "" {
			return reflect.Zero(targetType).Interface(), nil
		}
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float: %s", value)
		}
		if targetType.Kind() == reflect.Float32 {
			return float32(floatVal), nil
		}
		return floatVal, nil
	case reflect.Bool:
		if value == "" {
			return false, nil
		}
		if value == "1" {
			return true, nil
		} else if value == "0" {
			return false, nil
		}
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean: %s", value)
		}
		return boolVal, nil
	}
	return nil, fmt.Errorf("conversion not supported")
}

func (c BasicConverter) ConvertSlice(values []string, targetType reflect.Type) (interface{}, error) {
	elemType := targetType.Elem()
	slice := reflect.MakeSlice(targetType, len(values), len(values))
	
	for i, value := range values {
		converted, err := c.Convert(value, elemType)
		if err != nil {
			return nil, fmt.Errorf("error converting element %d: %w", i, err)
		}
		slice.Index(i).Set(reflect.ValueOf(converted))
	}
	
	return slice.Interface(), nil
}

// TimeConverter handles time.Time
type TimeConverter struct {
	Formats []string // Configurable time formats
}

func NewTimeConverter() *TimeConverter {
	return &TimeConverter{
		Formats: []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02 15:04:05",
			"01/02/2006",
			"01/02/2006 15:04:05",
		},
	}
}

func (c TimeConverter) CanConvert(targetType reflect.Type) bool {
	return targetType == reflect.TypeOf(time.Time{})
}

func (c TimeConverter) Convert(value string, targetType reflect.Type) (interface{}, error) {
	if value == "" {
		return time.Time{}, nil
	}
	
	for _, format := range c.Formats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}
	
	return nil, fmt.Errorf("invalid time format: %s", value)
}

func (c TimeConverter) ConvertSlice(values []string, targetType reflect.Type) (interface{}, error) {
	slice := make([]time.Time, len(values))
	for i, value := range values {
		converted, err := c.Convert(value, targetType.Elem())
		if err != nil {
			return nil, fmt.Errorf("error converting time element %d: %w", i, err)
		}
		slice[i] = converted.(time.Time)
	}
	return slice, nil
}

// Built-in validators

// RequiredValidator checks for required fields
type RequiredValidator struct{}

func (v RequiredValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.GetTagBool("required")
}

func (v RequiredValidator) Validate(ctx *FieldContext) []ParseError {
	switch val := ctx.Value.(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			return []ParseError{{
				Field: ctx.Name, Message: "This field is required", Code: "required",
			}}
		}
	case nil:
		return []ParseError{{
			Field: ctx.Name, Message: "This field is required", Code: "required",
		}}
	default:
		if reflect.ValueOf(val).IsZero() {
			return []ParseError{{
				Field: ctx.Name, Message: "This field is required", Code: "required",
			}}
		}
	}
	return nil
}

// LengthValidator checks string length constraints
type LengthValidator struct{}

func (v LengthValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.FieldType.Kind() == reflect.String && 
		   (ctx.GetTag("minlen", "") != "" || ctx.GetTag("maxlen", "") != "")
}

func (v LengthValidator) Validate(ctx *FieldContext) []ParseError {
	str, ok := ctx.Value.(string)
	if !ok || str == "" {
		return nil // Skip validation for non-strings or empty strings
	}
	
	var errors []ParseError
	length := utf8.RuneCountInString(str)
	
	if minLen := ctx.GetTagInt("minlen", 0); minLen > 0 && length < minLen {
		errors = append(errors, ParseError{
			Field: ctx.Name,
			Message: fmt.Sprintf("Must be at least %d characters long", minLen),
			Code: "min_length",
			Meta: map[string]interface{}{"min_length": minLen, "actual_length": length},
		})
	}
	
	if maxLen := ctx.GetTagInt("maxlen", 0); maxLen > 0 && length > maxLen {
		errors = append(errors, ParseError{
			Field: ctx.Name,
			Message: fmt.Sprintf("Must be no more than %d characters long", maxLen),
			Code: "max_length",
			Meta: map[string]interface{}{"max_length": maxLen, "actual_length": length},
		})
	}
	
	return errors
}

// RangeValidator checks numeric range constraints
type RangeValidator struct{}

func (v RangeValidator) CanValidate(ctx *FieldContext) bool {
	kind := ctx.FieldType.Kind()
	isNumeric := kind >= reflect.Int && kind <= reflect.Float64
	return isNumeric && (ctx.GetTag("min", "") != "" || ctx.GetTag("max", "") != "")
}

func (v RangeValidator) Validate(ctx *FieldContext) []ParseError {
	var numVal int64
	var isValid bool
	
	switch val := ctx.Value.(type) {
	case int: numVal, isValid = int64(val), true
	case int8: numVal, isValid = int64(val), true
	case int16: numVal, isValid = int64(val), true
	case int32: numVal, isValid = int64(val), true
	case int64: numVal, isValid = val, true
	case uint: numVal, isValid = int64(val), true
	case uint8: numVal, isValid = int64(val), true
	case uint16: numVal, isValid = int64(val), true
	case uint32: numVal, isValid = int64(val), true
	case uint64: numVal, isValid = int64(val), true
	}
	
	if !isValid {
		return nil
	}
	
	var errors []ParseError
	
	if minTag := ctx.GetTag("min", ""); minTag != "" {
		if min, err := strconv.ParseInt(minTag, 10, 64); err == nil && numVal < min {
			errors = append(errors, ParseError{
				Field: ctx.Name,
				Message: fmt.Sprintf("Must be at least %d", min),
				Code: "min_value",
				Meta: map[string]interface{}{"min_value": min, "actual_value": numVal},
			})
		}
	}
	
	if maxTag := ctx.GetTag("max", ""); maxTag != "" {
		if max, err := strconv.ParseInt(maxTag, 10, 64); err == nil && numVal > max {
			errors = append(errors, ParseError{
				Field: ctx.Name,
				Message: fmt.Sprintf("Must be no more than %d", max),
				Code: "max_value",
				Meta: map[string]interface{}{"max_value": max, "actual_value": numVal},
			})
		}
	}
	
	return errors
}

// EmailValidator validates email format
type EmailValidator struct{}

func (v EmailValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.FieldType.Kind() == reflect.String && ctx.GetTagBool("email")
}

func (v EmailValidator) Validate(ctx *FieldContext) []ParseError {
	str, ok := ctx.Value.(string)
	if !ok || str == "" {
		return nil
	}
	
	if _, err := mail.ParseAddress(str); err != nil {
		return []ParseError{{
			Field: ctx.Name,
			Message: "Must be a valid email address",
			Code: "invalid_email",
		}}
	}
	
	return nil
}

// CustomValidator allows for custom validation functions
type CustomValidator struct {
	ValidateFunc func(ctx *FieldContext) []ParseError
	CanValidateFunc func(ctx *FieldContext) bool
}

func (v CustomValidator) CanValidate(ctx *FieldContext) bool {
	if v.CanValidateFunc != nil {
		return v.CanValidateFunc(ctx)
	}
	return true // Default to always applicable
}

func (v CustomValidator) Validate(ctx *FieldContext) []ParseError {
	return v.ValidateFunc(ctx)
}

// StandardTagParser parses common validation tags
type StandardTagParser struct{}

func (p StandardTagParser) ParseTags(field reflect.StructField) map[string]string {
	tags := make(map[string]string)
	
	// Parse validate tag: "required,minlen=2,maxlen=100,email"
	if validateTag := field.Tag.Get("validate"); validateTag != "" {
		parts := strings.Split(validateTag, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				tags[kv[0]] = kv[1]
			} else {
				tags[part] = "true"
			}
		}
	}
	
	return tags
}

// Parser is the main parsing engine that coordinates all components
type Parser struct {
	converters    *ConverterRegistry
	validators    *ValidatorRegistry
	tagParsers    []TagParser
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
