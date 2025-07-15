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
	Logger           Logger            `json:"-"`
	CustomMessages   bool              `json:"custom_messages"`
	StopOnFirstError bool              `json:"stop_on_first"`
	Locale           string            `json:"locale"`
	CustomRules      map[string]Rule   `json:"-"`
	ErrorMessages    map[string]string `json:"error_messages"`
}

// ValidationConfig for request-level options
type ValidationConfig struct {
	Context   map[string]interface{}
	ForceJSON bool
	ForceForm bool
	Locale    string
	Rules     map[string]Rule
}

// ValidationOption for flexible configuration
type ValidationOption func(*ValidationConfig)

// ValidationError represents validation failure details
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value,omitempty"`
	Tag     string      `json:"tag"`
	Message string      `json:"message"`
}

// ValidationErrors collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Rule interface for custom validation rules
type Rule interface {
	Validate(value interface{}, params string, context map[string]interface{}) error
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

// Logger interface (reuse from response package)
type Logger interface {
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NoOpLogger implements Logger interface with no-op methods
type NoOpLogger struct{}

func (l *NoOpLogger) Error(msg string, fields ...interface{}) {}
func (l *NoOpLogger) Warn(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Info(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Debug(msg string, fields ...interface{}) {}

// JSONWriter interface for framework-agnostic response writing
type JSONWriter interface {
	WriteJSON(statusCode int, data interface{}) error
	GetStatusCode() int
}

// Default configuration
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
		Context: make(map[string]interface{}),
		Locale:  "en",
		Rules:   make(map[string]Rule),
	}

	for _, opt := range opts {
		opt(&config)
	}

	return config
}

// Validation options
func WithContext(ctx map[string]interface{}) ValidationOption {
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
