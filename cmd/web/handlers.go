package main

import (
	"net/http"

	"github.com/wumbabum/home_assist/internal/response"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	profileData := app.sessionManager.Get(r.Context(), "profile")
	profile, _ := profileData.(UserProfile)

	app.logger.Info("profile data", "profile", profile)

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	profileData := app.sessionManager.Get(r.Context(), "profile")

	profile, ok := profileData.(UserProfile)

	app.logger.Info("profile data", "profile", profile)

	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data["Profile"] = profile

	err := response.Page(w, http.StatusOK, data, "pages/user.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}
