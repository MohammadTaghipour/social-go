package main

import (
	"net/http"
)

func (app *application) statusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal error", "method", r.Method, "path", r.URL.Path,
		"error", err)
	writeJSONError(w, http.StatusInternalServerError, "Something went wrong in server")
}

func (app *application) statusBadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad request error", "method", r.Method, "path", r.URL.Path,
		"error", err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) statusNotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not found error", "method", r.Method, "path", r.URL.Path,
		"error", err)
	writeJSONError(w, http.StatusNotFound, "Not found")
}
