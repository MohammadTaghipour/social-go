package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt  int64 `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followedID, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`
	_, err := s.db.ExecContext(ctx, query, followedID, userID)
	return err
}

func (s *FollowerStore) UnFollow(ctx context.Context, followedID, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`
	_, err := s.db.ExecContext(ctx, query, followedID, userID)
	return err
}
