package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var ZapLogger *zap.Logger

// LoggerKeyType constant for context log name
const LoggerKeyType = "logger_key"

type Event map[string]interface{}

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "log_level"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.TimeKey = "timestamp_app"

	var config zap.Config

	if strings.EqualFold(os.Getenv("APP_DEV"), "true") {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig = encoderConfig
	config.OutputPaths = []string{"stdout"}
	config.DisableCaller = true

	if strings.EqualFold(os.Getenv("APP_DEBUG"), "true") {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	var err error
	ZapLogger, err = config.Build()
	if err != nil {
		panic(err)
	}

	ZapLogger = ZapLogger.With(
		zap.String("log_type", "APPLICATION"),
	)
}

func addEvent(event ...Event) zap.Field {
	if len(event) > 0 {
		return zap.Any("event", event[0])
	}
	return zap.Any("event", nil)
}

func Info(context context.Context, msg string, event ...Event) {
	WithContext(context).Info(msg, addEvent(event...))
}

func Warn(context context.Context, msg string, event ...Event) {
	WithContext(context).Warn(msg, addEvent(event...))
}

func Error(context context.Context, msg string, err error, event ...Event) {
	WithContext(context).Error(msg, zap.Error(err), addEvent(event...))
}

func Fatal(context context.Context, msg string, err error, event ...Event) {
	WithContext(context).Fatal(msg, zap.Error(err), addEvent(event...))
}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil || ctx == context.Background() || ctx == context.TODO() {
		return ZapLogger
	}

	if ctxLogger, ok := ctx.Value(LoggerKeyType).(*zap.Logger); ok {
		return ctxLogger
	}

	return ZapLogger
}

// Close flushing any buffered log entries
func Close() error {
	return ZapLogger.Sync()
}
