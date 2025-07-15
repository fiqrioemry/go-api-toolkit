package validation

import (
	"strings"

	"github.com/gin-gonic/gin"
	// Import your response package here
	// "github.com/fiqrioemry/go-api-toolkit/response"
)

// InitGin initializes validation for Gin framework
func InitGin(configs ...InitConfig) {
	var config InitConfig
	if len(configs) > 0 {
		config = configs[0]
	}

	// Apply defaults
	if config.Logger == nil {
		config.Logger = &NoOpLogger{}
	}
	if config.Locale == "" {
		config.Locale = "en"
	}
	if config.CustomRules == nil {
		config.CustomRules = make(map[string]Rule)
	}
	if config.ErrorMessages == nil {
		config.ErrorMessages = getDefaultMessages(config.Locale)
	}

	// Create handler config
	handlerConfig := Config{
		Logger:           config.Logger,
		CustomMessages:   config.CustomMessages,
		StopOnFirstError: config.StopOnFirstError,
		Locale:           config.Locale,
		CustomRules:      config.CustomRules,
		ErrorMessages:    config.ErrorMessages,
	}

	// Initialize global handler
	globalHandler = NewHandler(handlerConfig)
}

// BindAndValidate binds request data and validates struct
func BindAndValidate(c *gin.Context, obj interface{}, opts ...ValidationOption) error {
	if globalHandler == nil {
		// Auto-initialize with defaults if not initialized
		InitGin()
	}

	config := buildValidationConfig(opts...)

	// Smart binding based on request
	if err := smartBind(c, obj, config); err != nil {
		globalHandler.config.Logger.Error("validation: binding failed",
			"error", err.Error(),
			"path", c.FullPath(),
			"method", c.Request.Method,
		)

		// Return as BadRequest error (assuming response package integration)
		return createBindingError(err.Error())
	}

	// Validate the bound struct
	if err := globalHandler.ValidateStruct(obj, config.Context); err != nil {
		if validationErrs, ok := err.(ValidationErrors); ok {
			globalHandler.config.Logger.Debug("validation: validation failed",
				"errors", validationErrs.ToMap(),
				"path", c.FullPath(),
			)

			// Return as ValidationError (assuming response package integration)
			return createValidationError(validationErrs)
		}

		globalHandler.config.Logger.Error("validation: unexpected validation error",
			"error", err.Error(),
			"path", c.FullPath(),
		)

		return createValidationError(ValidationErrors{
			Errors: []ValidationError{{
				Field:   "unknown",
				Message: err.Error(),
			}},
		})
	}

	return nil
}

// smartBind performs smart binding based on request characteristics
func smartBind(c *gin.Context, obj interface{}, config ValidationConfig) error {
	// Force specific binding if requested
	if config.ForceJSON {
		return c.ShouldBindJSON(obj)
	}

	if config.ForceForm {
		return bindForm(c, obj)
	}

	// Auto-detect based on Content-Type and method
	contentType := c.GetHeader("Content-Type")
	method := c.Request.Method

	switch {
	case strings.Contains(contentType, "application/json"):
		return c.ShouldBindJSON(obj)
	case strings.Contains(contentType, "multipart/form-data"):
		return c.ShouldBind(obj)
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		return c.ShouldBind(obj)
	case method == "GET" || method == "DELETE":
		// For GET/DELETE, prefer query parameters
		return c.ShouldBindQuery(obj)
	default:
		// Fallback: try JSON for POST/PUT/PATCH, query for others
		if method == "POST" || method == "PUT" || method == "PATCH" {
			return c.ShouldBindJSON(obj)
		}
		return c.ShouldBindQuery(obj)
	}
}

// bindForm binds form data (query params or form-data)
func bindForm(c *gin.Context, obj interface{}) error {
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return c.ShouldBind(obj)
	}

	// Default to query parameters
	return c.ShouldBindQuery(obj)
}

// Gin-specific writer adapter
type ginJSONWriter struct {
	ctx        *gin.Context
	statusCode int
}

func (w *ginJSONWriter) WriteJSON(statusCode int, data interface{}) error {
	w.statusCode = statusCode
	w.ctx.JSON(statusCode, data)
	return nil
}

func (w *ginJSONWriter) GetStatusCode() int {
	return w.statusCode
}

// Helper functions to create errors (these would integrate with your response package)
func createBindingError(message string) error {
	// This would use your response package's BadRequest function
	// return response.BadRequest("Invalid request format: " + message)

	// Placeholder implementation
	return ValidationErrors{
		Errors: []ValidationError{{
			Field:   "request",
			Message: "Invalid request format: " + message,
			Tag:     "binding",
		}},
	}
}

func createValidationError(validationErrs ValidationErrors) error {
	// This would use your response package's ValidationError function
	// return response.ValidationError("Validation failed", validationErrs.ToMap())

	// Return the validation errors as-is for now
	return validationErrs
}

// Quick validation function for simple cases
func Validate(obj interface{}, opts ...ValidationOption) error {
	if globalHandler == nil {
		InitGin()
	}

	config := buildValidationConfig(opts...)
	return globalHandler.ValidateStruct(obj, config.Context)
}

// ValidateWithContext validates with context
func ValidateWithContext(obj interface{}, context map[string]interface{}) error {
	return Validate(obj, WithContext(context))
}
