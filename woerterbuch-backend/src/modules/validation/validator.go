package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// CustomValidator is an object to handle request validation.
type CustomValidator struct{}

// Validate validates object based on the 'validate' tag, obj must be a structure.
func (cv *CustomValidator) Validate(obj any) error {
	// Validation with reflect
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		return NewValidationError("request is not a struct")
	}

	fieldsNum := reflect.ValueOf(obj).NumField()
	structName := reflect.TypeOf(obj).Name()
	for i := 0; i < fieldsNum; i++ {
		field := reflect.ValueOf(obj).Type().Field(i)

		fieldValidateRules := strings.Split(field.Tag.Get("validate"), ",")
		if len(fieldValidateRules) == 0 {
			continue
		}

		fieldName := field.Name
		fieldValue := reflect.ValueOf(obj).FieldByName(fieldName).Interface()

		for _, ruleName := range fieldValidateRules {
			if ruleName == "" {
				continue
			}
			if valErr := cv.validationHandler(ruleName, fieldValue, structName, fieldName); valErr != nil {
				return valErr
			}
		}
	}

	return nil
}

// validationHandler handles checking of all the validation rules.
func (cv *CustomValidator) validationHandler(ruleName string, value any, structName, fieldName string) error {
	switch {
	case ruleName == "Required":
		if !cv.checkRequiredRule(value) {
			return NewValidationError(fmt.Sprintf("field '%s' of struct '%s' requires a value", fieldName, structName))
		}
	case strings.HasPrefix(ruleName, "MaxLength"):
		maxLength, err := strconv.Atoi(ruleName[len("MaxLength(") : len(ruleName)-1])
		if err != nil {
			return NewValidationError(fmt.Sprintf("MaxLength parameter of field '%s' in struct '%s' is not a number", fieldName, structName))
		}

		return cv.checkMaxLengthRule(value, maxLength, fieldName, structName)
	default:
		return NewValidationError(fmt.Sprintf("unknown validation rule '%s' for field '%s' of struct '%s'", ruleName, fieldName, structName))
	}

	return nil
}

// checkRequiredRule validates rule `validate:"Required"`.
func (cv *CustomValidator) checkRequiredRule(obj any) bool {
	if obj == nil {
		return false
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.Len() > 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.Interface:
		return !v.IsNil()
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !cv.checkRequiredRule(v.Field(i).Interface()) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// checkMaxLengthRule validates rule `validate:"MaxLength(number)"` for strings.
func (cv *CustomValidator) checkMaxLengthRule(value any, maxLength int, fieldName, structName string) error {
	strValue, ok := value.(string)
	if !ok {
		return NewValidationError(fmt.Sprintf("value of field '%s' in struct '%s' is not a string, but MaxLength rule was declared", fieldName, structName))
	}

	strLength := len(strValue)
	if strLength > maxLength {
		return NewValidationError(fmt.Sprintf("field '%s' of struct '%s' has a length of %d, but maximum length is %d", fieldName, structName, strLength, maxLength))
	}
	return nil
}
