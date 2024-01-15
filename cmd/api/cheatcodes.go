package main

import (
	"net/http"

	"github.com/arnokay/gta-vc-cheatcodes-crud/internal/data"
)

func (app *application) createCheatcodeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Code        string   `json:"code"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
    app.badRequestResponse(w, r, err)
    return
	}

	app.writeJSON(w, 200, envelope{"cheatcode": input}, nil)
}

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
