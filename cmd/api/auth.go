package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/MohammadTaghipour/social/internal/mailer"
	"github.com/MohammadTaghipour/social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plan token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// store the user
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.statusBadRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.statusBadRequestError(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	userWithToken := &UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "prod"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// send activation email
	err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.statusInternalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload credentials
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.statusBadRequestError(w, r, err)
		return
	}

	// fetch the user (check if the user exists) from the payload
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedError(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.jwt.expiration).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.jwt.issuer,
		"aud": app.config.auth.jwt.issuer,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	// send it to the client
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
