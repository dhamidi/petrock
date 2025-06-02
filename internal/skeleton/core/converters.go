package core

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

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
		case reflect.Int:
			return int(intVal), nil
		case reflect.Int8:
			return int8(intVal), nil
		case reflect.Int16:
			return int16(intVal), nil
		case reflect.Int32:
			return int32(intVal), nil
		case reflect.Int64:
			return intVal, nil
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
		case reflect.Uint:
			return uint(uintVal), nil
		case reflect.Uint8:
			return uint8(uintVal), nil
		case reflect.Uint16:
			return uint16(uintVal), nil
		case reflect.Uint32:
			return uint32(uintVal), nil
		case reflect.Uint64:
			return uintVal, nil
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
