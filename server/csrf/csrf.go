package csrf

import (
	"crypto/subtle"
	"encoding/hex"
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
	token, err := utils.ExtractTokenFromCookie( CookieName, r)
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
	headerToken := utils.GetHeader(r, HeaderName)
	if len(headerToken) == 0 {
			return errors.New("no token found in header")
	}
	// Reads the CSRF token from the cookie, otherwise returns an error.
	cookieToken, err := utils.ExtractTokenFromCookie( CookieName, r)
	if err != nil {
			return err
	}

	// Decodes the two tokens, returns error if token decode fails.
	decodedHeader, err := decode(headerToken)
	if err != nil  {
			return err
	}
	decodedCookie, err := decode(cookieToken)
	if err != nil  {
			return err
	}

	// Compares the two tokens to each other.
	if subtle.ConstantTimeCompare(decodedHeader, decodedCookie) != 1 {
			return errors.New("tokens do not match")
	}
	return nil
}

// Decodes the token string.
func decode(token string) ([]byte, error) {
	decodedToken, err := hex.DecodeString(token)
	if err != nil {
			return nil, err
	}

	// If the decoded token has a different length than required, return an error.
	if len(decodedToken) != tokenBytesLen {
			return nil, errors.New("unexpected token length")
	}

	return decodedToken, nil
}


