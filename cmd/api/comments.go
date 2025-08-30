package main

import (
	"context"
	"net/http"
	"time"

	"github.com/MohammadTaghipour/social/internal/store"
)

type CreateCommentPayload struct {
	PostID  int64  `json:"post_id" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=500"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	comment := store.Comment{
		Content: payload.Content,
		PostID:  payload.PostID,
		UserID:  payload.UserID,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := app.store.Comments.Create(ctx, &comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

}
