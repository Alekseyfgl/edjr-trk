package dto_validator

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

	// Define allowed prefixes for Base64 data URIs
	allowedPrefixes := []string{
		"data:image/jpeg;base64,", // JPEG
		"data:image/png;base64,",  // PNG
		"data:image/gif;base64,",  // GIF
		"data:image/webp;base64,", // WebP
	}

	// Remove any allowed prefix
	for _, prefix := range allowedPrefixes {
		if len(str) > len(prefix) && str[:len(prefix)] == prefix {
			str = str[len(prefix):] // Remove the prefix
			break
		}
	}

	// Try to decode the Base64 string
	_, err := base64.StdEncoding.DecodeString(str)
	return err == nil
}
