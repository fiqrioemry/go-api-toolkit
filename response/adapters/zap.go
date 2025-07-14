// ==================== response/adapters/zap.go ====================

package adapters

import (
	"github.com/fiqrioemry/go-api-toolkit/response"
	"go.uber.org/zap"
)

// ZapLogger implements Logger interface for Zap
type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) Debug(msg string, fields ...response.LogField) {
	z.logger.Debug(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...response.LogField) {
	z.logger.Info(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...response.LogField) {
	z.logger.Warn(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...response.LogField) {
	z.logger.Error(msg, z.convertFields(fields)...)
}

func (z *ZapLogger) convertFields(fields []response.LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
