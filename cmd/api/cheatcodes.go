package main

import (
	"net/http"

	"github.com/arnokay/gta-vc-cheatcodes-crud/internal/data"
)

// func (app *application) createCheatcodeHandler(w http.ResponseWriter, r *http.Request) {
//
// }

func (app *application) showCheatcodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
    return
	}

	cheatcode := data.Movie{
		ID: id,
	}

  err = app.writeJSON(w, http.StatusOK, envelope{"cheatcode": cheatcode}, nil)
	if err != nil {
	  app.serverErrorResponse(w, r, err)
  }
}
