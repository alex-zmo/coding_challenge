package auth

import (
	"errors"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
	"time"
)

// Sends response to nullify the cookies and deletes the tokens from sessions.
func DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Attempt to extract the token from the cookie, returns and error if fails.
	authTokenCookie, err := r.Cookie(authCookieName)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeInternalServerError(w)
	}
	// Retrieves sessions from context.
	sess, ok := r.Context().Value("sess").(map[string]Session)
	if !ok {
		utils.Logger.Println(errors.New("sess unset"))
		utils.ServeInternalServerError(w)
		return
	}
	authToken := authTokenCookie.Value
	delete(sess, authToken)

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
	csrfCookie := http.Cookie{
		Name:     csrfCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	}

	http.SetCookie(w, &csrfCookie)
}
