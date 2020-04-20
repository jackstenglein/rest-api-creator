package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// setLocation sets the file and line number of err to the point of err's generation.
// The file path is truncated to start at the module root directory.
func setLocation(err *err) {
	const leadingPath = "github.com/"
	const leadingLen = len(leadingPath)
	if pc, file, line, ok := runtime.Caller(2); ok {
		funcName := runtime.FuncForPC(pc).Name()
		endIndex := strings.Index(funcName[leadingLen:], ".") + leadingLen
		packageName := funcName[:endIndex]

		startIndex := strings.LastIndex(file, "/")
		file = file[startIndex:]
		err.file = packageName + file
		err.line = line
	}
}

// Wrap returns an error annotating previous with the supplied message and the file and line
// number of the point where Wrap was called. If previous is nil, Wrap returns nil.
func Wrap(previous error, message string) error {
	if previous == nil {
		return nil
	}
	err := &err{
		msg:  message,
		prev: previous,
		orig: Cause(previous),
		user: UserError(previous),
	}
	setLocation(err)
	return err
}

// NewClient returns a client-caused error with the supplied message. The error is annotated
// with the file and line number of the point where NewClient was called.
func NewClient(message string) error {
	err := &err{
		msg:  message,
		prev: nil,
	}
	err.orig = err
	err.user = err
	setLocation(err)
	return err
}

// NewServer returns a server-caused error with the supplied message. The error is annotated
// with the file and line number of the point where NewServer was called.
func NewServer(message string) error {
	err := &err{
		msg:  message,
		prev: nil,
	}
	err.orig = err
	setLocation(err)
	return err
}

// UserError returns the original error if err is a user-caused error. If err is not user-caused,
// UserError returns nil.
func UserError(err error) error {
	if uerr, ok := err.(user); ok {
		return uerr.userError()
	}
	return nil
}

// UserDetails returns the user-facing message and HTTP status code associated with err. If err
// is nil, the empty string and status code 200 are returned.
func UserDetails(err error) (string, int) {
	if err == nil {
		return "", 200
	}
	if uerr, ok := err.(user); ok {
		if underlying := uerr.userError(); underlying != nil {
			return Message(underlying), 400
		}
	}
	return Message(err), 500
}

// Cause returns the original error that led to err. If err does not implement the cause() method,
// err is assumed to be the original error and is returned.
func Cause(err error) error {
	if cerr, ok := err.(causer); ok {
		return cerr.cause()
	}
	return err
}

// Message returns a string description of err that may be more user-friendly than err.Error(). If
// err does not implement the message() method, Message returns err.Error().
func Message(err error) string {
	if err == nil {
		return ""
	}
	if merr, ok := err.(messager); ok {
		return merr.message()
	}
	return err.Error()
}

// Location returns a user-friendly string description of the file and line number where err was
// generated. If err does not implement the location() method, Location returns the string 'Unknown source'.
func Location(err error) string {
	if err == nil {
		return ""
	}
	if aerr, ok := err.(annotation); ok {
		file, line := aerr.location()
		if len(file) > 0 && line > 0 {
			return fmt.Sprintf("%s(%d)", file, line)
		}
	}
	return "Unknown source"
}

// Previous returns the direct ancestor of err in the error stack. If err does not implement the previous()
// method, Previous returns nil.
func Previous(err error) error {
	if aerr, ok := err.(annotation); ok {
		return aerr.previous()
	}
	return nil
}

// StackTrace returns a string description of the error stack described by err. The original cause of err
// is listed first. Each subsequent error is preceded by the file and line number of the point it was generated.
// Each error description is separated by a carriage return ('\r') and error descriptions after the original are
// preceded by a tab ('\t').
func StackTrace(err error) string {
	var b strings.Builder
	errStack := stack(err)
	for errStack.hasElements() {
		b.WriteString(errStack.pop())
		b.WriteString("\r\t")
	}
	return b.String()
}

// Equal returns true only if all errors in lhs's annotation stack have the same messages as the corresponding
// errors in rhs's error stack. File names and line numbers of the annotations are ignored. This function is
// intended to be used by tests in order to check returned error values.
func Equal(lhs error, rhs error) bool {
	for lhs != nil && rhs != nil {
		if Message(lhs) != Message(rhs) {
			return false
		}
		lhs = Previous(lhs)
		rhs = Previous(rhs)
	}
	return lhs == nil && rhs == nil
}
