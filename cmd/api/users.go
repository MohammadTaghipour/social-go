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

type userKey string

const userCtxKey userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user Profile
//	@Description	Returns a user's info
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.statusInternalServerError(w, r, err)
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		204		"User followed"
//	@Failure		400		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/user/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	// TODO: change this after auth
	var userID int64 = 1

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, userID); err != nil {
		app.statusInternalServerError(w, r, err)

	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.statusInternalServerError(w, r, err)
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user
//	@Description	UnFollows a user by ID
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		204		"User unfollowed"
//	@Failure		400		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/user/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	// TODO: change this after auth
	var userID int64 = 1

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	if err := app.store.Followers.UnFollow(ctx, followerUser.ID, userID); err != nil {
		app.statusInternalServerError(w, r, err)

	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.statusInternalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.statusBadRequestError(w, r, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.statusNotFoundError(w, r, err)
			default:
				app.statusInternalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)
	return user
}
