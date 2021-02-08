package auth

import (
	"errors"
	"net/http"
	"time"
)

const (
	authCookieName = "auth-cookie"
	csrfCookieName = "csrf-cookie"
)

// Validates the auth token and returns the account ID associated with the token.
func ValidateToken(r *http.Request) (string, error) {
	// Attempt to extract the token from the cookie, returns and error if fails.
	authCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return "", err
	}
	tokenStr := authCookie.Value
	// Retrieves sessions from context.
	sess, ok := r.Context().Value("sess").(map[string]Session)
	if !ok {
		return "", errors.New("unset sessions")
	}
	// Checks if the token is valid in the session and unexpired and returns account ID, otherwise returns error.
	if session, seen := sess[tokenStr]; seen {
		if session.Expires.Before(time.Now()) {
			delete(sess, tokenStr)
			return "", errors.New("unauthorized token")
		}
		accountID := session.AccountID
		return accountID, nil
	}
	return "", errors.New("unauthorized token")
}
