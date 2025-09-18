package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"` // Exclude password from JSON responses
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	p.text = &text
	p.hash = hash
	return err
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	switch {
	case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
		return ErrDuplicateEmail
	case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
		return ErrDuplicateUsername
	default:
		return err
	}
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string,
	invitationsExpDate time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create the user invite
		return s.createUserInvitation(ctx, tx, token, invitationsExpDate, user.ID)

	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string,
	invitationsExpDate time.Duration, userID int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationsExpDate))

	return err
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`
	var user User
	if err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
