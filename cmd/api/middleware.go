package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MohammadTaghipour/social/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ") // Basic sadkjgbsadfs...
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			// check the credentials
			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) JwtAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedJwtError(w, r, fmt.Errorf("authorization header missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ") // Bearer <token>
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.unauthorizedJwtError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			token := parts[1]
			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				app.unauthorizedJwtError(w, r, err)
				return
			}

			claims := jwtToken.Claims.(jwt.MapClaims)

			userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
			if err != nil {
				app.unauthorizedJwtError(w, r, err)
				return
			}

			ctx := r.Context()

			user, err := app.store.Users.GetByID(ctx, userID)
			if err != nil {
				app.unauthorizedJwtError(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (app *application) CheckPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		// if it belongs to user
		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		// role check
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.statusInternalServerError(w, r, err)
			return
		}
		if !allowed {
			app.statusForbiddenError(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) checkRolePrecedence(ctx context.Context,
	user *store.User, requiredRole string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, requiredRole)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}
