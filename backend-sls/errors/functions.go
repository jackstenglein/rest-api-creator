package errors

import (
	"fmt"
	"runtime"
	"strings"
)

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

func NewUserError(message string) error {
	err := &err{
		msg:  message,
		prev: nil,
	}
	err.orig = err
	err.user = err
	setLocation(err)
	return err
}

func UserError(err error) error {
	if uerr, ok := err.(user); ok {
		return uerr.userError()
	}
	return nil
}

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

func Cause(err error) error {
	if cerr, ok := err.(causer); ok {
		return cerr.cause()
	}
	return err
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	if merr, ok := err.(messager); ok {
		return merr.message()
	}
	return err.Error()
}

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

func Previous(err error) error {
	if aerr, ok := err.(annotation); ok {
		return aerr.previous()
	}
	return nil
}

func StackTrace(err error) string {
	var b strings.Builder
	errStack := stack(err)
	for errStack.hasElements() {
		b.WriteString(errStack.pop())
		b.WriteString("\r\t")
	}
	return b.String()
}
