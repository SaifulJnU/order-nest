package helper

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	Validator *validator.Validate
)

// InitializeValidator initializes the global validator with custom validations.
func InitializeValidator() {
	Validator = validator.New()

	// Register BD phone validation
	_ = Validator.RegisterValidation("regex_bd_phone", func(fl validator.FieldLevel) bool {
		regex := `^(01)[3-9]{1}[0-9]{8}$`
		matched, _ := regexp.MatchString(regex, fl.Field().String())
		return matched
	})
}

// JSONTagOrFieldName returns the json tag for a field, or field name when absent.
func JSONTagOrFieldName(obj interface{}, fieldName string) string {
	field, found := reflect.TypeOf(obj).FieldByName(fieldName)
	if !found {
		return fieldName // Return fieldName if JSON tag not found
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName // Return fieldName if JSON tag is empty
	}
	return jsonTag
}
