package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// Indonesian-specific validation rules

// NIKRule validates Indonesian NIK (Nomor Induk Kependudukan)
type NIKRule struct{}

func (r *NIKRule) Validate(value interface{}, params string, context map[string]interface{}) error {
	nik, ok := value.(string)
	if !ok {
		return fmt.Errorf("NIK must be a string")
	}

	// NIK must be 16 digits
	if len(nik) != 16 {
		return fmt.Errorf("NIK must be 16 digits")
	}

	// Check if all characters are digits
	for _, char := range nik {
		if char < '0' || char > '9' {
			return fmt.Errorf("NIK must contain only digits")
		}
	}

	// Additional NIK validation logic can be added here
	// (birth date validation, region code validation, etc.)

	return nil
}

func (r *NIKRule) GetMessage() string {
	return "field must be a valid NIK"
}

// IndonesianPhoneRule validates Indonesian phone numbers
type IndonesianPhoneRule struct{}

func (r *IndonesianPhoneRule) Validate(value interface{}, params string, context map[string]interface{}) error {
	phone, ok := value.(string)
	if !ok {
		return fmt.Errorf("phone number must be a string")
	}

	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Indonesian phone patterns
	patterns := []string{
		`^(\+62|62|0)8[1-9][0-9]{6,9}$`, // Mobile numbers
		`^(\+62|62|0)[2-9][0-9]{7,11}$`, // Landline numbers
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, phone); matched {
			return nil
		}
	}

	return fmt.Errorf("phone number must be a valid Indonesian phone number")
}

func (r *IndonesianPhoneRule) GetMessage() string {
	return "field must be a valid Indonesian phone number"
}

// URLRule validates URLs
type URLRule struct{}

func (r *URLRule) Validate(value interface{}, params string, context map[string]interface{}) error {
	url, ok := value.(string)
	if !ok {
		return fmt.Errorf("URL must be a string")
	}

	urlPattern := `^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`
	matched, _ := regexp.MatchString(urlPattern, url)
	if !matched {
		return fmt.Errorf("URL must be a valid URL")
	}

	return nil
}

func (r *URLRule) GetMessage() string {
	return "field must be a valid URL"
}

// UUIDRule validates UUID v4
type UUIDRule struct{}

func (r *UUIDRule) Validate(value interface{}, params string, context map[string]interface{}) error {
	uuid, ok := value.(string)
	if !ok {
		return fmt.Errorf("UUID must be a string")
	}

	uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(uuidPattern, strings.ToLower(uuid))
	if !matched {
		return fmt.Errorf("UUID must be a valid UUID v4")
	}

	return nil
}

func (r *UUIDRule) GetMessage() string {
	return "field must be a valid UUID"
}

// PasswordRule validates password strength
type PasswordRule struct{}

func (r *PasswordRule) Validate(value interface{}, params string, context map[string]interface{}) error {
	password, ok := value.(string)
	if !ok {
		return fmt.Errorf("password must be a string")
	}

	// Parse parameters for password requirements
	requirements := strings.Split(params, ",")

	for _, req := range requirements {
		req = strings.TrimSpace(req)
		switch req {
		case "upper":
			if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
				return fmt.Errorf("password must contain at least one uppercase letter")
			}
		case "lower":
			if !regexp.MustCompile(`[a-z]`).MatchString(password) {
				return fmt.Errorf("password must contain at least one lowercase letter")
			}
		case "number":
			if !regexp.MustCompile(`[0-9]`).MatchString(password) {
				return fmt.Errorf("password must contain at least one number")
			}
		case "special":
			if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
				return fmt.Errorf("password must contain at least one special character")
			}
		}
	}

	return nil
}

func (r *PasswordRule) GetMessage() string {
	return "password does not meet requirements"
}

// Helper functions to create rule instances
func NewNIKRule() Rule {
	return &NIKRule{}
}

func NewIndonesianPhoneRule() Rule {
	return &IndonesianPhoneRule{}
}

func NewURLRule() Rule {
	return &URLRule{}
}

func NewUUIDRule() Rule {
	return &UUIDRule{}
}

func NewPasswordRule() Rule {
	return &PasswordRule{}
}

// Built-in rules map for easy registration
func GetBuiltInRules() map[string]Rule {
	return map[string]Rule{
		"nik":      NewNIKRule(),
		"phone_id": NewIndonesianPhoneRule(),
		"url":      NewURLRule(),
		"uuid":     NewUUIDRule(),
		"password": NewPasswordRule(),
	}
}
