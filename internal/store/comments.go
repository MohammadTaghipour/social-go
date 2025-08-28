package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at 
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	return err
}
