package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Global validation handler instance
var globalHandler *ValidationHandler

// NewHandler creates a new validation handler
func NewHandler(config Config) *ValidationHandler {
	handler := &ValidationHandler{
		config: config,
		rules:  make(map[string]Rule),
	}

	// Register built-in rules
	handler.registerBuiltInRules()

	// Register custom rules
	for name, rule := range config.CustomRules {
		handler.rules[name] = rule
	}

	return handler
}

// ValidateStruct validates a struct using reflection
func (h *ValidationHandler) ValidateStruct(obj interface{}, context map[string]interface{}) error {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validation: expected struct, got %s", v.Kind())
	}

	var errors ValidationErrors
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get validation tag
		tag := getValidationTag(field)
		if tag == "" || tag == "-" {
			continue
		}

		// Validate field
		if errs := h.validateField(field, value, tag, context); len(errs) > 0 {
			errors.Errors = append(errors.Errors, errs...)

			// Stop on first error if configured
			if h.config.StopOnFirstError {
				break
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}

	return nil
}

// validateField validates a single struct field
func (h *ValidationHandler) validateField(field reflect.StructField, value reflect.Value, tag string, context map[string]interface{}) []ValidationError {
	var errors []ValidationError
	fieldName := getFieldName(field)
	customMessage := getCustomMessage(field)

	// Parse validation rules
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		// Parse rule and parameters
		parts := strings.SplitN(rule, "=", 2)
		ruleName := strings.TrimSpace(parts[0])
		params := ""
		if len(parts) > 1 {
			params = strings.TrimSpace(parts[1])
		}

		// Validate using rule
		if err := h.validateRule(fieldName, value.Interface(), ruleName, params, context); err != nil {
			message := err.Error()

			// Use custom message if provided
			if customMessage != "" {
				message = h.formatMessage(customMessage, params, value.Interface())
			} else if h.config.CustomMessages {
				message = h.getLocalizedMessage(ruleName, fieldName, params, value.Interface())
			}

			errors = append(errors, ValidationError{
				Field:   fieldName,
				Value:   value.Interface(),
				Tag:     ruleName,
				Message: message,
			})

			// Stop on first error if configured
			if h.config.StopOnFirstError {
				break
			}
		}
	}

	return errors
}

// validateRule validates using a specific rule
func (h *ValidationHandler) validateRule(field string, value interface{}, ruleName, params string, context map[string]interface{}) error {
	// Check if rule exists
	if rule, exists := h.rules[ruleName]; exists {
		return rule.Validate(value, params, context)
	}

	// Built-in rules
	switch ruleName {
	case "required":
		return h.validateRequired(value)
	case "min":
		return h.validateMin(value, params)
	case "max":
		return h.validateMax(value, params)
	case "email":
		return h.validateEmail(value)
	case "oneof":
		return h.validateOneOf(value, params)
	case "numeric":
		return h.validateNumeric(value)
	case "alpha":
		return h.validateAlpha(value)
	case "alphanum":
		return h.validateAlphaNum(value)
	default:
		h.config.Logger.Warn("validation: unknown rule", "rule", ruleName, "field", field)
		return nil // Don't fail on unknown rules
	}
}

// registerBuiltInRules registers all built-in validation rules
func (h *ValidationHandler) registerBuiltInRules() {
	// Built-in rules are implemented as methods
	// Custom rules can be added via CustomRules in config
}

// formatMessage formats custom message with parameters
func (h *ValidationHandler) formatMessage(message, params string, value interface{}) string {
	message = strings.ReplaceAll(message, "{value}", fmt.Sprintf("%v", value))

	// Handle parameters like min=5, max=10
	if strings.Contains(params, "=") {
		parts := strings.SplitN(params, "=", 2)
		if len(parts) == 2 {
			key := "{" + parts[0] + "}"
			message = strings.ReplaceAll(message, key, parts[1])
		}
	} else if params != "" {
		message = strings.ReplaceAll(message, "{param}", params)
	}

	return message
}

// getLocalizedMessage gets localized error message
func (h *ValidationHandler) getLocalizedMessage(rule, field, params string, value interface{}) string {
	messages := getDefaultMessages(h.config.Locale)

	if msg, exists := messages[rule]; exists {
		msg = strings.ReplaceAll(msg, "{field}", field)
		msg = strings.ReplaceAll(msg, "{value}", fmt.Sprintf("%v", value))
		msg = h.formatMessage(msg, params, value)
		return msg
	}

	return fmt.Sprintf("validation failed for field %s", field)
}

// Built-in validation methods
func (h *ValidationHandler) validateRequired(value interface{}) error {
	if value == nil {
		return fmt.Errorf("field is required")
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if strings.TrimSpace(v.String()) == "" {
			return fmt.Errorf("field is required")
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return fmt.Errorf("field is required")
		}
	case reflect.Ptr:
		if v.IsNil() {
			return fmt.Errorf("field is required")
		}
	}

	return nil
}

func (h *ValidationHandler) validateMin(value interface{}, params string) error {
	min, err := strconv.Atoi(params)
	if err != nil {
		return fmt.Errorf("invalid min parameter: %s", params)
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if len(v.String()) < min {
			return fmt.Errorf("field must be at least %d characters", min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < int64(min) {
			return fmt.Errorf("field must be at least %d", min)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() < uint64(min) {
			return fmt.Errorf("field must be at least %d", min)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() < float64(min) {
			return fmt.Errorf("field must be at least %d", min)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() < min {
			return fmt.Errorf("field must contain at least %d items", min)
		}
	}

	return nil
}

func (h *ValidationHandler) validateMax(value interface{}, params string) error {
	max, err := strconv.Atoi(params)
	if err != nil {
		return fmt.Errorf("invalid max parameter: %s", params)
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if len(v.String()) > max {
			return fmt.Errorf("field must be at most %d characters", max)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() > int64(max) {
			return fmt.Errorf("field must be at most %d", max)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() > uint64(max) {
			return fmt.Errorf("field must be at most %d", max)
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() > float64(max) {
			return fmt.Errorf("field must be at most %d", max)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() > max {
			return fmt.Errorf("field must contain at most %d items", max)
		}
	}

	return nil
}

func (h *ValidationHandler) validateEmail(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field must be a string")
	}

	// Simple email validation
	if !strings.Contains(str, "@") || !strings.Contains(str, ".") {
		return fmt.Errorf("field must be a valid email address")
	}

	return nil
}

func (h *ValidationHandler) validateOneOf(value interface{}, params string) error {
	str := fmt.Sprintf("%v", value)
	options := strings.Split(params, " ")

	for _, option := range options {
		if strings.TrimSpace(option) == str {
			return nil
		}
	}

	return fmt.Errorf("field must be one of: %s", params)
}

func (h *ValidationHandler) validateNumeric(value interface{}) error {
	str := fmt.Sprintf("%v", value)
	if _, err := strconv.ParseFloat(str, 64); err != nil {
		return fmt.Errorf("field must be numeric")
	}
	return nil
}

func (h *ValidationHandler) validateAlpha(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field must be a string")
	}

	for _, r := range str {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return fmt.Errorf("field must contain only letters")
		}
	}

	return nil
}

func (h *ValidationHandler) validateAlphaNum(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field must be a string")
	}

	for _, r := range str {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return fmt.Errorf("field must contain only letters and numbers")
		}
	}

	return nil
}
