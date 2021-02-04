package server

import (
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server/auth"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"net/http"
)

// TODO CHOOSE BEARER TOKEN FOR METRIC to be hard coded once metrics functional.
// TODO Implement SHORT POLLING and then implement WEBSOCKET real time event system.
// TODO NOTE: Integration testing should be done in the future for handlers, outside POC scope.
//				Unit tests could typically be done with httptest library.

// TODO NOTE: Full implementation would have Get, Post, and Delete methods as well, POC only needs patch for plan.
func accountHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableContentSecurityPolicy(&w)
	switch r.Method {
	case http.MethodPatch:
		patchAccountHandler(w, r)
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableContentSecurityPolicy(&w)
	switch r.Method {
	case http.MethodPost:
		auth.PostTokenHandler(w, r)
	case http.MethodGet:
		auth.ValidateTokenHandler(w, r)
	case http.MethodDelete:
		auth.NullifyTokens(w, r)
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableContentSecurityPolicy(&w)
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, "index", &model.IndexPage{})
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableContentSecurityPolicy(&w)
	switch r.Method {
	case http.MethodGet:
		CSRFToken, err := utils.ExtractTokenFromCookie( "csrf-cookie", r)
		if err != nil {
				utils.LogError(err)
		}
		RenderTemplate(w, "dashboard", &model.DashboardPage{CSRFToken: CSRFToken})
	default:
		utils.ServeMethodNotAllowed(w)
	}
}

func InitServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/account/", accountHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/token/", tokenHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("server/static"))))
	err := http.ListenAndServeTLS(":443", "certs/server-cert.pem", "certs/server-key.pem", nil)
	if err != nil {
			panic(err)
	}
}
