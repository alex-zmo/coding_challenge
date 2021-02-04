package auth

import (
	"errors"
	"github.com/gmo-personal/coding_challenge/server/csrf"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
	"time"
)

const (
	authCookieName = "auth-cookie"
	csrfCookieName = "csrf-cookie"
)

// Validates that the auth token exists as a valid session token
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Verifies that the CSRF cookie and header match. If they do not match, the request is unauthorized.
	err := csrf.VerifyCSRF(r)
	if err != nil {
			utils.LogError(err)
			utils.ServeUnauthorized(w)
			return
	}

	// Verifies that the auth token exists as a valid session token,
	// otherwise the request is unauthorized.
	_, err = ValidateToken(r)
	if err != nil {
			utils.LogError(err)
			utils.ServeUnauthorized(w)
			return
	}
}

// Validates the auth token and returns the account ID associated with the token.
func ValidateToken(r *http.Request) (string, error) {
	// Attempt to extract the token from the cookie, returns and error if fails.
	tokenStr, err := utils.ExtractTokenFromCookie( authCookieName, r)
	if err != nil {
			return "", err
	}
	// Checks if the token is valid in the session and unexpired and returns account ID, otherwise returns error.
	if session, seen := Sessions[tokenStr]; seen {
		if session.Expires.Unix() <= time.Now().Unix() {
				delete(Sessions, tokenStr)
				return "", errors.New("unauthorized")
		}
		accountID := session.AccountID
		return accountID, nil
	}
	return "", errors.New("unauthorized")
}
