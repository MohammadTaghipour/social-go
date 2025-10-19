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

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized error", "method", r.Method, "path", r.URL.Path,
		"error", err)
	writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized basic error", "method", r.Method, "path", r.URL.Path,
		"error", err)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) unauthorizedJwtError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized token error", "method", r.Method, "path", r.URL.Path,
		"error", err)

	writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) statusForbiddenError(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden error", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusForbidden, "Forbidden")
}
