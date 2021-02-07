package server

import (
	"github.com/gmo-personal/coding_challenge/server/utils"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("./server/templates/index.tmpl", "./server/templates/dashboard.tmpl"))

// Renders a given template to response writer.
func RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	// Renders page from template, if failed, returns internal server error.
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", p)
	if err != nil {
		utils.LogError(err)
		utils.ServeInternalServerError(w)
	}
}
