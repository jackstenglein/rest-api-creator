package errors

import "strings"

type annotation interface {
	location() (string, int)
	previous() error
}

type causer interface {
	cause() error
}

type messager interface {
	message() string
}

type user interface {
	userError() error
}

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
