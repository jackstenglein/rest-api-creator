// Package auth provides utilities for authentication and authorization.
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// TODO: dynamically pull this key from AWS KMS
const key = "TODO:changeThisKey"

// UserGetter wraps the GetUser function for database dependency injection when verifying cookies.
type UserGetter interface {
	GetUserInfo(string) (*dao.User, error)
}

// splitCookie takes a cookie in the following format
//		email#token#mac
// and separates it into its component parts. If the cookie is not in the correct format,
// splitCookie returns emtpy strings and an error.
func splitCookie(cookie string) (email string, token string, mac string, err error) {
	slice := strings.Split(cookie, "#")
	if len(slice) != 3 || len(slice[0]) == 0 || len(slice[1]) == 0 || len(slice[2]) == 0 {
		return "", "", "", errors.NewClient("Incorrect cookie format")
	}
	return slice[0], slice[1], slice[2], nil
}

// GenerateToken returns a session token created by hex encoding a random array of 8 bytes.
// The array is created using the crypto/rand package. If an error occurs, GenerateToken
// returns the empty string along with the error.
func GenerateToken() (token string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateCookie returns a cookie created using the following format:
// 		email#token#mac
// where mac is the SHA256 hash of email#token. If an error occurs, GenerateCookie returns
// the empty string along with the error.
func GenerateCookie(email string, token string) (cookie string, err error) {
	macString := email + "#" + token
	macBytes, err := computeMAC([]byte(macString))
	if err != nil {
		return "", err
	}
	return macString + "#" + hex.EncodeToString(macBytes), nil
}

// ExtractCookie extracts the actual session cookie value from a `Cookie` header value. The header
// must be in the following format:
//		session=<cookie value>
// If the header value is not in the correct format, the empty string will be returned.
func ExtractCookie(cookieHeader string) string {
	startIndex := strings.Index(cookieHeader, "session=")
	if startIndex != 0 {
		return ""
	}
	return cookieHeader[len("session="):]
}

// VerifyCookie checks that cookie is in the correct format, its mac is correct, and its
// contained (email, token) tuple matches the (email, token) tuple stored in db. VerifyCookie
// returns the email contained in the cookie if the cookie is valid. If the cookie is invalid,
// an error is returned and email is the empty string.
//
// db must implement the GetUser(string) (*dao.User, error) method
func VerifyCookie(cookie string, db UserGetter) (email string, err error) {
	email, token, mac, err := splitCookie(cookie)
	if err != nil {
		return "", errors.NewClient("Not authenticated")
	}

	expectedMac, err := computeMAC([]byte(email + "#" + token))
	if err != nil {
		return "", errors.Wrap(err, "Failed to compute verification MAC")
	}

	messageMac, err := hex.DecodeString(mac)
	if err != nil {
		return "", errors.Wrap(err, "Failed to hex decode message mac")
	}
	if !hmac.Equal(expectedMac, messageMac) {
		return "", errors.NewClient("Not authenticated")
	}

	user, err := db.GetUserInfo(email)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get user from database")
	}
	if user.Token != token {
		return "", errors.NewClient("Not authenticated")
	}
	return email, nil
}

// computeMAC returns the sha256 hmac of the given byte slice.
func computeMAC(b []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(key))
	_, err := mac.Write(b)
	return mac.Sum(nil), err
}
