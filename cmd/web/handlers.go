package main

import (
	"net/http"

	"github.com/wumbabum/home_assist/internal/response"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) restricted(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/restricted.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}
