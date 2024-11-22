package zlog

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	lg      *zap.Logger
	opt     Option
	ctxHook func(ctx context.Context) []zapcore.Field
}

func (l *Logger) Debug(ctx ...context.Context) *Event {
	return l.newEvent(withCtx(ctx), zapcore.DebugLevel)
}

func (l *Logger) Info(ctx ...context.Context) *Event {
	return l.newEvent(withCtx(ctx), zapcore.InfoLevel)
}

func (l *Logger) Warn(ctx ...context.Context) *Event {
	return l.newEvent(withCtx(ctx), zapcore.WarnLevel)
}

func (l *Logger) Error(ctx ...context.Context) *Event {
	return l.newEvent(withCtx(ctx), zapcore.ErrorLevel)
}

func (l *Logger) Fatal(ctx ...context.Context) *Event {
	return l.newEvent(withCtx(ctx), zapcore.FatalLevel)
}

func (l *Logger) Println(args ...interface{}) {
	l.newEvent(nil, zapcore.InfoLevel).Msg(strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.newEvent(nil, zapcore.InfoLevel).Msgf(format, args)
}

func (l *Logger) With(key string, val interface{}) *Logger {
	return &Logger{lg: l.lg.With(zap.Any(key, val)), opt: l.opt}
}

func (l *Logger) AddCtxHook(fn CtxHook) *Logger {
	l.ctxHook = func(ctx context.Context) []zapcore.Field {
		data := fn(ctx)
		fields := make([]zapcore.Field, 0, len(data))

		for k, v := range data {
			fields = append(fields, zap.Any(k, v))
		}

		return fields
	}

	return l
}

func (l *Logger) Sync() error {
	return l.lg.Sync()
}

func (l *Logger) newEvent(ctx context.Context, level zapcore.Level) *Event {
	e := newEvent(ctx, l.lg, level)

	if ctx != nil && l.ctxHook != nil {
		e.fields = l.ctxHook(ctx)
	}

	return e
}

func newLogger(opt Option) *zap.Logger {
	writeSyncer := newWriter(opt)
	encoder := newEncoder()

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	zlg := zap.New(core, zap.AddCaller())
	zlg = zlg.WithOptions(zap.AddCallerSkip(opt.SkipLevel))

	return zlg
}

func New(opt Option) *Logger {
	return &Logger{lg: newLogger(opt), opt: opt}
}
