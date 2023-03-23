package mlog

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	DefaultLogger *logrus.Logger
)

func init() {
	DefaultLogger = logrus.New()
	l, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		l = logrus.InfoLevel
	}
	DefaultLogger.Level = l
	DefaultLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
}

func Trace(a ...any)                 { DefaultLogger.Trace(a...) }
func Debug(a ...any)                 { DefaultLogger.Debug(a...) }
func Info(a ...any)                  { DefaultLogger.Info(a...) }
func Warn(a ...any)                  { DefaultLogger.Warn(a...) }
func Error(a ...any)                 { DefaultLogger.Error(a...) }
func Fatal(a ...any)                 { DefaultLogger.Fatal(a...) }
func Tracef(format string, a ...any) { DefaultLogger.Tracef(format, a...) }
func Debugf(format string, a ...any) { DefaultLogger.Debugf(format, a...) }
func Infof(format string, a ...any)  { DefaultLogger.Infof(format, a...) }
func Warnf(format string, a ...any)  { DefaultLogger.Warnf(format, a...) }
func Errorf(format string, a ...any) { DefaultLogger.Errorf(format, a...) }
func Fatalf(format string, a ...any) { DefaultLogger.Fatalf(format, a...) }
