// ==================== response/logger.go ====================
package response

// LogLevel represents log levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// LogField represents a log field
type LogField struct {
	Key   string
	Value any
}

// Logger interface for flexible logging
type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
}

// NoOpLogger implements Logger interface but does nothing
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, fields ...LogField) {}
func (n *NoOpLogger) Info(msg string, fields ...LogField)  {}
func (n *NoOpLogger) Warn(msg string, fields ...LogField)  {}
func (n *NoOpLogger) Error(msg string, fields ...LogField) {}
