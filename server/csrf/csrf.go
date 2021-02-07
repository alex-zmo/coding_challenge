package csrf

import (
	"crypto/subtle"
	"errors"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
)

const (
	CookieName    = "csrf-cookie"
	HeaderName    = "X-CSRF-Token"
	tokenBytesLen = 32
)

// Creates a new CSRF token if one does not exist.
func CreateCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
	token, err := utils.ExtractTokenFromCookie(CookieName, r)
	// If there was an error retrieving the token, the token doesn't exist, so create a new one.
	if err != nil || len(token) == 0 {
		token, err = utils.CryptoRandomString(tokenBytesLen)
		if err != nil {
			return "", err
		}
	}
	// Sends response to save the cookie.
	utils.SaveCookie(CookieName, token, w)
	return token, nil
}

// Verifies the CSRF header and the CSRF Cookie against each other.
func VerifyCSRF(r *http.Request) error {
	// Reads the CSRF token from the X-CSRF-Token header, otherwise returns an error.
	headerToken := r.Header.Get(HeaderName)
	if len(headerToken) == 0 {
		return errors.New("no token found in header")
	}
	// Reads the CSRF token from the cookie, otherwise returns an error.
	csrfCookie, err := r.Cookie(CookieName)
	if err != nil {
		return err
	}
	cookieToken := csrfCookie.Value
	// Compares the two tokens to each other.
	if subtle.ConstantTimeCompare([]byte(headerToken), []byte(cookieToken)) != 1 {
		return errors.New("tokens do not match")
	}
	return nil
}
