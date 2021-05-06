package logger

import (
	"context"
	"os"
	"time"

	"github.com/piyuo/libsrv/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger global instance, cause log use frequently
//
var logger = initial()

// initLogger init zap logger
//
func initial() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	atom := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if env.Debug {
		atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(env.AppName)
	}
	encoder := zapcore.NewConsoleEncoder(config)
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom), zap.AddCaller(), zap.AddCallerSkip(2))
}

// addContextInformation add context information to zap fields
//
//	fields := addContextInformation(ctx)
//
func addContextInformation(ctx context.Context) []zap.Field {
	var fields []zap.Field
	user := env.GetUserID(ctx)
	if user != "" {
		fields = []zap.Field{zap.String("user", user)}
	}
	return fields
}

// Debug only print message when os.Getenv("DEBUG") is defined
//
//	Debug(ctx,"server start")
//
func Debug(ctx context.Context, message string) {
	logger.Debug(message, addContextInformation(ctx)...)
}

// Info as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	Info(ctx,"server start")
//
func Info(ctx context.Context, message string) {
	logger.Info(message, addContextInformation(ctx)...)
}

// Warn events might cause problems.
//
//	Warning(ctx,"hi")
//
func Warn(ctx context.Context, message string) {
	logger.Warn(message, addContextInformation(ctx)...)
}

// Error write error log
//
//	Error(ctx,"error")
//
func Error(ctx context.Context, message string) {
	logger.Error(message, addContextInformation(ctx)...)
}
