// ==================== response/context.go ====================
package response

// Context represents request context for logging
type Context struct {
	Path      string
	Method    string
	ClientIP  string
	UserAgent string
	UserID    string
	TraceID   string
	Headers   map[string]string
}

// ContextExtractor defines how to extract context from request
type ContextExtractor func(req any) *Context
