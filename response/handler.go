// ==================== response/handler.go ====================
package response

import "net/http"

// Handler handles HTTP responses with logging
type Handler struct {
	logger           Logger
	contextExtractor ContextExtractor
	config           *Config
}

// Config represents handler configuration
type Config struct {
	LogSuccessResponses bool
	LogErrorResponses   bool
	LogLevel            LogLevel
	IncludeStackTrace   bool
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		LogSuccessResponses: false,
		LogErrorResponses:   true,
		LogLevel:            LogLevelInfo,
		IncludeStackTrace:   false,
	}
}

// JSONWriter interface for framework-agnostic JSON responses
type JSONWriter interface {
	JSON(statusCode int, obj any)
}

// NewHandler creates a new response handler
func NewHandler(options ...Option) *Handler {
	h := &Handler{
		logger: &NoOpLogger{},
		config: DefaultConfig(),
	}

	for _, opt := range options {
		opt(h)
	}

	return h
}

// Option represents handler option
type Option func(*Handler)

// WithLogger sets the logger
func WithLogger(logger Logger) Option {
	return func(h *Handler) {
		h.logger = logger
	}
}

// WithContextExtractor sets the context extractor
func WithContextExtractor(extractor ContextExtractor) Option {
	return func(h *Handler) {
		h.contextExtractor = extractor
	}
}

// WithConfig sets the configuration
func WithConfig(config *Config) Option {
	return func(h *Handler) {
		h.config = config
	}
}

// HandleError handles error responses
func (h *Handler) HandleError(w JSONWriter, req any, err error) {
	ctx := h.extractContext(req)

	if appErr, ok := IsAppError(err); ok {
		response := ErrorResponse{
			Success: false,
			Message: appErr.Message,
			Code:    appErr.Code,
		}

		if appErr.Context != nil {
			if errorDetails, exists := appErr.Context["errors"]; exists {
				response.Errors = errorDetails.(map[string]any)
			}
		}

		if h.config.LogErrorResponses {
			h.logError(ctx, appErr)
		}

		w.JSON(appErr.HTTPStatus, response)
		return
	}

	// Handle unknown errors
	response := ErrorResponse{
		Success: false,
		Message: "Internal server error",
		Code:    ErrCodeInternalServer,
	}

	if h.config.LogErrorResponses {
		h.logUnknownError(ctx, err)
	}

	w.JSON(http.StatusInternalServerError, response)
}

// Success sends success response
func (h *Handler) Success(w JSONWriter, req any, statusCode int, message string, data any) {
	response := SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	if h.config.LogSuccessResponses {
		ctx := h.extractContext(req)
		h.logSuccess(ctx, statusCode, message)
	}

	w.JSON(statusCode, response)
}

// SuccessWithMeta sends success response with metadata
func (h *Handler) SuccessWithMeta(w JSONWriter, req any, statusCode int, message string, data any, meta *Meta) {
	response := SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	if h.config.LogSuccessResponses {
		ctx := h.extractContext(req)
		h.logSuccess(ctx, statusCode, message)
	}

	w.JSON(statusCode, response)
}

// OK sends 200 OK response
func (h *Handler) OK(w JSONWriter, req any, message string, data any) {
	h.Success(w, req, http.StatusOK, message, data)
}

// Created sends 201 Created response
func (h *Handler) Created(w JSONWriter, req any, message string, data any) {
	h.Success(w, req, http.StatusCreated, message, data)
}

// OKWithPagination sends 200 OK response with pagination
func (h *Handler) OKWithPagination(w JSONWriter, req any, message string, data any, pagination *Pagination) {
	h.SuccessWithMeta(w, req, http.StatusOK, message, data, &Meta{
		Pagination: pagination,
	})
}

// OKWithPaginationAndPermissions sends 200 OK response with pagination and permissions
func (h *Handler) OKWithPaginationAndPermissions(w JSONWriter, req any, message string, data any, pagination *Pagination, permissions map[string]bool) {
	h.SuccessWithMeta(w, req, http.StatusOK, message, data, &Meta{
		Pagination:  pagination,
		Permissions: permissions,
	})
}

// extractContext extracts context from request
func (h *Handler) extractContext(req any) *Context {
	if h.contextExtractor != nil {
		return h.contextExtractor(req)
	}
	return &Context{}
}

// logError logs application errors
func (h *Handler) logError(ctx *Context, appErr *AppError) {
	fields := h.buildLogFields(ctx)
	fields = append(fields,
		LogField{Key: "error_code", Value: string(appErr.Code)},
		LogField{Key: "error_message", Value: appErr.Message},
	)

	if IsServerError(appErr) {
		if appErr.Err != nil {
			fields = append(fields, LogField{Key: "underlying_error", Value: appErr.Err.Error()})
		}
		h.logger.Error("Server error occurred", fields...)
	} else {
		h.logger.Warn("Client error occurred", fields...)
	}
}

// logUnknownError logs unknown errors
func (h *Handler) logUnknownError(ctx *Context, err error) {
	fields := h.buildLogFields(ctx)
	fields = append(fields, LogField{Key: "error", Value: err.Error()})
	h.logger.Error("Unknown error occurred", fields...)
}

// logSuccess logs successful responses
func (h *Handler) logSuccess(ctx *Context, statusCode int, message string) {
	fields := h.buildLogFields(ctx)
	fields = append(fields,
		LogField{Key: "status_code", Value: statusCode},
		LogField{Key: "message", Value: message},
	)
	h.logger.Info("Success response", fields...)
}

// buildLogFields builds common log fields
func (h *Handler) buildLogFields(ctx *Context) []LogField {
	fields := []LogField{
		{Key: "path", Value: ctx.Path},
		{Key: "method", Value: ctx.Method},
		{Key: "client_ip", Value: ctx.ClientIP},
		{Key: "user_agent", Value: ctx.UserAgent},
	}

	if ctx.UserID != "" {
		fields = append(fields, LogField{Key: "user_id", Value: ctx.UserID})
	}

	if ctx.TraceID != "" {
		fields = append(fields, LogField{Key: "trace_id", Value: ctx.TraceID})
	}

	return fields
}
