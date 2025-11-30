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

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	profile := app.sessionManager.Get(r.Context(), "profile")

	data := app.newTemplateData(r)
	data["Profile"] = profile

	err := response.Page(w, http.StatusOK, data, "pages/user.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}
