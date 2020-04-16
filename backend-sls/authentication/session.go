package authentication

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

const key = "TODO:changeThisKey"

type Authenticator interface {
	GenerateToken() (string, error)
	GenerateCookie(string, string) (string, error)
	VerifyCookie(string, string, string) (bool, error)
}

type SessionAuthenticator struct {}

func NewSessionAuthenticator() *SessionAuthenticator {
	return &SessionAuthenticator{}
}

// GenerateToken returns a string created by hex encoding a random array of 8 bytes.
// The array is created using the crypto/rand package. If an error occurs, GenerateToken
// returns the empty string along with the error.
func (auth *SessionAuthenticator) GenerateToken() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	sessionId := hex.EncodeToString(b)
	return sessionId, nil
}

// GenerateCookie returns a string created using the following format:
// 		email#token#hmac(email#token)
// If an error occurs, GenerateCookie returns the empty string along with the error.
func (auth *SessionAuthenticator) GenerateCookie(email string, token string) (string, error) {
	hmacString := email + "#" + token
	hmacBytes, err := computeMAC([]byte(hmacString))
	if err != nil {
		return "", err
	}
	return hmacString + "#" + hex.EncodeToString(hmacBytes), nil
}

// VerifyCookie returns true if the hmac matches the given email and token and false otherwise.
func (auth *SessionAuthenticator) VerifyCookie(email string, token string, messageMac string) (bool, error) {
	hmacString := email + "#" + token
	expectedHmac, err := computeMAC([]byte(hmacString))
	if err != nil {
		return false, err
	}
	messageMacBytes, err := hex.DecodeString(messageMac)
	if err != nil {
		return false, err
	}
	return hmac.Equal(expectedHmac, messageMacBytes), nil
}

// computeMAC returns the sha256 hmac of the given byte slice.
func computeMAC(b []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(key))
	_, err := mac.Write(b)
	return mac.Sum(nil), err
}
