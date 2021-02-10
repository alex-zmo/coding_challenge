package server

import (
	"database/sql"
	"errors"
	"github.com/gmo-personal/coding_challenge/database"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server/auth"
	"github.com/gmo-personal/coding_challenge/server/csrf"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
)

// Serves the metrics and account plan info in json format to the client.
func getMetricHandler(w http.ResponseWriter, r *http.Request) {
	// Verifies that the CSRF cookie and header match. If they do not match, the request is unauthorized.
	err := csrf.VerifyCSRF(r)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeUnauthorized(w)
		return
	}

	// Verifies that the auth token exists as a valid session token and returns the accountID,
	// otherwise the request is unauthorized.
	accountID, err := auth.ValidateToken(r)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeUnauthorized(w)
		return
	}

	// Retrieves db from context.
	db, ok := r.Context().Value("db").(*sql.DB)
	if !ok {
		utils.Logger.Println(errors.New("db unset"))
		utils.ServeInternalServerError(w)
		return
	}

	// Gets the dashboard info associated with an account, otherwise serves not found.
	dashboardJson, err := getDashboardInfo(db, accountID)
	if err != nil {
		utils.Logger.Println(err)
		utils.ServeInternalServerError(w)
		return
	}

	// Attempts to serve the dashboard json.
	utils.ServeJson(w, dashboardJson)
}

// Gets the metrics count and account plan associated with an account.
func getDashboardInfo(db *sql.DB, accountID string) (*model.DashboardInfo, error) {
	// Gets the account associated with an account ID.
	existingAccount, err := database.SelectAccount(db, accountID)
	if err != nil {
		return nil, err
	}
	// Gets the metrics count associated with an account ID.
	metricCount, err := database.CountMetrics(db, accountID)
	if err != nil {
		return nil, err
	}
	dashboardJson := &model.DashboardInfo{
		MetricsCount: metricCount,
		Plan:         existingAccount.Plan,
	}
	return dashboardJson, nil
}
