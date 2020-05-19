// Package errors implements functions to manipulate errors. These functions provide extra
// functionality (beyond the github.com/pkg/errors package) that are helpful in a REST API
// context, including line and file numbers, HTTP status codes and user-friendly messages.
package errors

import "strings"

// annotation wraps the location and previous methods.
// location returns the file name and line number where the annotation was generated.
// previous returns the previous error in the annotation stack.
type annotation interface {
	location() (string, int)
	previous() error
}

// causer wraps the cause method.
// cause returns the original error in the error stack.
type causer interface {
	cause() error
}

// messager wraps the message method.
// message returns a string description of its implementor. This description can be
// more user-friendly than that returned by the Error or String methods.
type messager interface {
	message() string
}

// user wraps the userError method.
// userError returns the original error in the error stack if the implementor is a
// user-caused error. Otherwise, userError returns nil.
type user interface {
	userError() error
}

// err represents an error that is annotated with additional context, file names and
// line numbers, and details of user causes of the error.
//
// err implements the annotation, causer, error, messager and user interfaces. See
// those interfaces for documentation on the methods of type err.
type err struct {
	orig error
	file string
	line int
	msg  string
	prev error
	user error
}

func (e *err) cause() error {
	if e == nil {
		return nil
	}
	return e.orig
}

func (e *err) Error() string {
	if e == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(e.msg)
	err := e.prev
	for err != nil {
		b.WriteString(": ")
		b.WriteString(Message(err))
		err = Previous(err)
	}
	return b.String()
}

func (e *err) location() (string, int) {
	if e == nil {
		return "", 0
	}
	return e.file, e.line
}

func (e *err) message() string {
	if e == nil {
		return ""
	}
	return e.msg
}

func (e *err) previous() error {
	if e == nil {
		return nil
	}
	return e.prev
}

func (e *err) userError() error {
	if e == nil {
		return nil
	}
	return e.user
}
