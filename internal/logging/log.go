package logging

import (
	"context"
	"io"
)

const (
	FIELD_SERVICE = "service"
	FIELD_ERROR   = "error"
)

type Logger interface {
	SimpleLog
	FormatedLog
	LineLog
	CustomLog
}

type SimpleLog interface {
	Log(level string, args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

type FormatedLog interface {
	Logf(level string, format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type LineLog interface {
	Logln(level string, msg ...interface{})
	Infoln(msg ...interface{})
	Debugln(msg ...interface{})
	Errorln(msg ...interface{})
	Warnln(msg ...interface{})
}

type CustomLog interface {
	WithFields(fs map[string]interface{}) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger
	WriterLevel(level string) io.Writer
}

type LogMessage struct {
	Timestamp      string `json:"timestamp"`
	Message        string `json:"message"`
	Severity       string `json:"severity"`
	ReportLocation struct {
		FilePath     string `json:"filePath,omitempty"`
		LineNumber   int    `json:"lineNumber,omitempty"`
		FunctionName string `json:"functionName,omitempty"`
	} `json:"reportLocation,omitempty"`
	Service  interface{}            `json:"service,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type LogParam struct {
	AppCtxEnabled bool
	AppName       string
	AppVersion    string

	DebuggingEnabled   bool
	PrettyPrintEnabled bool

	StackSkip []string
}

type LogOption func(*LogParam)

func WithAppContext(name, version string) LogOption {
	return func(lo *LogParam) {
		lo.AppCtxEnabled = true
		lo.AppName = name
		lo.AppVersion = version
	}
}

func EnableDebugging() LogOption {
	return func(lo *LogParam) {
		lo.DebuggingEnabled = true
	}
}

func EnablePrettyPrint() LogOption {
	return func(lo *LogParam) {
		lo.PrettyPrintEnabled = true
	}
}

func AddStackSkip(pkg string) LogOption {
	return func(lo *LogParam) {
		lo.StackSkip = append(lo.StackSkip, pkg)
	}
}
