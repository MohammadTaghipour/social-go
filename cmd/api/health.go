package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Healthcheck for API
//	@Description	Healthcheck for API
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	"Status OK"
//	@Failure		500	{object}	error	"Internal Error"
//	@Security		ApiKeyAuth
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
