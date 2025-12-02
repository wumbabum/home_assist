package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	// Create oidc request and create session state

	// Stub: use fixed state (INSECURE - for testing only)
	csrfToken, err := csrfToken()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "oauth_state", csrfToken)

	http.Redirect(w, r, app.auth0.AuthCodeURL(csrfToken), http.StatusTemporaryRedirect)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
	savedState := app.sessionManager.GetString(r.Context(), "oauth_state")
	if r.URL.Query().Get("state") != savedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	token, err := app.auth0.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	idToken, err := app.auth0.VerifyIDToken(r.Context(), token)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var profile UserProfile

	if err := idToken.Claims(&profile); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("auth0 profile data", "profile", profile)

	app.sessionManager.Put(r.Context(), "access_token", token.AccessToken)
	app.sessionManager.Put(r.Context(), "profile", profile)
	// Stub: redirect to stubbed profile
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.Destroy(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	logoutURL := "https://" + app.config.auth0.domain + "/v2/logout?returnTo=" +
		app.config.baseURL + "&client_id=" + app.config.auth0.clientID

	http.Redirect(w, r, logoutURL, http.StatusSeeOther)
}

func csrfToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
