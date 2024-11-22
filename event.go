package zlog

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	initCap = 4
)

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{fields: make([]zapcore.Field, 0, initCap)}
	},
}

type Event struct {
	lg     *zap.Logger
	level  zapcore.Level
	fields []zapcore.Field
}

func (e *Event) Int(key string, val int) *Event {
	e.fields = append(e.fields, zap.Int(key, val))
	return e
}

func (e *Event) Int64(key string, val int64) *Event {
	e.fields = append(e.fields, zap.Int64(key, val))
	return e
}

func (e *Event) Float64(key string, val float64) *Event {
	e.fields = append(e.fields, zap.Float64(key, val))
	return e
}

func (e *Event) Bool(key string, val bool) *Event {
	e.fields = append(e.fields, zap.Bool(key, val))
	return e
}

func (e *Event) Str(key string, val string) *Event {
	e.fields = append(e.fields, zap.String(key, val))
	return e
}

func (e *Event) Time(key string, val time.Time) *Event {
	e.fields = append(e.fields, zap.Time(key, val))
	return e
}

func (e *Event) Err(err error) *Event {
	e.fields = append(e.fields, zap.Error(err))
	return e
}

func (e *Event) Any(key string, val interface{}) *Event {
	e.fields = append(e.fields, zap.Any(key, val))
	return e
}

func (e *Event) Msg(msg string) {
	e.lg.Log(e.level, msg, e.fields...)
	putEvent(e)
}

func (e *Event) Msgf(format string, args ...interface{}) {
	e.lg.Log(e.level, fmt.Sprintf(format, args...), e.fields...)
	putEvent(e)
}

func (e *Event) Send() {
	e.lg.Log(e.level, "", e.fields...)
	putEvent(e)
}

func newEvent(_ context.Context, logger *zap.Logger, level zapcore.Level) *Event {
	event := eventPool.Get().(*Event)
	event.lg = logger
	event.level = level

	return event
}

func putEvent(e *Event) {
	e.lg = nil
	e.fields = e.fields[:0]
	eventPool.Put(e)
}
