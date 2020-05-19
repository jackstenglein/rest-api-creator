package auth

import (
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
	"testing"
)

type userGetterFunc func(string) (*dao.User, error)

func (f userGetterFunc) GetUserInfo(email string) (*dao.User, error) {
	return f(email)
}

func getUserMock(email string, token string, tokenMatches bool, mockErr error) userGetterFunc {
	return func(input string) (*dao.User, error) {
		if input == email {
			if !tokenMatches {
				token = ""
			}
			return &dao.User{Email: email, Token: token}, mockErr
		}
		return nil, errors.NewServer("Incorrect input to getUserMock")
	}
}

var sessionTests = []struct {
	name         string
	email        string
	tokenMatches bool
	mockErr      error
	wantErr      error
	wantEmail    string
}{
	{
		name:         "DatabaseError",
		email:        "test@example.com",
		tokenMatches: true,
		mockErr:      errors.NewServer("Database failure"),
		wantErr:      errors.Wrap(errors.NewServer("Database failure"), "Failed to get user from database"),
	},
	{
		name:         "IncorrectToken",
		email:        "test@example.com",
		tokenMatches: false,
		wantErr:      errors.NewClient("Not authenticated"),
	},
	{
		name:         "SuccessfulInvocation",
		email:        "test@example.com",
		tokenMatches: true,
		wantEmail:    "test@example.com",
	},
}

func TestAllMethods(t *testing.T) {
	for _, test := range sessionTests {
		t.Run(test.name, func(t *testing.T) {
			token, err := GenerateToken()
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}

			cookie, err := GenerateCookie(test.email, token)
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}

			db := getUserMock(test.email, token, test.tokenMatches, test.mockErr)
			email, err := VerifyCookie(cookie, db)
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
			if email != test.wantEmail {
				t.Errorf("Got email %s; want email %s", email, test.wantEmail)
			}
		})
	}
}

var extractCookieTests = []struct {
	name       string
	header     string
	wantCookie string
}{
	{
		name:       "MissingSession",
		header:     "cookievalue;HttpOnly;SameSite=strict;Secure",
		wantCookie: "",
	},
	{
		name:       "SessionNotFirst",
		header:     "Asdfsession=cookievalue",
		wantCookie: "",
	},
	{
		name:       "NoSessionValue",
		header:     "session=",
		wantCookie: "",
	},
	{
		name:       "CorrectHeader",
		header:     "session=cookievalue",
		wantCookie: "cookievalue",
	},
}

func TestExtractCookie(t *testing.T) {
	for _, test := range extractCookieTests {
		t.Run(test.name, func(t *testing.T) {
			cookie := ExtractCookie(test.header)
			if cookie != test.wantCookie {
				t.Errorf("Got cookie '%s'; want '%s'", cookie, test.wantCookie)
			}
		})
	}
}

func TestInvalidCookie(t *testing.T) {
	t.Run("IncorrectFormat", func(t *testing.T) {
		_, _, _, err := splitCookie("")
		if err == nil {
			t.Errorf("Empty string cookie considered valid")
		}

		_, _, _, err = splitCookie("#token#mac")
		if err == nil {
			t.Errorf("Cookie `#token#mac` considered valid")
		}

		_, _, _, err = splitCookie("token#mac#")
		if err == nil {
			t.Errorf("Cookie `token#mac#` considered valid")
		}

		_, _, _, err = splitCookie("email#token#mac#")
		if err == nil {
			t.Errorf("Cookie `email#token#mac#` considered valid")
		}

		_, _, _, err = splitCookie("email#token##mac")
		if err == nil {
			t.Errorf("Cookie `email#token##mac` considered valid")
		}

		email, err := VerifyCookie("email#token##mac", nil)
		if email != "" {
			t.Errorf("Got email %s; want empty string", email)
		}
		if errors.Message(err) != "Not authenticated" {
			t.Errorf("Got error %v", err)
		}
	})

	t.Run("IncorrectEncoding", func(t *testing.T) {
		email, err := VerifyCookie("email#token#mac", nil)
		if email != "" {
			t.Errorf("Got email %s; want empty string", email)
		}
		if errors.Message(err) != "Failed to hex decode message mac" {
			t.Errorf("Got error %v", err)
		}

		email, err = VerifyCookie("email#token#0123456789abcdef", nil)
		if email != "" {
			t.Errorf("Got email %s; want empty string", email)
		}
		if errors.Message(err) != "Not authenticated" {
			t.Errorf("Got error %v", err)
		}
	})
}
