package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Custom (User) error messages
var customMessages = map[string]string{
	"Username.username": "Username can only contain Latin characters, no spaces, no special characters",
	"Username.required": "Username is required",
	"Username.min":      "Username must be at least 8 characters long",
	"Username.max":      "Username must be at most 32 characters long",
	"Email.required":    "Email is required",
	"Email.email":       "Email must be a valid email address",
	"Password.required": "Password is required",
	"Password.min":      "Password must be at least 8 characters long",
}

func usernameValidator(fl validator.FieldLevel) bool {
	// Regular expression to allow only Latin characters (a-zA-Z), no spaces, no special characters
	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	username := fl.Field().String()

	// Check if the username matches the regex
	if !regex.MatchString(username) {
		return false
	}

	return true
}

func ValidateUser(data interface{}) error {
	validate = validator.New()

	// Registering the custom validation function
	validate.RegisterValidation("username", usernameValidator)

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.StructField()
			tag := err.Tag()
			customErrorKey := fmt.Sprintf("%s.%s", field, tag)

			// If there is a custom message for the validation error
			if customMsg, exists := customMessages[customErrorKey]; exists {
				return fmt.Errorf(customMsg)
			}

			// Fallback to default error message if no custom message is defined
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}
