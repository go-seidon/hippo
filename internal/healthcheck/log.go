package healthcheck

import (
	"fmt"

	gLog "github.com/InVisionApp/go-logger"
	"github.com/go-seidon/provider/logging"
)

type GoHealthLog struct {
	client logging.Logger
}

func (l *GoHealthLog) Info(args ...interface{}) {
	l.client.Info(args...)
}

func (l *GoHealthLog) Debug(args ...interface{}) {
	l.client.Debug(args...)
}

func (l *GoHealthLog) Error(args ...interface{}) {
	l.client.Error(args...)
}

func (l *GoHealthLog) Warn(args ...interface{}) {
	l.client.Warn(args...)
}

func (l *GoHealthLog) Infof(format string, args ...interface{}) {
	l.client.Infof(format, args...)
}

func (l *GoHealthLog) Debugf(format string, args ...interface{}) {
	l.client.Debugf(format, args...)
}

func (l *GoHealthLog) Errorf(format string, args ...interface{}) {
	l.client.Errorf(format, args...)
}

func (l *GoHealthLog) Warnf(format string, args ...interface{}) {
	l.client.Warnf(format, args...)
}

func (l *GoHealthLog) Infoln(args ...interface{}) {
	l.client.Infoln(args...)
}

func (l *GoHealthLog) Debugln(args ...interface{}) {
	l.client.Debugln(args...)
}

func (l *GoHealthLog) Errorln(args ...interface{}) {
	l.client.Errorln(args...)
}

func (l *GoHealthLog) Warnln(args ...interface{}) {
	l.client.Warnln(args...)
}

func (l *GoHealthLog) WithFields(fs gLog.Fields) gLog.Logger {
	client := l.client.WithFields(fs)
	return &GoHealthLog{
		client: client,
	}
}

func NewGoHealthLog(logger logging.Logger) (*GoHealthLog, error) {
	if logger == nil {
		return nil, fmt.Errorf("invalid logger")
	}
	return &GoHealthLog{logger}, nil
}
