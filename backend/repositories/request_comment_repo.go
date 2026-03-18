package repositories

import (
	"context"
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type RequestCommentRepo struct {
	db *sql.DB
}

func NewRequestCommentRepo(db *sql.DB) *RequestCommentRepo {
	return &RequestCommentRepo{
		db: db,
	}
}

func (r *RequestCommentRepo) GetCommentsByRequestID(ctx context.Context, requestID int) ([]models.RequestCommentDTO, error) {
	query := `
		SELECT c.id, c.request_id, c.user_id, c.comment, c.created_at, c.updated_at, u.first_name, u.last_name
		FROM document_comments c JOIN users u ON c.user_id = u.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := make([]models.RequestCommentDTO, 0)
	for rows.Next() {
		var c models.RequestCommentDTO
		err := rows.Scan(
			&c.ID,
			&c.RequestID,
			&c.UserID,
			&c.Comment,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.UserFirstName,
			&c.UserLastName,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, rows.Err()
}

func (r *RequestCommentRepo) GetCommentByID(ctx context.Context, commentID int) (models.RequestCommentDTO, error) {
	query := `
		SELECT c.id, c.request_id, c.user_id, c.comment, c.created_at, c.updated_at, u.first_name, u.last_name
		FROM document_comments c JOIN users u ON c.user_id = u.id WHERE c.id = $1
	`

	row, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return models.RequestCommentDTO{}, err
	}
	defer row.Close()

	var c models.RequestCommentDTO
	err = row.Scan(
		&c.ID,
		&c.RequestID,
		&c.UserID,
		&c.Comment,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.UserFirstName,
		&c.UserLastName,
	)

	return c, err
}

func (r *RequestCommentRepo) AddComment(ctx context.Context, comment models.RequestComment) (int, error) {
	query := `
		INSERT INTO document_comments(request_id, user_id, comment, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		comment.RequestID,
		comment.UserID,
		comment.Comment,
		comment.CreatedAt,
		comment.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
