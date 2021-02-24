package logrusfmt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type textFormatter struct {
	txtfmt *logrus.TextFormatter
}

func (t *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b = new(bytes.Buffer)
	var levelColor int
	var (
		green  = 32
		gray   = 90
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
	if levelText == "WARNING" {
		levelText = "WARN"
	}

	//if entry.HasCaller() {

	contextData := make(map[string]interface{})

	fmt.Fprintf(b, "[%s] %-16s %s %-44s", colorText(gray, entry.Time.Format("15:04:05 MST")), fmt.Sprintf("[%s]",colorText(levelColor,levelText)), colorText(gray, "-"), entry.Message)
	if len(entry.Data) > 0 {
		for k, v := range entry.Data {
			fmt.Fprintf(b, "%s %v ", colorText(blue, k+":"), v)
		}
	}
	//b.WriteByte('\n')
	//b.WriteByte('\t')
	if entry.HasCaller() {
		b.WriteByte('\n')
		b.WriteByte('\t')
		c := entry.Caller
		fmt.Fprintf(b, "in %s at %s", colorText(blue, c.Function+"()"), colorText(green, fmt.Sprintf("%s:%d", c.File, c.Line)))
		if entry.Context != nil {
			setRequestContext(entry.Context, contextData)
		}
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

func colorText(color int, text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
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
