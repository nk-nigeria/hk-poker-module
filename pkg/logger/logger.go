package logger

import (
	"io"
	"sync"

	"github.com/evalphobia/logrus_sentry"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var instance *Logger
var once sync.Once

// Logrus : implement Logger
type Logger struct {
	log *logrus.Logger
}

func getLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			log: logrus.StandardLogger(),
		}
		instance.log.SetFormatter(new(logrus.JSONFormatter))
	})
	return instance
}

// GetEchoLogger for e.Logger
func GetEchoLogger() Logger {
	return Logger{
		log: getLogger().log,
	}
}

func AddSentry(dns string) {
	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
	hook, err := logrus_sentry.NewSentryHook(dns, levels)
	if err != nil {
		getLogger().Errorf("new sentry hook err: %v", err)
		return
	}

	hook.Timeout = 0

	hook.StacktraceConfiguration = logrus_sentry.StackTraceConfiguration{
		Enable:            true,
		Level:             logrus.ErrorLevel,
		Skip:              8,
		Context:           7,
		InAppPrefixes:     nil,
		SendExceptionType: true,
	}

	getLogger().log.Hooks.Add(hook)

}

func Print(i ...interface{}) {
	getLogger().Print(i...)
}
func Printf(format string, args ...interface{}) {
	getLogger().Printf(format, args...)
}
func Debug(i ...interface{}) {
	getLogger().Debug(i...)
}
func Debugf(format string, args ...interface{}) {
	getLogger().Debugf(format, args...)
}
func Info(i ...interface{}) {
	getLogger().Info(i...)
}
func Infof(format string, args ...interface{}) {
	getLogger().Infof(format, args...)
}
func Warn(i ...interface{}) {
	getLogger().Warn(i...)
}
func Warnf(format string, args ...interface{}) {
	getLogger().Warnf(format, args...)
}
func Warning(i ...interface{}) {
	getLogger().Warn(i...)
}
func Warningf(format string, args ...interface{}) {
	getLogger().Warnf(format, args...)
}
func Error(i ...interface{}) {
	getLogger().Error(i...)
}
func Errorf(format string, args ...interface{}) {
	getLogger().Errorf(format, args...)
}
func Fatal(i ...interface{}) {
	getLogger().Fatal(i...)
}
func Fatalf(format string, args ...interface{}) {
	getLogger().Fatalf(format, args...)
}
func Panic(i ...interface{}) {
	getLogger().Panic(i...)
}
func Panicf(format string, args ...interface{}) {
	getLogger().Panicf(format, args...)
}

func WithField(key string, value interface{}) *Logger {
	getLogger().log.WithField(key, value)
	return getLogger()
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.log.WithField(key, value)
	return l
}

// implement Logger Echo interface
func (l Logger) Output() io.Writer {
	return getLogger().log.Out
}
func (l Logger) SetOutput(w io.Writer) {
	getLogger().SetOutput(w)
}
func (l Logger) Prefix() string {
	return ""
}

func (l Logger) SetPrefix(p string) {

}
func (l Logger) Level() log.Lvl {
	switch getLogger().log.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		getLogger().Panic("Invalid level")
	}

	return log.OFF
}
func (l Logger) SetLevel(v log.Lvl) {
	switch v {
	case log.DEBUG:
		getLogger().log.SetLevel(logrus.DebugLevel)
	case log.WARN:
		getLogger().log.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		getLogger().log.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		getLogger().log.SetLevel(logrus.InfoLevel)
	default:
		getLogger().log.Panic("Invalid level")
	}
}
func (l Logger) SetHeader(h string) {}

func (l Logger) Print(i ...interface{}) {
	l.log.Print(i...)
}
func (l Logger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}
func (l Logger) Printj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Print()
}
func (l Logger) Debug(i ...interface{}) {
	l.log.Debug(i...)
}
func (l Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}
func (l Logger) Debugj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Debug()
}
func (l Logger) Info(i ...interface{}) {
	l.log.Info(i...)
}
func (l Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}
func (l Logger) Infoj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Info()
}
func (l Logger) Warn(i ...interface{}) {
	l.log.Warn(i...)
}
func (l Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}
func (l Logger) Warnj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Warn()
}
func (l Logger) Error(i ...interface{}) {
	l.log.Error(i...)
}
func (l Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
func (l Logger) Errorj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Error()
}
func (l Logger) Fatal(i ...interface{}) {
	l.log.Fatal(i...)
}
func (l Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}
func (l Logger) Fatalj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Fatal()
}
func (l Logger) Panic(i ...interface{}) {
	l.log.Panic(i...)
}
func (l Logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}
func (l Logger) Panicj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Panic()
}
