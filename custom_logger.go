package main

import (
	"fmt"
	"os"

	"github.com/withmandala/go-log"
	"github.com/withmandala/go-log/colorful"
)

var (
	plainDebug_ignoreFile  = []byte("[DEBUG] ")
	DebugPrefix_ignoreFile = log.Prefix{
		Plain: plainDebug_ignoreFile,
		Color: colorful.Purple(plainDebug_ignoreFile),
		File:  false,
	}
)

type customLogger struct {
	*log.Logger
}

func (l *customLogger) Debug(v ...interface{}) {
	if l.IsDebug() {
		l.Output(1, DebugPrefix_ignoreFile, fmt.Sprintln(v...))
	}
}

func (l *customLogger) Debugf(format string, v ...interface{}) {
	if l.IsDebug() {
		l.Output(1, DebugPrefix_ignoreFile, fmt.Sprintf(format, v...))
	}
}

func NewCustomLogger(out log.FdWriter) *customLogger {
	return &customLogger{log.New(out)}
}

var logger *customLogger = NewCustomLogger(os.Stdout)
