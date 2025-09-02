package main

import (
	"context"
	"net/http"
	"time"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filterns and etc

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(1)) // TODO: get userID from auth
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
