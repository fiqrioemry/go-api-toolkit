package validation

import (
	"reflect"

	"go.uber.org/zap"
)

// Logger interface (same as before)
type Logger interface {
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NoOpLogger (same as before)
type NoOpLogger struct{}

func (l *NoOpLogger) Error(msg string, fields ...interface{}) {}
func (l *NoOpLogger) Warn(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Info(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Debug(msg string, fields ...interface{}) {}

// ZapSugarAdapter (optimized version)
type ZapSugarAdapter struct {
	logger *zap.SugaredLogger
}

func (z *ZapSugarAdapter) Error(msg string, fields ...interface{}) {
	z.logger.Errorw(msg, fields...)
}

func (z *ZapSugarAdapter) Warn(msg string, fields ...interface{}) {
	z.logger.Warnw(msg, fields...)
}

func (z *ZapSugarAdapter) Info(msg string, fields ...interface{}) {
	z.logger.Infow(msg, fields...)
}

func (z *ZapSugarAdapter) Debug(msg string, fields ...interface{}) {
	z.logger.Debugw(msg, fields...)
}

// adaptLogger automatically detects logger type and creates appropriate adapter
func adaptLogger(logger interface{}) Logger {
	if logger == nil {
		return &NoOpLogger{}
	}

	// Check if it already implements our Logger interface
	if validationLogger, ok := logger.(Logger); ok {
		return validationLogger
	}

	// Auto-detect and adapt different logger types
	loggerValue := reflect.ValueOf(logger)
	loggerType := reflect.TypeOf(logger)

	// Handle pointer types
	if loggerType.Kind() == reflect.Ptr {
		loggerType = loggerType.Elem()
		if loggerValue.IsNil() {
			return &NoOpLogger{}
		}
	}

	// Detect zap.Logger
	if loggerType.Name() == "Logger" && loggerType.PkgPath() == "go.uber.org/zap" {
		if zapLogger, ok := logger.(*zap.Logger); ok {
			return &ZapSugarAdapter{logger: zapLogger.Sugar()}
		}
	}

	// Detect zap.SugaredLogger
	if loggerType.Name() == "SugaredLogger" && loggerType.PkgPath() == "go.uber.org/zap" {
		if sugarLogger, ok := logger.(*zap.SugaredLogger); ok {
			return &ZapSugarAdapter{logger: sugarLogger}
		}
	}

	// Try to detect by method signature (duck typing)
	if hasCompatibleMethods(loggerValue) {
		return createDuckTypedAdapter(logger)
	}

	// Fallback to NoOpLogger if can't adapt
	return &NoOpLogger{}
}

// hasCompatibleMethods checks if logger has compatible method signatures
func hasCompatibleMethods(loggerValue reflect.Value) bool {
	required := []string{"Error", "Warn", "Info", "Debug"}

	for _, methodName := range required {
		method := loggerValue.MethodByName(methodName)
		if !method.IsValid() {
			return false
		}

		methodType := method.Type()
		if methodType.NumIn() < 1 {
			return false
		}

		// First parameter should be string
		if methodType.In(0).Kind() != reflect.String {
			return false
		}
	}

	return true
}

// createDuckTypedAdapter creates adapter for loggers with compatible interface
func createDuckTypedAdapter(logger interface{}) Logger {
	return &DuckTypedAdapter{logger: logger}
}

// DuckTypedAdapter adapts any logger with compatible methods
type DuckTypedAdapter struct {
	logger interface{}
}

func (d *DuckTypedAdapter) Error(msg string, fields ...interface{}) {
	d.callMethod("Error", msg, fields...)
}

func (d *DuckTypedAdapter) Warn(msg string, fields ...interface{}) {
	d.callMethod("Warn", msg, fields...)
}

func (d *DuckTypedAdapter) Info(msg string, fields ...interface{}) {
	d.callMethod("Info", msg, fields...)
}

func (d *DuckTypedAdapter) Debug(msg string, fields ...interface{}) {
	d.callMethod("Debug", msg, fields...)
}

func (d *DuckTypedAdapter) callMethod(methodName, msg string, fields ...interface{}) {
	loggerValue := reflect.ValueOf(d.logger)
	method := loggerValue.MethodByName(methodName)

	if !method.IsValid() {
		return
	}

	// Prepare arguments
	args := make([]reflect.Value, 1+len(fields))
	args[0] = reflect.ValueOf(msg)

	for i, field := range fields {
		args[i+1] = reflect.ValueOf(field)
	}

	// Call method safely
	defer func() {
		if r := recover(); r != nil {
			// Silently ignore if method call fails
		}
	}()

	method.Call(args)
}
