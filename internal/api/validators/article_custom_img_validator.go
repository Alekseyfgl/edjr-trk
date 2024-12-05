// internal/api/validators/article_validators.go
package validators

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
)

// RegArticleValidators registers custom validators for articles.
func RegArticleValidators(validate *validator.Validate) {
	// Register the custom validator 'img_base64_or_null'.
	validate.RegisterValidation("img_base64_or_null", imgBase64OrNull)
}

// imgBase64OrNull checks if a string is either empty, null, or a valid Base64-encoded string.
func imgBase64OrNull(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()

	// Allow null values
	if value == nil {
		return true
	}

	// Check if the value is a string
	str, ok := value.(string)
	if !ok {
		return false
	}

	// Empty string is considered valid
	if str == "" {
		return true
	}

	// Try to decode the Base64 string
	_, err := base64.StdEncoding.DecodeString(str)
	return err == nil
}
