package logger

import (
	"fmt"
	"os"

	"github.com/withmandala/go-log"
	"github.com/withmandala/go-log/colorful"
)

var (
	plaindebugIgnorefile  = []byte("[DEBUG] ")
	DebugprefixIgnorefile = log.Prefix{
		Plain: plaindebugIgnorefile,
		Color: colorful.Purple(plaindebugIgnorefile),
		File:  false,
	}
)

type CustomLogger struct {
	*log.Logger
}

func (l *CustomLogger) Debug(v ...interface{}) {
	if l.IsDebug() {
		l.Output(1, DebugprefixIgnorefile, fmt.Sprintln(v...))
	}
}

func (l *CustomLogger) Debugf(format string, v ...interface{}) {
	if l.IsDebug() {
		l.Output(1, DebugprefixIgnorefile, fmt.Sprintf(format, v...))
	}
}

func NewCustomLogger(out log.FdWriter) *CustomLogger {
	return &CustomLogger{log.New(out)}
}

var Logger = NewCustomLogger(os.Stdout)
