package ray

import "github.com/kataras/golog"

type Logger interface {
	SetLevel(level string)
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	*golog.Logger
}

func NewDefaultLogger() Logger {
	return &logger{
		golog.New(),
	}
}

func (l *logger) SetLevel(level string) {
	l.SetLevel(level)
}
