package main

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

type customLogger struct {
	*log.Logger
}

func (l *customLogger) Debug(v ...interface{}) {
	if l.IsDebug() {
		l.Output(1, DebugprefixIgnorefile, fmt.Sprintln(v...))
	}
}

func (l *customLogger) Debugf(format string, v ...interface{}) {
	if l.IsDebug() {
		l.Output(1, DebugprefixIgnorefile, fmt.Sprintf(format, v...))
	}
}

func NewCustomLogger(out log.FdWriter) *customLogger {
	return &customLogger{log.New(out)}
}

var logger = NewCustomLogger(os.Stdout)
