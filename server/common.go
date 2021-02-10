package server

import (
	"context"
	"database/sql"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server/auth"
	"github.com/gmo-personal/coding_challenge/server/csrf"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
	"os"
)

// TODO Implement SHORT POLLING and then implement WEBSOCKET real time event system.
// TODO NOTE: Integration testing should be done in the future for handlers, outside POC scope.
//				Unit tests could typically be done with httptest library.

const (
	csrfCookieName = "csrf-cookie"
)

// TODO NOTE: Full implementation would have Get, Post, and Delete methods as well, POC only needs patch for plan.
// Handles account related requests.
func accountHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPatch:
		patchAccountHandler(w, r)
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

// Handles token related requests.
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.PostTokenHandler(w, r)
	case http.MethodDelete:
		auth.DeleteTokenHandler(w, r)
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postMetricHandler(w, r)
	case http.MethodGet:
		getMetricHandler(w, r)
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

// Renders login page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		csrfToken, err := csrf.CreateCSRF(w, r)
		if err != nil {
			utils.Logger.Println(err)
			utils.ServeInternalServerError(w)
			return
		}
		RenderTemplate(w, "index", &model.IndexPage{CSRFToken: csrfToken})
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

// Renders dashboard page.
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, err := auth.ValidateToken(r)
		if err != nil {
			utils.Logger.Println(err)
			utils.ServeUnauthorized(w)
			return
		}
		csrfCookie, err := r.Cookie(csrfCookieName)
		if err != nil {
			utils.Logger.Println(err)
			utils.ServeInternalServerError(w)
			return
		}
		RenderTemplate(w, "dashboard", &model.DashboardPage{CSRFToken: csrfCookie.Value})
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

// Wrapper that adds content security policy to handler.
func csp(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.EnableContentSecurityPolicy(w)
		fn(w, r)
	}
}

// Wrapper that injects DB to handler.
func injectDB(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// TODO This should be a database table to avoid memory and concurrency issues as well as scaling.
// Wrapper that injects Sessions to handler.
func injectSessions(sess map[string]auth.Session, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "sess", sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Initializes server.
func InitServer(db *sql.DB) {
	sess := make(map[string]auth.Session)
	http.HandleFunc("/", csp(indexHandler))
	http.HandleFunc("/account/", injectSessions(sess, injectDB(db, csp(accountHandler))))
	http.HandleFunc("/dashboard", injectSessions(sess, csp(dashboardHandler)))
	http.HandleFunc("/token/", injectSessions(sess, injectDB(db, csp(tokenHandler))))
	http.HandleFunc("/metrics", injectSessions(sess, injectDB(db, csp(metricsHandler))))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("server/static"))))
	err := http.ListenAndServeTLS(":"+os.Getenv("SRV_PORT"), os.Getenv("SRV_CERT_PATH"), os.Getenv("SRV_KEY_PATH"), nil)
	if err != nil {
		utils.Logger.Println(err)
		return
	}
}
