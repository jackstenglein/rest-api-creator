package errors

type UserError struct {
	message string
}

func NewUserError(message string) error {
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
