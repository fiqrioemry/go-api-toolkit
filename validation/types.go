package validation

import (
	"reflect"
	"strings"
)

// ValidationHandler handles framework-agnostic validation
type ValidationHandler struct {
	config Config
	rules  map[string]Rule
}

// Config holds validation configuration
type Config struct {
	Logger           Logger
	CustomMessages   bool
	StopOnFirstError bool
	Locale           string
	CustomRules      map[string]Rule
	ErrorMessages    map[string]string
}

// InitConfig for initialization
type InitConfig struct {
	Logger           any               `json:"-"`
	CustomMessages   bool              `json:"custom_messages"`
	StopOnFirstError bool              `json:"stop_on_first"`
	Locale           string            `json:"locale"`
	CustomRules      map[string]Rule   `json:"-"`
	ErrorMessages    map[string]string `json:"error_messages"`
}

// ValidationConfig for request-level options
type ValidationConfig struct {
	Context   map[string]any
	ForceJSON bool
	ForceForm bool
	Locale    string
	Rules     map[string]Rule
}

// ValidationOption for flexible configuration
type ValidationOption func(*ValidationConfig)

// ValidationError represents validation failure details
type ValidationError struct {
	Field   string `json:"field"`
	Value   any    `json:"value,omitempty"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// ValidationErrors collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// applyConfigDefaults applies defaults and adapts logger automatically
func applyConfigDefaults(config InitConfig) Config {
	result := getDefaultConfig()

	// 🚀 AUTO-ADAPT LOGGER (this is where the magic happens)
	result.Logger = adaptLogger(config.Logger)

	// Apply other config values if provided
	if config.Locale != "" {
		result.Locale = config.Locale
		// Update error messages for the locale
		result.ErrorMessages = getDefaultMessages(config.Locale)
	}

	if config.CustomRules != nil {
		result.CustomRules = config.CustomRules
	}

	if config.ErrorMessages != nil {
		// Merge custom messages with defaults
		for key, value := range config.ErrorMessages {
			result.ErrorMessages[key] = value
		}
	}

	// Apply boolean settings
	result.CustomMessages = config.CustomMessages
	result.StopOnFirstError = config.StopOnFirstError

	return result
}

// Rule interface for custom validation rules
type Rule interface {
	Validate(value any, params string, context map[string]any) error
	GetMessage() string
}

// FieldValidator holds field validation information
type FieldValidator struct {
	Field    reflect.StructField
	Value    reflect.Value
	Tag      string
	Rules    []string
	Messages map[string]string
	Required bool
}

// BindingType represents the type of data binding
type BindingType int

const (
	BindingAuto BindingType = iota
	BindingJSON
	BindingForm
	BindingQuery
	BindingMultipart
)

// JSONWriter interface for framework-agnostic response writing
type JSONWriter interface {
	WriteJSON(statusCode int, data any) error
	GetStatusCode() int
}

// getDefaultConfig returns default configuration
func getDefaultConfig() Config {
	return Config{
		Logger:           &NoOpLogger{},
		CustomMessages:   true,
		StopOnFirstError: false,
		Locale:           "en",
		CustomRules:      make(map[string]Rule),
		ErrorMessages:    getDefaultMessages("en"),
	}
}

// Build validation config from options
func buildValidationConfig(opts ...ValidationOption) ValidationConfig {
	config := ValidationConfig{
		Context: make(map[string]any),
		Locale:  "en",
		Rules:   make(map[string]Rule),
	}

	for _, opt := range opts {
		opt(&config)
	}

	return config
}

// Validation options
func WithContext(ctx map[string]any) ValidationOption {
	return func(config *ValidationConfig) {
		config.Context = ctx
	}
}

func ForceJSON() ValidationOption {
	return func(config *ValidationConfig) {
		config.ForceJSON = true
	}
}

func ForceForm() ValidationOption {
	return func(config *ValidationConfig) {
		config.ForceForm = true
	}
}

func ForceQuery() ValidationOption {
	return func(config *ValidationConfig) {
		config.ForceForm = true // Query params use form binding
	}
}

func WithLocale(locale string) ValidationOption {
	return func(config *ValidationConfig) {
		config.Locale = locale
	}
}

func WithCustomRules(rules map[string]Rule) ValidationOption {
	return func(config *ValidationConfig) {
		for k, v := range rules {
			config.Rules[k] = v
		}
	}
}

// Helper to extract validation tag
func getValidationTag(field reflect.StructField) string {
	tag := field.Tag.Get("validate")
	if tag == "" {
		tag = field.Tag.Get("validation")
	}
	return tag
}

// Helper to extract custom message
func getCustomMessage(field reflect.StructField) string {
	return field.Tag.Get("message")
}

// Helper to get field name for error
func getFieldName(field reflect.StructField) string {
	// Try json tag first, then form tag, then field name
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		if name := strings.Split(jsonTag, ",")[0]; name != "-" {
			return name
		}
	}
	if formTag := field.Tag.Get("form"); formTag != "" {
		if name := strings.Split(formTag, ",")[0]; name != "-" {
			return name
		}
	}
	return strings.ToLower(field.Name)
}

// Error implements error interface
func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}
	return ve.Errors[0].Message
}

// HasErrors checks if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// ToMap converts validation errors to map for JSON response
func (ve ValidationErrors) ToMap() map[string]string {
	result := make(map[string]string)
	for _, err := range ve.Errors {
		result[err.Field] = err.Message
	}
	return result
}
