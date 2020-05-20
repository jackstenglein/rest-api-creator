// Package log provides helpers for logging. It currently assumes that logs are written to AWS CloudWatch.
// CloudWatch splits log entries using the newline character (\n), so the log package replaces newlines in
// the middle of strings with carriage returns (\r). Log entries are ended with newlines.
package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

// Log levels control which functions produce output and which are ignored.
const (
	// Silent allows logs from none of the logging functions.
	Silent = iota

	// Failure allows logs from the Fail function and the Error function when used with server errors.
	Failure

	// Warning allows logs from the Warn, Fail and Error functions.
	Warning

	// Debugging allows logs from the Debug, Warn, Fail and Error functions.
	Debugging

	// Information allows logs from all functions.
	Information
)

// level contains the currently in-use log level.
var level = Silent

// writer contains the location to write logs.
var writer io.Writer = os.Stdout

// init uses the DEPLOYMENT_STAGE environment variable to set the log level. The log level can be overridden by using
// the SetLevel function. If DEPLOYMENT_STAGE is not set or is invalid, the log level will be Silent.
func init() {
	switch os.Getenv("DEPLOYMENT_STAGE") {
	case "dev":
		level = Information
	case "alpha":
		level = Debugging
	case "prod":
		level = Failure
	}
}

// printCarriageReturn formats its operands using their default formats. printCarriageReturn then replaces any
// newlines in the resulting string with carriage returns. The final result is written to standard output and
// a newline is appended.
func printCarriageReturn(a ...interface{}) {
	s := fmt.Sprintln(a...)
	s = s[0 : len(s)-1]
	s = strings.Replace(s, "\n", "\r", -1)
	io.WriteString(writer, s)
	io.WriteString(writer, "\n")
}

// SetLevel sets the log level to the provided value. If the given level is invalid, the current log
// level is unchanged.
func SetLevel(l int) {
	if l >= Silent && l <= Information {
		level = l
	}
}

// Error writes the given error to standard output. If the error is nil, nothing is logged. If the error is a server error,
// the log level must be Failure or higher for the error to be logged. If the error is a client error,
// the log level must be Warning or higher for the error to be logged.
func Error(err error) {
	_, status := errors.UserDetails(err)
	if status == 200 {
		return
	}

	stack := errors.StackTrace(err)
	if status == 500 {
		Fail(stack)
	} else if status == 400 {
		Warn(stack)
	}
}

// Fail formats using the default formats for its operands and writes to standard output if the log level is greater than or equal to
// Failure. Spaces are added between operands.
func Fail(a ...interface{}) {
	if level >= Failure {
		io.WriteString(writer, "[FAIL]: ")
		printCarriageReturn(a...)
	}
}

// Warn formats using the default formats for its operands and writes to standard output if the log level is greater than or equal to
// Warning. Spaces are added between operands.
func Warn(a ...interface{}) {
	if level >= Warning {
		io.WriteString(writer, "[WARN]: ")
		printCarriageReturn(a...)
	}
}

// Debug formats using the default formats for its operands and writes to standard output if the log level is greater than or equal to
// Debugging. Spaces are added between operands.
func Debug(a ...interface{}) {
	if level >= Debugging {
		io.WriteString(writer, "[DEBUG]: ")
		printCarriageReturn(a...)
	}
}

// Info formats using the default formats for its operands and writes to standard output if the log level is greater than or equal to
// Information. Spaces are added between operands.
func Info(a ...interface{}) {
	if level >= Information {
		io.WriteString(writer, "[INFO]: ")
		printCarriageReturn(a...)
	}
}
