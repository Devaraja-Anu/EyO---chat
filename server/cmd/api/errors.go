package main

import "net/http"

type errorResponse struct {
	Error string `json:"error"`
}

func (app *application) errorResponse(w http.ResponseWriter, status int, message string) {
	env := envelope{"error": message}
	_ = app.writeJSON(w, status, env, nil)
}

func (app *application) badRequestResponse(w http.ResponseWriter, message string) {
	app.errorResponse(w, http.StatusBadRequest, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusNotFound, "resource not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, message string) {
	app.errorResponse(w, http.StatusConflict, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter,err error) {
	app.logger.Error("internal server error", map[string]any{
		"error": err.Error(),
	})
	app.errorResponse(w, http.StatusInternalServerError, "the  server encountered an issue")
}
