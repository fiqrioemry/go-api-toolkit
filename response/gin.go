// ==================== response/gin.go ====================
package response

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Global handler - initialized once
var globalHandler *Handler

// InitConfig for simple initialization
type InitConfig struct {
	Logger              *zap.Logger
	LogSuccessResponses bool
	LogErrorResponses   bool
}

// GinJSONWriter implements JSONWriter for Gin framework
type GinJSONWriter struct {
	ctx *gin.Context
}

func (g *GinJSONWriter) JSON(statusCode int, obj any) {
	g.ctx.JSON(statusCode, obj)
}

// GinContextExtractor extracts context from Gin request
func GinContextExtractor(req any) *Context {
	if ginCtx, ok := req.(*gin.Context); ok {
		return &Context{
			Path:      ginCtx.Request.URL.Path,
			Method:    ginCtx.Request.Method,
			ClientIP:  ginCtx.ClientIP(),
			UserAgent: ginCtx.Request.UserAgent(),
			UserID:    ginCtx.GetString("user_id"),
			TraceID:   ginCtx.GetString("trace_id"),
		}
	}
	return &Context{}
}

// InitGin initializes the global response handler for Gin
func InitGin(config InitConfig) {
	logger := NewZapLogger(config.Logger)

	handlerConfig := &Config{
		LogSuccessResponses: config.LogSuccessResponses,
		LogErrorResponses:   config.LogErrorResponses,
		LogLevel:            LogLevelInfo,
	}

	globalHandler = NewHandler(
		WithLogger(logger),
		WithContextExtractor(GinContextExtractor),
		WithConfig(handlerConfig),
	)
}

// ============ RESPONSE FUNCTIONS ============

func Error(c *gin.Context, err error) {
	writer := &GinJSONWriter{ctx: c}
	globalHandler.HandleError(writer, c, err)
}

func OK(c *gin.Context, message string, data any) {
	writer := &GinJSONWriter{ctx: c}
	globalHandler.OK(writer, c, message, data)
}

func Created(c *gin.Context, message string, data any) {
	writer := &GinJSONWriter{ctx: c}
	globalHandler.Created(writer, c, message, data)
}

func BadRequestMsg(c *gin.Context, message string) {
	err := NewBadRequest(message)
	Error(c, err)
}

func NotFoundMsg(c *gin.Context, message string) {
	err := NewNotFound(message)
	Error(c, err)
}

func UnauthorizedMsg(c *gin.Context, message string) {
	err := NewUnauthorized(message)
	Error(c, err)
}

func ForbiddenMsg(c *gin.Context, message string) {
	err := NewForbidden(message)
	Error(c, err)
}
