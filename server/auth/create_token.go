package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gmo-personal/coding_challenge/server/csrf"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"io/ioutil"
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

// Generates auth token and return to the client as a cookie.
func PostTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Verifies that the CSRF cookie and header match. If they do not match, the request is unauthorized.
	err := csrf.VerifyCSRF(r)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeUnauthorized(w)
		return
	}

	// Attempts to read the body, if unable to do so, return bad request.
	accountJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeBadRequest(w)
		return
	}
	defer r.Body.Close()

	// Attempts to parse json into account model, if unable to do so, return bad request.
	account := &model.Account{}
	err = json.Unmarshal(accountJson, account)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeBadRequest(w)
		return
	}
	// Retrieves db from context.
	db, ok := r.Context().Value("db").(*sql.DB)
	if !ok {
		utils.Logger.Println(errors.New("db unset"))
		utils.ServeInternalServerError(w)
		return
	}
	// Attempts to find an existing account by username and compares input and existing passwords,
	// Serves unauthorized if they do not match.
	existingAccount, err := database.SelectAccountByUsername(db, account.Username)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeUnauthorized(w)
	}
	// Uses bcrypt library to compare the passwords, returns unauthorized if comparison fails.
	err = bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(account.Password))
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeUnauthorized(w)
		return
	}
	// Retrieves sessions from context.
	sess, ok := r.Context().Value("sess").(map[string]Session)
	if !ok {
		utils.Logger.Println(errors.New("sess unset"))
		utils.ServeInternalServerError(w)
		return
	}
	// Attempts to create an auth token for the account, serves internal server error on fail.
	authToken, err := createAuth(existingAccount.ID, sess)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeInternalServerError(w)
		return
	}
	// Sends response to save the cookie.
	utils.SaveCookie(authCookieName, authToken, w)
}

// Creates an auth token.
func createAuth(accountID string, sess map[string]Session) (string, error) {
	token, err := utils.CryptoRandomString(32)
	if err != nil {
		return "", err
	}
	// Adds the token to Sessions with a 60 minute expiry.
	sess[token] = Session{accountID, time.Now().Add(time.Minute * 60)}
	return token, nil
}
