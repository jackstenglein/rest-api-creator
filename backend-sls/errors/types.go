package errors

import "strings"

type annotation interface {
	Location() (string, int)
	Previous() error
}

type causer interface {
	Cause() error
}

type messager interface {
	Message() string
}

type user interface {
	UserError() error
}

type Err struct {
	cause    error
	file     string
	line     int
	message  string
	previous error
	user     error
}

func (e *Err) Cause() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Err) Error() string {
	if e == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(e.message)
	err := e.previous
	for err != nil {
		b.WriteString(": ")
		b.WriteString(Message(err))
		err = Previous(err)
	}
	return b.String()
}

func (e *Err) Location() (string, int) {
	if e == nil {
		return "", 0
	}
	return e.file, e.line
}

func (e *Err) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *Err) Previous() error {
	if e == nil {
		return nil
	}
	return e.previous
}

func (e *Err) UserError() error {
	if e == nil {
		return nil
	}
	return e.user
}
