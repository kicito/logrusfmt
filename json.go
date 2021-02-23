package logrusfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var ProductionFormatter logrus.Formatter = &jsonFormatter{
	formatter: &logrus.JSONFormatter{
		DataKey:          "message",
		CallerPrettyfier: stackPrint,
	},
}

type jsonFormatter struct {
	formatter *logrus.JSONFormatter
}

func (f *jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+4)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	if f.formatter.DataKey != "" {
		temp := make(logrus.Fields, 4)
		data["message"] = entry.Message
		temp[f.formatter.DataKey] = data
		data = temp
	}

	if entry.Context != nil {
		setRequestContext(entry.Context, data)
	}

	timestampFormat := f.formatter.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	if !f.formatter.DisableTimestamp {
		data["timestamp"] = entry.Time.Format(timestampFormat)
	}
	data["severity"] = entry.Level.String()
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		if f.formatter.CallerPrettyfier != nil {
			funcVal, _ = f.formatter.CallerPrettyfier(entry.Caller)
		}
		if funcVal != "" {
			data["trace"] = funcVal
		}
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.formatter.DisableHTMLEscape)
	if f.formatter.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	return b.Bytes(), nil
}

func stackPrint(f *runtime.Frame) (function string, file string) {
	//s := strings.Split(f.Function, ".")
	//fn := s[:len(s)-1]
	//funcName := strings.Join(fn, ".")
	function = fmt.Sprintf("in %s() at %s:%v ", f.Func.Name(), f.File, f.Line)
	return function, ""
}
