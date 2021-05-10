package logger

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

type (
	ILogger interface {
		Debug(v ...interface{})
		Info(v ...interface{})
		Warning(v ...interface{})
		Error(v ...interface{})
		Critical(v ...interface{})
	}

	logger struct {
		logger   *log.Logger
		writer   io.Writer
		logLevel LogLevel
	}
)

var Log ILogger

func New(writer io.Writer, logLevel LogLevel, prefix ...string) ILogger {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return &logger{
		logger:   log.New(writer, p, log.LstdFlags),
		writer:   writer,
		logLevel: logLevel,
	}
}

type LogLevel int

func Convert(level string) LogLevel {
	switch strings.ToLower(level) {
	case "critical":
		return Critical
	case "error", "err":
		return Error
	case "warn", "warnign":
		return Warn
	case "info", "prod":
		return Info
	case "debug", "dev":
		return Debug
	default:
		return Info
	}
}

const (
	Critical LogLevel = iota
	Error
	Warn
	Info
	Debug
)

func (l *logger) isEnabledLevel(level LogLevel) bool {
	return level <= l.logLevel
}

func getTrace() string {
	if pt, file, line, ok := runtime.Caller(3); ok {
		funcName := runtime.FuncForPC(pt).Name()
		return fmt.Sprintf("%s - %d\tfunc:%s", file, line, funcName)
	}
	return ""
}

func (l *logger) print(loglevel LogLevel, v ...interface{}) {
	if l.isEnabledLevel(loglevel) {
		now := time.Now().Format("2006/01/02 - 15:04:05")
		// trace := getTrace()
		print := fmt.Sprint(v...)
		prints := strings.Split(print, "\n")
		for _, p := range prints {
			if strings.TrimSpace(p) != "" {
				values := []interface{}{now, p}
				value := fmt.Sprintf("time:%s\tmessage:%v", values...)
				fmt.Fprintln(l.writer, value)
			}
		}
		// l.logger.Print(v...)
	}
}

// func (l *logger) printf(loglevel int, format string, v ...interface{}) {
// 	if l.isEnabledLevel(loglevel) {
// 		l.logger.Printf(format, v...)
// 	}
// }

// func (l *logger) println(loglevel int, v ...interface{}) {
// 	if l.isEnabledLevel(loglevel) {
// 		l.logger.Println(v...)
// 	}
// }

func (l *logger) Debug(v ...interface{})    { l.print(Debug, v...) }
func (l *logger) Info(v ...interface{})     { l.print(Info, v...) }
func (l *logger) Warning(v ...interface{})  { l.print(Warn, v...) }
func (l *logger) Error(v ...interface{})    { l.print(Error, v...) }
func (l *logger) Critical(v ...interface{}) { l.print(Critical, v...) }

// func (l *logger) Debugf(format string, v ...interface{})    { l.printf(debug, format, v...) }
// func (l *logger) Infof(format string, v ...interface{})     { l.printf(info, format, v...) }
// func (l *logger) Warningf(format string, v ...interface{})  { l.printf(warn, format, v...) }
// func (l *logger) Errorf(format string, v ...interface{})    { l.printf(err, format, v...) }
// func (l *logger) Criticalf(format string, v ...interface{}) { l.printf(critical, format, v...) }

// func (l *logger) Debugln(v ...interface{})    { l.println(debug, v...) }
// func (l *logger) Infoln(v ...interface{})     { l.println(info, v...) }
// func (l *logger) Warningln(v ...interface{})  { l.println(warn, v...) }
// func (l *logger) Errorln(v ...interface{})    { l.println(err, v...) }
// func (l *logger) Criticalln(v ...interface{}) { l.println(critical, v...) }
