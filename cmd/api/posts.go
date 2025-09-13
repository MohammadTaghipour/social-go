package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MohammadTaghipour/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// createPostHandler godoc
//
//	@Summary		Create a new post
//	@Description	Creates a new post with title, content, and optional tags
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			post	body		CreatePostPayload	true	"Post data"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error	"Invalid request payload"
//	@Failure		500		{object}	error	"Internal Server Error"
//	@Router			/post/create [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  1, // TODO: change after auth
		Tags:    payload.Tags,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// getPostHandler godoc
//
//	@Summary		Get a post
//	@Description	Returns a single post with metadata and comments
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		500		{object}	error	"Internal Server Error"
//	@Router			/post/{postID}/ [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	comments, err := app.store.Comments.GetByPostID(ctx, post.ID)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// deletePostHandler godoc
//
//	@Summary		Delete a post
//	@Description	Deletes a post by its ID
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int	true	"Post ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	error	"Bad Request"
//	@Failure		404		{object}	error	"Post Not Found"
//	@Failure		500		{object}	error	"Internal Server Error"
//	@Router			/post/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.statusNotFoundError(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags"`
}

// updatePostHandler godoc
//
//	@Summary		Update an existing post
//	@Description	Updates fields of an existing post (title, content, tags)
//	@Tags			post
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"
//	@Param			post	body		UpdatePostPayload	true	"Updated post data"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error	"Invalid request payload"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error	"Internal Server Error"
//	@Router			/post/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.statusInternalServerError(w, r, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.statusNotFoundError(w, r, err)
			default:
				app.statusInternalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
