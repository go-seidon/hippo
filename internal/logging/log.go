package logging

import "context"

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
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

type FormatedLog interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type LineLog interface {
	Infoln(msg ...interface{})
	Debugln(msg ...interface{})
	Errorln(msg ...interface{})
	Warnln(msg ...interface{})
}

type CustomLog interface {
	WithFields(fs map[string]interface{}) Logger
	WithError(err error) Logger
	WithContext(ctx context.Context) Logger
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

type LogOption struct {
	AppCtxEnabled bool
	AppName       string
	AppVersion    string

	DebuggingEnabled   bool
	PrettyPrintEnabled bool

	StackSkip []string
}

type Option func(*LogOption)

func WithAppContext(name, version string) Option {
	return func(lo *LogOption) {
		lo.AppCtxEnabled = true
		lo.AppName = name
		lo.AppVersion = version
	}
}

func EnableDebugging() Option {
	return func(lo *LogOption) {
		lo.DebuggingEnabled = true
	}
}

func EnablePrettyPrint() Option {
	return func(lo *LogOption) {
		lo.PrettyPrintEnabled = true
	}
}

func AddStackSkip(pkg string) Option {
	return func(lo *LogOption) {
		lo.StackSkip = append(lo.StackSkip, pkg)
	}
}
