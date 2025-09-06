package main

import (
	"context"
	"net/http"
	"time"

	"github.com/MohammadTaghipour/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// default
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
		Since:  "",
		Until:  "",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.statusBadRequestError(w, r, err)
		return

	}

	if err := validate.Struct(fq); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(1), fq) // TODO: get userID from auth
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
