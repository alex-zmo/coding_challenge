package auth

import (
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
	"time"
)


// Sends response to nullify the cookies and deletes the tokens from sessions.
func NullifyTokens(w http.ResponseWriter, r *http.Request) {
	// Attempt to extract the token from the cookie, returns and error if fails.
	authToken, err := utils.ExtractTokenFromCookie( authCookieName, r)
	if err != nil {
			utils.LogError(err)
			utils.ServeInternalServerError(w)
			return
	}
	delete(Sessions, authToken)

	// Sets auth cookie expire time to -1000 hours from now time with no value.
	authCookie := http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(w, &authCookie)

	// Sets CSRF cookie expire time to -1000 hours from now time with no value.
	CSRFCookie := http.Cookie{
		Name:     csrfCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	http.SetCookie(w, &CSRFCookie)
}
