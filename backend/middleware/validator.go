package middleware

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct using `validate` tags.
// Supported tags: required, min, max, gte, lte, gt, email, oneof
// Tags are comma-separated; each tag is key=value or just key.
func ValidateStruct(s interface{}) []ValidationError {
	var errs []ValidationError

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return errs
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" || tag == "-" {
			continue
		}

		fieldName := field.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}
		// Strip json options like ",omitempty"
		if idx := strings.Index(fieldName, ","); idx != -1 {
			fieldName = fieldName[:idx]
		}

		for _, rule := range strings.Split(tag, ",") {
			if err := validateRule(fieldName, value, rule); err != nil {
				errs = append(errs, *err)
				break // one error per field
			}
		}
	}

	return errs
}

func validateRule(fieldName string, value reflect.Value, rule string) *ValidationError {
	parts := strings.SplitN(rule, "=", 2)
	ruleName := parts[0]
	ruleParam := ""
	if len(parts) > 1 {
		ruleParam = parts[1]
	}

	switch ruleName {
	case "required":
		if isZero(value) {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s is required", fieldName)}
		}

	case "min":
		n, _ := strconv.Atoi(ruleParam)
		if value.Kind() == reflect.String && len(value.String()) < n {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be at least %d characters", fieldName, n)}
		}

	case "max":
		n, _ := strconv.Atoi(ruleParam)
		if value.Kind() == reflect.String && len(value.String()) > n {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be at most %d characters", fieldName, n)}
		}

	case "gte":
		limit, _ := strconv.ParseFloat(ruleParam, 64)
		if numVal, ok := toFloat64(value); ok && numVal < limit {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be at least %s", fieldName, ruleParam)}
		}

	case "lte":
		limit, _ := strconv.ParseFloat(ruleParam, 64)
		if numVal, ok := toFloat64(value); ok && numVal > limit {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be at most %s", fieldName, ruleParam)}
		}

	case "gt":
		limit, _ := strconv.ParseFloat(ruleParam, 64)
		if numVal, ok := toFloat64(value); ok && numVal <= limit {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be greater than %s", fieldName, ruleParam)}
		}

	case "email":
		s := value.String()
		if s != "" && (!strings.Contains(s, "@") || !strings.Contains(s, ".")) {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be a valid email address", fieldName)}
		}

	case "oneof":
		s := fmt.Sprint(value.Interface())
		allowed := strings.Fields(ruleParam)
		found := false
		for _, a := range allowed {
			if s == a {
				found = true
				break
			}
		}
		if !found {
			return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be one of: %s", fieldName, ruleParam)}
		}
	}

	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Map, reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func toFloat64(v reflect.Value) (float64, bool) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}
