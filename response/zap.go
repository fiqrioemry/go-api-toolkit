// ==================== response/zap_adapter.go ====================
package response

// ZapLogger implements Logger interface for Zap
// This is a placeholder for the actual implementation
// which would use the zap package for structured logging.
import "go.uber.org/zap"

// ZapLogger implements Logger interface for Zap
type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) Debug(msg string, fields ...LogField) {
	z.logger.Debug(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...LogField) {
	z.logger.Info(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...LogField) {
	z.logger.Warn(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...LogField) {
	z.logger.Error(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) convertFields(fields []LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
