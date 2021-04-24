package logger

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type (
	ILogger interface {
		Debug(v ...interface{})
		Info(v ...interface{})
		Warning(v ...interface{})
		Error(v ...interface{})
		Critical(v ...interface{})

		// Debugf(format string, v ...interface{})
		// Infof(format string, v ...interface{})
		// Warningf(format string, v ...interface{})
		// Errorf(format string, v ...interface{})
		// Criticalf(format string, v ...interface{})

		// Debugln(v ...interface{})
		// Infoln(v ...interface{})
		// Warningln(v ...interface{})
		// Errorln(v ...interface{})
		// Criticalln(v ...interface{})
	}

	logger struct {
		logger   *log.Logger
		writer   io.Writer
		logLevel int
	}
)

var Log ILogger

func New(writer io.Writer, logLevel int) ILogger {
	return &logger{
		logger:   log.New(writer, "", log.LstdFlags),
		writer:   writer,
		logLevel: logLevel,
	}
}

const (
	critical = iota
	err
	warn
	info
	debug
)

func (l *logger) isEnabledLevel(level int) bool {
	return level <= l.logLevel
}

func getTrace() string {
	if pt, file, line, ok := runtime.Caller(3); ok {
		funcName := runtime.FuncForPC(pt).Name()
		return fmt.Sprintf("%20s:%d | %20s | ", file, line, funcName)
	}
	return ""
}

func (l *logger) print(loglevel int, v ...interface{}) {
	if l.isEnabledLevel(loglevel) {
		now := time.Now().Format("2006/01/02 - 15:04:05")
		trace := getTrace()
		v = append([]interface{}{now, trace}, v...)
		fmt.Fprintln(l.writer, v...)
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

func (l *logger) Debug(v ...interface{})    { l.print(debug, v...) }
func (l *logger) Info(v ...interface{})     { l.print(info, v...) }
func (l *logger) Warning(v ...interface{})  { l.print(warn, v...) }
func (l *logger) Error(v ...interface{})    { l.print(err, v...) }
func (l *logger) Critical(v ...interface{}) { l.print(critical, v...) }

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
