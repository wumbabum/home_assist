package main

import (
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	// Create oidc request and create session state

	// Stub: use fixed state (INSECURE - for testing only)
	state := "stub-state-replace-later"

	http.Redirect(w, r, app.auth0.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
	// TODO: Verify state parameter matches session
	// TODO: Exchange authorization code for tokens
	// TODO: Verify ID token
	// TODO: Extract user profile from claims
	// TODO: Store profile and tokens in session

	// Stub profile
	profile := UserProfile{
		Sub:   "auth0-stubbed-user-id",
		Email: "stub@example.com",
		Name:  "Stub User",
	}

	app.sessionManager.Put(r.Context(), "profile", profile)
	// Stub: redirect to stubbed profile
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Destroy session
	// TODO: Redirect to Auth0 logout endpoint with returnTo URL

	// Stub: redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
