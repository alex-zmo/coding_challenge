package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gmo-personal/coding_challenge/database"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"github.com/microcosm-cc/bluemonday"
	"io/ioutil"
	"net/http"
	"strings"
)

// The metrics token is a hard coded randomly generated opaque token (utils.CryptoRandomString)
const (
	HardCodedMetricsToken = "0b2a4a30602b7f6d1079593cabfdb9386eff2041b863a9c678915158eb60fdac"
)

// Inserts the new metric into the database if applicable to plan and count.
func postMetricHandler(w http.ResponseWriter, r *http.Request) {
	// Sets bluemonday sanitize whitelist policy to strict, ensures anything into DB does not contain html
	// restricted inputs.
	policy := bluemonday.StrictPolicy()
	// Attempts to read the body, if unable to do so, return bad request.
	metricsJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeBadRequest(w)
		return
	}
	// Attempts to parse json into metrics model, if unable to do so, return bad request.
	metric := &model.Metric{}
	err = json.Unmarshal(metricsJson, &metric)
	if err != nil || len(metric.UserID) == 0 || len(metric.AccountID) == 0 {
		utils.Logger.Println(err)
		utils.ServeBadRequest(w)
		return
	}

	// TODO in future change from hard-coded token to more secure randomly generated opaque token
	// Validates the metrics authorization bearer token, otherwise serves forbidden.
	auth := strings.Split(r.Header["Authorization"][0], "Bearer")
	if len(auth) < 2 || strings.TrimSpace(auth[1]) != HardCodedMetricsToken {
		w.Header().Set("Content-Type", "application/json")
		utils.Logger.Println(errors.New("metrics bearer invalid"))
		utils.ServeForbidden(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Sanitizes metrics post input to prevent undesired inputs into database.
	metric.AccountID = policy.Sanitize(metric.AccountID)
	metric.UserID = policy.Sanitize(metric.UserID)
	metric.Timestamp = policy.Sanitize(metric.Timestamp)
	metric.Timestamp = metric.Timestamp[:len(metric.Timestamp)-1]

	// Retrieves db from context.
	db, ok := r.Context().Value("db").(*sql.DB)
	if !ok {
		utils.Logger.Println(errors.New("db unset"))
		utils.ServeInternalServerError(w)
		return
	}

	// Begins database transaction, serves error if failed.
	tx, err := database.StartTransaction(db)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeInternalServerError(w)
		return
	}

	// Defers the transaction commit / rollback till end of function
	defer database.ResolveTransaction(tx)
	// Checks to see if there is an existing account of the same account ID in the database,
	// if not, return not found.
	existingAccount, err := database.SelectAccount(tx, db, metric.AccountID)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeUnauthorized(w)
		return
	}

	// Returns count of metrics associated with account ID,serves internal server error if failed.
	metricsCount, err := database.CountMetrics(tx, db, metric.AccountID)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeInternalServerError(w)
		return
	}

	// Checks if able to insert metric based on plan and count, otherwise serve forbidden.
	if metricsCount < 100 || existingAccount.Plan == 1 && metricsCount < 1000 {
		// Inserts new metric, serves internal server error if failed.
		err = database.InsertMetric(tx, db, metric)
		if err != nil {
			utils.Logger.Println(err)
			utils.ServeInternalServerError(w)
			return
		}
		utils.ServeCreated(w)
	} else {
		utils.ServeForbidden(w)
	}
}
