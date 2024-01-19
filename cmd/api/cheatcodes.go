package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnokay/gta-vc-cheatcodes-crud/internal/data"
	"github.com/arnokay/gta-vc-cheatcodes-crud/internal/validator"
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

	cheatcode := &data.Cheatcode{
		Code:        input.Code,
		Description: input.Description,
		Tags:        input.Tags,
	}

	v := validator.New()

	if data.ValidateCheatcode(v, cheatcode); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cheatcodes.Insert(cheatcode)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/cheatcodes/%d", cheatcode.ID))

	err = app.writeJSON(w, 200, envelope{"cheatcode": input}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCheatcodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	cheatcode, err := app.models.Cheatcodes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cheatcode": cheatcode}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCheatcodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	cheatcode, err := app.models.Cheatcodes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(cheatcode.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Code        *string  `json:"code"`
		Description *string  `json:"description"`
		Tags        []string `json:"tags"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Code != nil {
		cheatcode.Code = *input.Code
	}
	if input.Description != nil {
		cheatcode.Description = *input.Description
	}
	if input.Tags != nil {
		cheatcode.Tags = input.Tags
	}

	v := validator.New()

	if data.ValidateCheatcode(v, cheatcode); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cheatcodes.Update(cheatcode)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cheatcode": cheatcode}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteCheatcodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Cheatcodes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "cheatcode successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCheatcodesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Code        string
		Description string
		Tags        []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Code = app.readString(qs, "code", "")
	input.Description = app.readString(qs, "description", "")
	input.Tags = app.readCSV(qs, "tags", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "code", "-id", "-code"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cheatcodes, err := app.models.Cheatcodes.GetAll(input.Code, input.Description, input.Tags, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cheatcodes": cheatcodes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
