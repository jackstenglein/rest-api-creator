package log

import (
	"github.com/jackstenglein/rest_api_creator/backend/errors"
	"strings"
	"testing"
)

type errorTest struct {
	name       string
	err        error
	wantString string
}

type logTest struct {
	name       string
	action     func(a ...interface{})
	args       []interface{}
	wantString string
}

var logTests = []struct {
	name       string
	level      int
	errorTests []errorTest
	logTests   []logTest
}{
	{
		name:  "Silent",
		level: Silent,
		errorTests: []errorTest{
			{
				name:       "ServerError",
				err:        errors.NewServer("ServerError"),
				wantString: "",
			},
			{
				name:       "ClientError",
				err:        errors.NewClient("ClientError"),
				wantString: "",
			},
		},
		logTests: []logTest{
			{
				name:       "Fail",
				action:     Fail,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Warn",
				action:     Warn,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Debug",
				action:     Debug,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Info",
				action:     Info,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
		},
	},
	{
		name:  "Failure",
		level: Failure,
		errorTests: []errorTest{
			{
				name:       "ServerError",
				err:        errors.NewServer("ServerError"),
				wantString: "[FAIL]: ServerError\r\t\n",
			},
			{
				name:       "ClientError",
				err:        errors.NewClient("ClientError"),
				wantString: "",
			},
		},
		logTests: []logTest{
			{
				name:       "Fail",
				action:     Fail,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[FAIL]: This is a test\n",
			},
			{
				name:       "Warn",
				action:     Warn,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Debug",
				action:     Debug,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Info",
				action:     Info,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
		},
	},
	{
		name:  "Warning",
		level: Warning,
		errorTests: []errorTest{
			{
				name:       "ServerError",
				err:        errors.NewServer("Server\nError"),
				wantString: "[FAIL]: Server\rError\r\t\n",
			},
			{
				name:       "ClientError",
				err:        errors.NewClient("Client\nError"),
				wantString: "[WARN]: Client\rError\r\t\n",
			},
		},
		logTests: []logTest{
			{
				name:       "Fail",
				action:     Fail,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[FAIL]: This is a test\n",
			},
			{
				name:       "Warn",
				action:     Warn,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[WARN]: This is a test\n",
			},
			{
				name:       "Debug",
				action:     Debug,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Info",
				action:     Info,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
		},
	},
	{
		name:  "Debugging",
		level: Debugging,
		errorTests: []errorTest{
			{
				name:       "ServerError",
				err:        errors.NewServer("Server\nError"),
				wantString: "[FAIL]: Server\rError\r\t\n",
			},
			{
				name:       "ClientError",
				err:        errors.NewClient("Client\nError"),
				wantString: "[WARN]: Client\rError\r\t\n",
			},
		},
		logTests: []logTest{
			{
				name:       "Fail",
				action:     Fail,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[FAIL]: This is a test\n",
			},
			{
				name:       "Warn",
				action:     Warn,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[WARN]: This is a test\n",
			},
			{
				name:       "Debug",
				action:     Debug,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[DEBUG]: This is a test\n",
			},
			{
				name:       "Info",
				action:     Info,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
		},
	},
	{
		name:  "Information",
		level: Information,
		errorTests: []errorTest{
			{
				name:       "ServerError",
				err:        errors.NewServer("Server\nError"),
				wantString: "[FAIL]: Server\rError\r\t\n",
			},
			{
				name:       "ClientError",
				err:        errors.NewClient("Client\nError"),
				wantString: "[WARN]: Client\rError\r\t\n",
			},
		},
		logTests: []logTest{
			{
				name:       "Fail",
				action:     Fail,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[FAIL]: This is a test\n",
			},
			{
				name:       "Warn",
				action:     Warn,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[WARN]: This is a test\n",
			},
			{
				name:       "Debug",
				action:     Debug,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[DEBUG]: This is a test\n",
			},
			{
				name:       "Info",
				action:     Info,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "[INFO]: This is a test\n",
			},
		},
	},
	{
		name:  "Miscellaneous",
		level: Information,
		errorTests: []errorTest{
			{
				name:       "NilError",
				err:        nil,
				wantString: "",
			},
		},
	},
	{
		name:  "InvalidLevel",
		level: -1,
		errorTests: []errorTest{
			{
				name:       "ServerError",
				err:        errors.NewServer("ServerError"),
				wantString: "",
			},
			{
				name:       "ClientError",
				err:        errors.NewClient("ClientError"),
				wantString: "",
			},
		},
		logTests: []logTest{
			{
				name:       "Fail",
				action:     Fail,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Warn",
				action:     Warn,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Debug",
				action:     Debug,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
			{
				name:       "Info",
				action:     Info,
				args:       []interface{}{"This", "is", "a", "test"},
				wantString: "",
			},
		},
	},
}

func TestLog(t *testing.T) {
	for _, test := range logTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			SetLevel(test.level)
			defer SetLevel(Silent)
			var buf strings.Builder
			writer = &buf

			// Execute Error tests
			for _, errorTest := range test.errorTests {
				t.Run(errorTest.name, func(t *testing.T) {
					// Setup
					buf.Reset()

					// Execute
					Error(errorTest.err)

					// Verify
					gotString := buf.String()
					if gotString != errorTest.wantString {
						t.Errorf("Got string `%s`; want `%s`", gotString, errorTest.wantString)
					}
				})
			}

			// Execute log level tests
			for _, logTest := range test.logTests {
				t.Run(logTest.name, func(t *testing.T) {
					// Setup
					buf.Reset()

					// Execute
					logTest.action(logTest.args...)

					// Verify
					gotString := buf.String()
					if gotString != logTest.wantString {
						t.Errorf("Got string `%s`; want `%s`", gotString, logTest.wantString)
					}
				})
			}
		})
	}
}
