// Gin adapter
package adapters

import (
	"github.com/fiqrioemry/go-api-toolkit/response"
	"github.com/gin-gonic/gin"
)

// GinJSONWriter implements JSONWriter for Gin framework
type GinJSONWriter struct {
	ctx *gin.Context
}

func (g *GinJSONWriter) JSON(statusCode int, obj any) {
	g.ctx.JSON(statusCode, obj)
}

// GinContextExtractor extracts context from Gin request
func GinContextExtractor(req any) *response.Context {
	if ginCtx, ok := req.(*gin.Context); ok {
		return &response.Context{
			Path:      ginCtx.Request.URL.Path,
			Method:    ginCtx.Request.Method,
			ClientIP:  ginCtx.ClientIP(),
			UserAgent: ginCtx.Request.UserAgent(),
			UserID:    ginCtx.GetString("user_id"),  // Assuming user_id is set in context
			TraceID:   ginCtx.GetString("trace_id"), // Assuming trace_id is set in context
		}
	}
	return &response.Context{}
}

// NewGinHandler creates a response handler for Gin
func NewGinHandler(options ...response.Option) *response.Handler {
	options = append(options, response.WithContextExtractor(GinContextExtractor))
	return response.NewHandler(options...)
}

// Helper functions for Gin
func HandleError(handler *response.Handler, c *gin.Context, err error) {
	writer := &GinJSONWriter{ctx: c}
	handler.HandleError(writer, c, err)
}

func OK(handler *response.Handler, c *gin.Context, message string, data any) {
	writer := &GinJSONWriter{ctx: c}
	handler.OK(writer, c, message, data)
}

func Created(handler *response.Handler, c *gin.Context, message string, data any) {
	writer := &GinJSONWriter{ctx: c}
	handler.Created(writer, c, message, data)
}

func OKWithPagination(handler *response.Handler, c *gin.Context, message string, data any, pagination *response.Pagination) {
	writer := &GinJSONWriter{ctx: c}
	handler.OKWithPagination(writer, c, message, data, pagination)
}
