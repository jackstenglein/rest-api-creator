package errors

import (
	"fmt"
	"runtime"
	"strings"
)

func setLocation(err *Err) {
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

func Wrap(previous error, message string) error {
	if previous == nil {
		return nil
	}
	err := &Err{
		message:  message,
		previous: previous,
		cause:    Cause(previous),
		user:     UserError(previous),
	}
	setLocation(err)
	return err
}

func NewUserError(message string) error {
	err := &Err{
		message:  message,
		previous: nil,
	}
	err.cause = err
	err.user = err
	setLocation(err)
	return err
}

func UserError(err error) error {
	if uerr, ok := err.(user); ok {
		return uerr.UserError()
	}
	return nil
}

func UserDetails(err error) (string, int) {
	if err == nil {
		return "", 200
	}
	if uerr, ok := err.(user); ok {
		if underlying := uerr.UserError(); underlying != nil {
			return Message(underlying), 400
		}
	}
	return Message(err), 500
}

func Cause(err error) error {
	if cerr, ok := err.(causer); ok {
		return cerr.Cause()
	}
	return err
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	if merr, ok := err.(messager); ok {
		return merr.Message()
	}
	return err.Error()
}

func Location(err error) string {
	if err == nil {
		return ""
	}
	if aerr, ok := err.(annotation); ok {
		file, line := aerr.Location()
		if len(file) > 0 && line > 0 {
			return fmt.Sprintf("%s(%d)", file, line)
		}
	}
	return "Unknown source"
}

func Previous(err error) error {
	if aerr, ok := err.(annotation); ok {
		return aerr.Previous()
	}
	return nil
}

func StackTrace(err error) string {
	var b strings.Builder
	errStack := Stack(err)
	for errStack.HasElements() {
		b.WriteString(errStack.Pop())
		b.WriteString("\r\t")
	}
	return b.String()
}
