package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-stack/stack"
	"github.com/sirupsen/logrus"
)

type GoFormatter struct {
	PrettyPrint bool
	StackSkip   []string
}

func (f *GoFormatter) reportLocation() (stack.Call, error) {
	skip := func(pkg string) bool {
		for _, skip := range f.StackSkip {
			if pkg == skip {
				return true
			}
		}
		return false
	}

	// We start at 2 to skip this call and our caller's call.
	for i := 2; ; i++ {
		c := stack.Caller(i)
		// ErrNoFunc indicates we're over traversing the stack.
		if _, err := c.MarshalText(); err != nil {
			return stack.Call{}, nil
		}

		pkg := fmt.Sprintf("%+k", c)
		// Remove vendoring from package path.
		parts := strings.SplitN(pkg, "/vendor/", 2)
		pkg = parts[len(parts)-1]
		if !skip(pkg) {
			return c, nil
		}
	}
}

func (f *GoFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := entry.Data

	message := entry.Message
	logErr, ok := data[FIELD_ERROR]
	if ok {
		message = fmt.Sprintf("%s: %s", message, logErr)
		delete(data, FIELD_ERROR)
	}

	service := data[FIELD_SERVICE]
	delete(data, FIELD_SERVICE)

	reportLoc := struct {
		FilePath     string `json:"filePath,omitempty"`
		LineNumber   int    `json:"lineNumber,omitempty"`
		FunctionName string `json:"functionName,omitempty"`
	}{}
	c, err := f.reportLocation()
	if err == nil {
		lineNumber, _ := strconv.ParseInt(fmt.Sprintf("%d", c), 10, 64)
		reportLoc.FilePath = fmt.Sprintf("%+s", c)
		reportLoc.LineNumber = int(lineNumber)
		reportLoc.FunctionName = fmt.Sprintf("%n", c)
	}

	logMessage := LogMessage{
		Timestamp:      entry.Time.UTC().Format(time.RFC3339),
		Message:        message,
		Severity:       entry.Level.String(),
		Metadata:       data,
		Service:        service,
		ReportLocation: reportLoc,
	}

	if len(data) > 0 {
		logMessage.Metadata = data
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	err = encoder.Encode(logMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return b.Bytes(), nil
}
