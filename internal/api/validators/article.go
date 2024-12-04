package validators

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateImg(fl validator.FieldLevel) bool {
	// Извлекаем значение поля
	field := fl.Field().Interface()

	// Обрабатываем случай, если это указатель на строку (*string)
	if imgPtr, ok := field.(*string); ok {
		if imgPtr == nil {
			fmt.Println("Img field is nil - validation passed")
			return true
		}
		return validateBase64(*imgPtr)
	}

	// Обрабатываем случай, если это строка (string)
	if img, ok := field.(string); ok {
		if img == "" {
			fmt.Println("Img field is empty - validation passed")
			return true
		}
		return validateBase64(img)
	}

	fmt.Println("Img field is not of type string or *string")
	return false
}

// validateBase64 проверяет, является ли строка корректной Base64 или имеет префикс data:image
func validateBase64(img string) bool {
	// Проверяем, содержит ли строка префикс data:image/
	if strings.HasPrefix(img, "data:image/") {
		parts := strings.SplitN(img, ",", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid data URL format")
			return false
		}
		img = parts[1] // Используем только Base64 часть
	}

	// Декодируем Base64
	decoded, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		fmt.Println("Base64 decode error:", err)
		return false
	}

	if len(decoded) == 0 {
		fmt.Println("Decoded Base64 is empty")
		return false
	}

	fmt.Println("Base64 validation passed")
	return true
}

func RegArticleValidators(validate *validator.Validate) {
	// Register the custom "img_base64_or_null" rule.
	err := validate.RegisterValidation("img_base64_or_null", ValidateImg)
	if err != nil {
		panic(fmt.Sprintf("Failed to register custom validator: %v", err))
	}
}
