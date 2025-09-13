package main

import (
	"context"
	"net/http"
	"time"

	"github.com/MohammadTaghipour/social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Returns paginated posts for a given user (with filters, tags, etc.)
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query	int			false	"Max items per page"
//	@Param			offset	query	int			false	"Pagination offset"
//	@Param			sort	query	string		false	"Sort order (asc or desc)"
//	@Param			tags	query	[]string	false	"Filter by tags (comma separated)"
//	@Param			search	query	string		false	"Full-text search in title/content"
//	@Param			since	query	string		false	"Filter posts created since (RFC3339)"
//	@Param			until	query	string		false	"Filter posts created until (RFC3339)"
//	@Success		200		{array}	store.PostWithMetadata
//	@Failure		400
//	@Failure		500
//	@Router			/feed [get]
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
