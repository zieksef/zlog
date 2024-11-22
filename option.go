package zlog

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultSkipLevel = 1
	defaultMaxSize   = 20
	defaultMaxAge    = 15
)

type Option struct {
	Dir            string
	Filename       string
	DisableConsole bool
	DailyRotation  bool
	MaxSize        int // in MB
	MaxAge         int // in day
	SkipLevel      int
	Timezone       *time.Location
	Writer         io.Writer
}

// DefaultOption with Asia/Shanghai timezone has no daily rotation and only enables console output.
func DefaultOption() Option {
	timezone, _ := time.LoadLocation("Asia/Shanghai")

	return Option{
		SkipLevel: defaultSkipLevel,
		MaxSize:   defaultMaxSize,
		MaxAge:    defaultMaxAge,
		Timezone:  timezone,
		Writer:    os.Stdout,
	}
}

func (o Option) SetDir(dir string) Option {
	o.Dir = dir
	return o
}

func (o Option) SetFilename(filename string) Option {
	o.Filename = filename
	return o
}

func (o Option) SetSkipLevel(level int) Option {
	o.SkipLevel = level
	return o
}

func (o Option) SetTimezone(tz *time.Location) Option {
	o.Timezone = tz
	return o
}

func (o Option) SetMaxSize(size int) Option {
	o.MaxSize = size
	return o
}

func (o Option) SetMaxAge(age int) Option {
	o.MaxAge = age
	return o
}

func (o Option) SetWriter(w io.Writer) Option {
	o.Writer = w
	return o
}

func (o Option) SetDailyRotation(rotation bool) Option {
	o.DailyRotation = rotation
	return o
}

func (o Option) SetDisableConsole(disable bool) Option {
	o.DisableConsole = disable
	return o
}

func newEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(cfg)
}

func newWriter(opt Option) zapcore.WriteSyncer {
	if opt.Writer == nil {
		opt.Writer = os.Stdout
	}

	consoleSyncer := zapcore.AddSync(opt.Writer)

	if opt.Filename == "" {
		return zapcore.NewMultiWriteSyncer(consoleSyncer) // log file name not provided, output to console
	}

	if opt.MaxSize < 0 {
		opt.MaxSize = defaultMaxSize
	}

	if opt.MaxAge < 0 {
		opt.MaxAge = defaultMaxAge
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename: filepath.Join(opt.Dir, opt.Filename),
		MaxSize:  defaultMaxSize,
		MaxAge:   opt.MaxAge,
	}

	if opt.DailyRotation {
		go func() {
			for {
				now := time.Now()
				nowStr := now.Format(time.DateOnly)

				t2, _ := time.ParseInLocation(time.DateOnly, nowStr, opt.Timezone)

				next := t2.AddDate(0, 0, 1)

				// subtract 1 nanosecond to ensure the timer triggers before the next rotation time point
				after := next.UnixNano() - now.UnixNano() - 1

				<-time.After(time.Duration(after) * time.Nanosecond)

				_ = lumberJackLogger.Rotate()
			}
		}()
	}

	fileSyncer := zapcore.AddSync(lumberJackLogger)

	if opt.DisableConsole {
		return zapcore.NewMultiWriteSyncer(fileSyncer)
	}
	return zapcore.NewMultiWriteSyncer(consoleSyncer, fileSyncer)
}
