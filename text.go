package logrusfmt

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type textFormatter struct {
	txtfmt *logrus.TextFormatter
}

func (t *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b = new(bytes.Buffer)
	var levelColor int
	var (
		gray   = 37
		yellow = 33
		red    = 31
		blue   = 36
	)
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	levelText := strings.ToUpper(entry.Level.String())

	var caller string
	if entry.HasCaller() {
		if t.txtfmt.CallerPrettyfier != nil {
			caller, _ = t.txtfmt.CallerPrettyfier(entry.Caller)
		}
	}
	contextData := make(map[string]interface{})
	if entry.Context != nil {
		setRequestContext(entry.Context, contextData)
	}

	fmt.Fprintf(b, "[\x1b[%dm%s\x1b[0m] [\x1b[%dm%s\x1b[0m] %-44s", gray, entry.Time.Format("2006-01-02 15:04:05 MST"), levelColor, levelText, entry.Message)
	if len(entry.Data) > 0 {
		b.WriteByte('\n')
		b.WriteByte('\t')
		for k, v := range entry.Data {
			fmt.Fprintf(b, "\x1b[%dm%s:\x1b[0m %v ", blue, k, v)
		}
	}
	if entry.HasCaller() {
		b.WriteByte('\n')
		b.WriteByte('\t')
		fmt.Fprintf(b, "%s", caller)
	}
	if len(contextData) > 0 {
		b.WriteByte('\n')
		b.WriteByte('\t')
		for k, v := range contextData {
			fmt.Fprintf(b, "\x1b[%dm%s:\x1b[0m %v ", blue, k, v)
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

var LocalFormatter logrus.Formatter = &textFormatter{
	txtfmt: &logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          stackPrint,
	},
}
