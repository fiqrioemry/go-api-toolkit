// File: validation/gin.go - Updated to match your response package pattern

package validation

import (
	"strings"

	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

// InitGin initializes validation for Gin framework
func InitGin(configs ...InitConfig) {
	var config InitConfig
	if len(configs) > 0 {
		config = configs[0]
	}

	handlerConfig := applyConfigDefaults(config)
	globalHandler = NewHandler(handlerConfig)
}

// BindAndValidate binds request data and validates struct
func BindAndValidate(c *gin.Context, obj any, opts ...ValidationOption) error {
	if globalHandler == nil {
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

		// 🚀 USE YOUR RESPONSE PACKAGE PATTERN
		return response.BadRequest("Invalid request format: " + err.Error())
	}

	// Validate the bound struct
	if err := globalHandler.ValidateStruct(obj, config.Context); err != nil {
		if validationErrs, ok := err.(ValidationErrors); ok {
			globalHandler.config.Logger.Debug("validation: validation failed",
				"errors", validationErrs.ToMap(),
				"path", c.FullPath(),
			)

			// 🚀 CREATE VALIDATION ERROR USING YOUR PATTERN
			return createValidationError(validationErrs)
		}

		globalHandler.config.Logger.Error("validation: unexpected validation error",
			"error", err.Error(),
			"path", c.FullPath(),
		)

		return response.InternalServerError("Validation error", err)
	}

	return nil
}

// 🚀 HELPER FUNCTION TO CREATE VALIDATION ERROR USING YOUR PATTERN
func createValidationError(validationErrs ValidationErrors) error {
	// Create BadRequest with validation failed message
	err := response.NewBadRequest("Validation failed")

	// Add validation errors as context (same pattern as your utils)
	err.WithContext("errors", validationErrs.ToMap())

	return err
}

// smartBind performs smart binding based on request characteristics
func smartBind(c *gin.Context, obj any, config ValidationConfig) error {
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
		return c.ShouldBindQuery(obj)
	default:
		if method == "POST" || method == "PUT" || method == "PATCH" {
			return c.ShouldBindJSON(obj)
		}
		return c.ShouldBindQuery(obj)
	}
}

// bindForm binds form data (query params or form-data)
func bindForm(c *gin.Context, obj any) error {
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return c.ShouldBind(obj)
	}

	return c.ShouldBindQuery(obj)
}

// Quick validation function for simple cases
func Validate(obj any, opts ...ValidationOption) error {
	if globalHandler == nil {
		InitGin()
	}

	config := buildValidationConfig(opts...)
	return globalHandler.ValidateStruct(obj, config.Context)
}

// ValidateWithContext validates with context
func ValidateWithContext(obj any, context map[string]any) error {
	return Validate(obj, WithContext(context))
}

// BindAndValidateJSON as drop-in replacement for utils.BindAndValidateJSON
func BindAndValidateJSON[T any](c *gin.Context, req *T) bool {
	if err := BindAndValidate(c, req, ForceJSON()); err != nil {
		response.Error(c, err)
		return false
	}
	return true
}

// BindAndValidateForm as drop-in replacement for utils.BindAndValidateForm
func BindAndValidateForm[T any](c *gin.Context, req *T) bool {
	if err := BindAndValidate(c, req, ForceForm()); err != nil {
		response.Error(c, err)
		return false
	}
	return true
}

// BindAndValidateQuery for query parameters
func BindAndValidateQuery[T any](c *gin.Context, req *T) bool {
	if err := BindAndValidate(c, req, ForceQuery()); err != nil {
		response.Error(c, err)
		return false
	}
	return true
}
