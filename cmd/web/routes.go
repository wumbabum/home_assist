package main

import (
	"net/http"

	"github.com/wumbabum/home_assist/assets"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(app.notFound)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)
	mux.Use(app.sessionManager.LoadAndSave)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("/static/*", fileServer)

	mux.Get("/", app.home)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireBasicAuthentication)

		mux.Get("/restricted-basic-auth", app.restricted)
	})

	return mux
}
