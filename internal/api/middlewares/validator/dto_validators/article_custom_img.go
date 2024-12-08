package dto_validators

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
)

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

	// Check if the string starts with any allowed prefix
	hasValidPrefix := false
	for _, prefix := range allowedPrefixes {
		if len(str) > len(prefix) && str[:len(prefix)] == prefix {
			hasValidPrefix = true
			str = str[len(prefix):] // Remove the prefix
			break
		}
	}

	// If no valid prefix was found, the string is invalid
	if !hasValidPrefix {
		return false
	}

	// Try to decode the Base64 string
	_, err := base64.StdEncoding.DecodeString(str)
	return err == nil
}
