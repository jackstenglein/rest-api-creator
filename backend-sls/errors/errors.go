package errors

type ApiError interface {
	Error() string
	StatusCode() int
}

type UserError struct {
	message string
}

type ServerError struct {
	message string
}

func NewUserError(message string) ApiError {
	return &UserError{message}
}

func (e *UserError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *UserError) StatusCode() int {
	if e == nil {
		return 200
	}
	return 400
}

func NewServerError(message string) ApiError {
	return &ServerError{message}
}

func (e *ServerError) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

func (e *ServerError) StatusCode() int {
	if e == nil {
		return 200
	}
	return 500
}
