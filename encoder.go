package zlog

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	red   = "\033[31m"
	reset = "\033[0m"
)

type coloredConsoleEncoder struct {
	zapcore.Encoder
}

func (c coloredConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	lineBuffer, _ := c.Encoder.EncodeEntry(ent, fields)

	fds := strings.Split(lineBuffer.String(), "\t")
	if len(fds) == 0 {
		return lineBuffer, nil
	}

	lastIdx := len(fds) - 1
	fieldData := fds[lastIdx]

	var fieldDataMap map[string]interface{}
	if err := json.Unmarshal([]byte(fieldData), &fieldDataMap); err != nil {
		fmt.Println(err)
		return lineBuffer, nil
	}

	i := 0
	result := "{"
	for k, v := range fieldDataMap {
		if k == "error" {
			result += fmt.Sprintf(`"%s": "%s%v%s"`, k, red, v, reset)
			fmt.Println(result)
		} else {
			result += fmt.Sprintf(`"%s": "%v"`, k, v)
		}

		if i < len(fieldDataMap)-1 {
			result += ", "
		}
		i++
	}
	result += "}\n"

	fds[lastIdx] = result

	lineBuffer.Reset()
	lineBuffer.AppendString(strings.Join(fds, "\t"))

	return lineBuffer, nil
}

func newEncoder(opt Option) zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder

	if opt.Pretty {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return &coloredConsoleEncoder{Encoder: zapcore.NewConsoleEncoder(cfg)}
	}

	return zapcore.NewJSONEncoder(cfg)
}
