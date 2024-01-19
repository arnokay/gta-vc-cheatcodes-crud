package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(http.HandlerFunc(app.notFoundResponse))
	r.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)
		r.Route("/cheatcodes", func(r chi.Router) {
			r.Get("/", app.listCheatcodesHandler)
			r.Post("/", app.createCheatcodeHandler)
			r.Get("/{id}", app.showCheatcodeHandler)
			r.Patch("/{id}", app.updateCheatcodeHandler)
			r.Delete("/{id}", app.deleteCheatcodeHandler)
		})
	})
	return r
}
