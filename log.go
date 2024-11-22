package log

import (
	"context"

	"go.uber.org/zap/zapcore"
)

var logger = New(DefaultOption())

func Debug(ctx ...context.Context) *Event {
	return logger.newEvent(withCtx(ctx), zapcore.DebugLevel)
}

func Info(ctx ...context.Context) *Event {
	return logger.newEvent(withCtx(ctx), zapcore.InfoLevel)
}

func Warn(ctx ...context.Context) *Event {
	return logger.newEvent(withCtx(ctx), zapcore.WarnLevel)
}

func Error(ctx ...context.Context) *Event {
	return logger.newEvent(withCtx(ctx), zapcore.ErrorLevel)
}

func Fatal(ctx ...context.Context) *Event {
	return logger.newEvent(withCtx(ctx), zapcore.FatalLevel)
}

func Println(args ...interface{}) {
	logger.Println(args...)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func With(key string, val interface{}) *Logger {
	return logger.With(key, val)
}

func AddCtxHook(fn CtxHook) *Logger {
	return logger.AddCtxHook(fn)
}

func Sync() error {
	return logger.Sync()
}

func Init(opt Option) {
	logger = &Logger{lg: newLogger(opt), opt: opt}
}
