package server

import (
	"encoding/json"
	"github.com/gmo-personal/coding_challenge/database"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server/auth"
	"github.com/gmo-personal/coding_challenge/server/csrf"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
)

// Updates the existing account plan to the new account plan if applicable.
func patchAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Verifies that the CSRF cookie and header match. If they do not match, the request is unauthorized.
	err := csrf.VerifyCSRF(r)
	if err != nil {
			utils.LogError(err)
			utils.ServeUnauthorized(w)
			return
	}

	// Verifies that the auth token exists as a valid session token and returns the accountID,
	// otherwise the request is unauthorized.
	accountID, err := auth.ValidateToken(r)
	if err != nil {
			utils.LogError(err)
			utils.ServeUnauthorized(w)
			return
	}

	// Attempts to read the body, if unable to do so, return bad request.
	accountJson, err := utils.GetBody(r)
	if err != nil {
			utils.LogError(err)
			utils.ServeBadRequest(w)
			return
	}

	// Attempts to parse json into account model, if unable to do so, return bad request.
	newAccount := &model.Account{}
	err = json.Unmarshal(accountJson, newAccount)
	if err != nil {
			utils.LogError(err)
			utils.ServeBadRequest(w)
			return
	}

	// Checks to see if there is an existing account of the same account ID in the database,
	// if not, return not found.
	existingAccount, err := database.SelectAccount(accountID)
	if err != nil {
			utils.LogError(err)
			utils.ServeNotFound(w)
			return
	}

	// Updates the existing account's plan to the new plan.
	if newAccount.Plan != existingAccount.Plan {
		existingAccount.Plan = newAccount.Plan
	}

	// Updates the account plan in the database.
	err = database.UpdateAccount(existingAccount)
	if err != nil {
			utils.LogError(err)
			utils.ServeInternalServerError(w)
			return
	}
}
