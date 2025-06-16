package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
)

var (
	levelStrings = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	levelColors  = []string{cyan, green, yellow, red, purple}
)

type Logger struct {
	*log.Logger
	mu       sync.Mutex
	level    LogLevel
	output   io.Writer
	colorful bool
}

var std = New(os.Stdout, INFO, true)

func New(out io.Writer, level LogLevel, colorful bool) *Logger {
	return &Logger{
		Logger:   log.New(out, "", log.LstdFlags|log.Lshortfile),
		level:    level,
		output:   out,
		colorful: colorful,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
	l.Logger = log.New(l.output, "", log.LstdFlags|log.Lshortfile)
}

func (l *Logger) log(level LogLevel, depth int, msg string) {
	if level < l.level {
		return
	}

	if l.colorful {
		color := levelColors[level]
		// 整行使用相同颜色
		l.Output(depth+1, fmt.Sprintf("%s%s%s", color, msg, reset))
	} else {
		l.Output(depth+1, msg)
	}
}

// Package-level functions
func SetLevel(level LogLevel) { std.SetLevel(level) }
func SetOutput(w io.Writer)   { std.SetOutput(w) }
func SetColorful(b bool)      { std.colorful = b }

func Debug(v ...interface{}) {
	msg := formatLog(DEBUG, 2, fmt.Sprint(v...))
	std.log(DEBUG, 2, msg)
}

func Debugf(format string, v ...interface{}) {
	msg := formatLog(DEBUG, 2, fmt.Sprintf(format, v...))
	std.log(DEBUG, 2, msg)
}

func Info(v ...interface{}) {
	msg := formatLog(INFO, 2, fmt.Sprint(v...))
	std.log(INFO, 2, msg)
}

func Infof(format string, v ...interface{}) {
	msg := formatLog(INFO, 2, fmt.Sprintf(format, v...))
	std.log(INFO, 2, msg)
}

func Warn(v ...interface{}) {
	msg := formatLog(WARNING, 2, fmt.Sprint(v...))
	std.log(WARNING, 2, msg)
}

func Warnf(format string, v ...interface{}) {
	msg := formatLog(WARNING, 2, fmt.Sprintf(format, v...))
	std.log(WARNING, 2, msg)
}

func Error(v ...interface{}) {
	msg := formatLog(ERROR, 2, fmt.Sprint(v...))
	std.log(ERROR, 2, msg)
}

func Errorf(format string, v ...interface{}) {
	msg := formatLog(ERROR, 2, fmt.Sprintf(format, v...))
	std.log(ERROR, 2, msg)
}

func Fatal(v ...interface{}) {
	msg := formatLog(FATAL, 2, fmt.Sprint(v...))
	std.log(FATAL, 2, msg)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	msg := formatLog(FATAL, 2, fmt.Sprintf(format, v...))
	std.log(FATAL, 2, msg)
	os.Exit(1)
}

func formatLog(level LogLevel, depth int, msg string) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	file = filepath.Base(file)
	return fmt.Sprintf("[%s] %s:%d: %s", levelStrings[level], file, line, msg)
}
