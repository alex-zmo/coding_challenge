package auth

import (
	"encoding/json"
	"github.com/gmo-personal/coding_challenge/server/csrf"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"time"

	"github.com/gmo-personal/coding_challenge/database"
	"github.com/gmo-personal/coding_challenge/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Structure of a session, containing account ID and expiry time.
type Session struct {
	AccountID string
	Expires   time.Time
}

// In memory sessions table to keep track of active sessions.
// TODO This should be a database table to avoid memory and concurrency issues as well as scaling.
var Sessions = make(map[string]Session)

// Generates auth and CSRF token and returns them to the client.
func PostTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Attempts to read the body, if unable to do so, return bad request.
	accountJson, err := utils.GetBody(r)
	if err != nil {
			utils.LogError(err)
			utils.ServeBadRequest(w)
			return
	}

	// Attempts to parse json into account model, if unable to do so, return bad request.
	account := &model.Account{}
	err = json.Unmarshal(accountJson, account)
	if err != nil {
			utils.LogError(err)
			utils.ServeBadRequest(w)
			return
	}

	// Attempts to find an existing account by username and compares input and existing passwords,
	// Serves unauthorized if they do not match.
	existingAccount, err := database.SelectAccountByUsername(account.Username)
	if err != nil {
			utils.LogError(w)
			utils.ServeNotFound(w)
			return
	}
	if existingAccount != nil {
		// Uses bcrypt library to compare the passwords, returns unauthorized if comparison fails.
		err = bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(account.Password))
		if err != nil {
				utils.LogError(w)
				utils.ServeUnauthorized(w)
				return
		}
		// Attempts to create an auth token for the account, serves internal server error on fail.
		authToken, err := createAuth(existingAccount.ID)
		if err != nil {
				utils.LogError(w)
				utils.ServeInternalServerError(w)
				return
		}
		// Sends response to save the cookie.
		utils.SaveCookie(authCookieName, authToken, w)

		w.Header().Set("Content-Type", "application/json")

		// Attempts to create an CSRF token for the account, serves internal server error on fail.
		CSRFToken, err := csrf.CreateCSRF(w, r)
		if err != nil {
				utils.LogError(w)
				utils.ServeInternalServerError(w)
				return
		}
		utils.ServeJson(w, CSRFToken)
		return
	}
	utils.ServeUnauthorized(w)
}

// Creates an auth token.
func createAuth(accountID string) (string, error) {
	token, err := utils.CryptoRandomString(32)
	if err != nil {
			return "", err
	}
	// Adds the token to Sessions with a 60 minute expiry.
	Sessions[token] = Session{accountID, time.Now().Add(time.Minute * 60)}
	return token, nil
}


