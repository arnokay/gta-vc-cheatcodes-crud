package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(app.recoverPanic)

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
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.registerUserHandler)
		})
	})
	return r
}
