package errors

import (
	origerr "github.com/pkg/errors"
	"testing"
)

var baseErr = origerr.New("Base error")
var userErr = NewClient("Invalid input")
var serverErr = NewServer("Server error")
var wrapErr = Wrap(userErr, "Additional context 1")

func TestErrType(t *testing.T) {
	for _, test := range []struct {
		name         string
		errptr       *err
		wantCause    error
		wantError    string
		wantFile     string
		wantLine     int
		wantMessage  string
		wantPrevious error
		wantUser     error
	}{
		{
			name: "NilErr",
		},
		{
			name:         "WrapperErr",
			errptr:       &err{orig: userErr, msg: "Additional context 2", prev: wrapErr, user: userErr, file: "testFile", line: 13},
			wantCause:    userErr,
			wantError:    "Additional context 2: Additional context 1: Invalid input",
			wantFile:     "testFile",
			wantLine:     13,
			wantMessage:  "Additional context 2",
			wantPrevious: wrapErr,
			wantUser:     userErr,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := test.errptr
			if err.cause() != test.wantCause {
				t.Errorf("Got cause %v; want %v", err.cause(), test.wantCause)
			}
			if err.Error() != test.wantError {
				t.Errorf("Got error '%s'; want '%s'", err.Error(), test.wantError)
			}
			if err.message() != test.wantMessage {
				t.Errorf("Got message '%s'; want '%s'", err.message(), test.wantMessage)
			}
			if err.previous() != test.wantPrevious {
				t.Errorf("Got previous %v; want %v", err.previous(), test.wantPrevious)
			}
			if err.userError() != test.wantUser {
				t.Errorf("Got user %v; want %v", err.userError(), test.wantUser)
			}
			file, line := err.location()
			if file != test.wantFile {
				t.Errorf("Got file %v; want %v", file, test.wantFile)
			}
			if line != test.wantLine {
				t.Errorf("Got line %d; want %d", line, test.wantLine)
			}
		})
	}
}

func TestErrorFunctions(t *testing.T) {
	for _, test := range []struct {
		name            string
		err             error
		wantUserError   error
		wantUserMessage string
		wantUserStatus  int
		wantCause       error
		wantMessage     string
		wantLocation    string
		wantPrevious    error
	}{
		{
			name:           "NilError",
			wantUserStatus: 200,
		},
		{
			name:            "BaseError",
			err:             baseErr,
			wantUserMessage: "Base error",
			wantUserStatus:  500,
			wantCause:       baseErr,
			wantMessage:     "Base error",
			wantLocation:    "Unknown source",
		},
		{
			name:            "WrappedError",
			err:             Wrap(baseErr, "Additional context"),
			wantUserMessage: "Additional context",
			wantUserStatus:  500,
			wantCause:       baseErr,
			wantMessage:     "Additional context",
			wantLocation:    "github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(95)",
			wantPrevious:    baseErr,
		},
		{
			name:            "ClientError",
			err:             userErr,
			wantUserError:   userErr,
			wantUserMessage: "Invalid input",
			wantUserStatus:  400,
			wantCause:       userErr,
			wantMessage:     "Invalid input",
			wantLocation:    "github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(9)",
		},
		{
			name:            "ServerError",
			err:             serverErr,
			wantUserMessage: "Server error",
			wantUserStatus:  500,
			wantCause:       serverErr,
			wantMessage:     "Server error",
			wantLocation:    "github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(10)",
		},
		{
			name:            "DoubleWrappedClientError",
			err:             Wrap(wrapErr, "Additional context 2"),
			wantUserError:   userErr,
			wantUserMessage: "Invalid input",
			wantUserStatus:  400,
			wantCause:       userErr,
			wantMessage:     "Additional context 2",
			wantLocation:    "github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(124)",
			wantPrevious:    wrapErr,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if result := UserError(test.err); result != test.wantUserError {
				t.Errorf("UserError returned %v; want nil", result)
			}
			if message, status := UserDetails(test.err); message != test.wantUserMessage || status != test.wantUserStatus {
				t.Errorf("UserDetails returned (%s, %d); want ('%s', %d)", message, status, test.wantUserMessage, test.wantUserStatus)
			}
			if result := Cause(test.err); result != test.wantCause {
				t.Errorf("Cause returned %v; want %v", result, test.wantCause)
			}
			if result := Message(test.err); result != test.wantMessage {
				t.Errorf("Message returned %v; want '%s'", result, test.wantMessage)
			}
			if result := Location(test.err); result != test.wantLocation {
				t.Errorf("Location returned %v; want '%s'", result, test.wantLocation)
			}
			if result := Previous(test.err); result != test.wantPrevious {
				t.Errorf("Previous returned %v; want %v", result, test.wantPrevious)
			}
		})
	}

	t.Run("WrapNil", func(t *testing.T) {
		if result := Wrap(nil, "Additional context"); result != nil {
			t.Errorf("Wrap returned %v; want nil", result)
		}
	})

	t.Run("Equals", func(t *testing.T) {
		if !Equal(nil, nil) {
			t.Errorf("Nil errors are not equal")
		}

		err1 := Wrap(Wrap(NewServer("Base error"), "Additional context 1"), "Additional context 2")
		err2 := Wrap(Wrap(baseErr, "Additional context 1"), "Additional context 2")
		if !Equal(err1, err2) {
			t.Errorf("Equal returned false for errors with the same annotation stack")
		}

		err3 := Wrap(baseErr, "Additional context 2")
		if Equal(err1, err3) {
			t.Errorf("Equal returned true for errors with different annotation stacks")
		}
	})
}

func TestErrorStack(t *testing.T) {
	t.Run("NilStack", func(t *testing.T) {
		var stack *errorStack = nil
		if stack.hasElements() {
			t.Error("Nil stack has elements")
		}
		stack.push("Nil")
		stack.push("Nil")
		result := stack.pop()
		if result != "" {
			t.Errorf("Got result '%s'; want ''", result)
		}
	})

	t.Run("StackFromError", func(t *testing.T) {
		err1 := Wrap(baseErr, "Message err 1\n")
		err2 := Wrap(err1, "Message err 2\t")
		err3 := Wrap(err2, "Message err 3\r")
		stack := stack(err3)

		if !stack.hasElements() {
			t.Error("Filled stack has no elements")
		}

		wantLines := []string{
			baseErr.Error(),
			"github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(195): Message err 1",
			"github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(196): Message err 2",
			"github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(197): Message err 3",
		}
		for _, wantLine := range wantLines {
			gotLine := stack.pop()
			if gotLine != wantLine {
				t.Errorf("Got line '%s'; want '%s'", gotLine, wantLine)
			}
		}

		wantTrace := baseErr.Error() + "\r\t" +
			"github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(195): Message err 1\r\t" +
			"github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(196): Message err 2\r\t" +
			"github.com/jackstenglein/rest_api_creator/backend/errors/errors_test.go(197): Message err 3\r\t"
		if gotTrace := StackTrace(err3); gotTrace != wantTrace {
			t.Errorf("Got trace:\n%v\nWant trace:\n%v\n", gotTrace, wantTrace)
		}
	})
}
