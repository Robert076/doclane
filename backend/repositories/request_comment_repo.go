package repositories

import (
	"context"
	"database/sql"
	"time"

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
		FROM document_comments c JOIN users u ON c.user_id = u.id and c.request_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, requestID)
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

	return id, err
}

func (r *RequestCommentRepo) GetLastCommentFromUser(ctx context.Context, userID int) (models.RequestComment, error) {
	query := `
        SELECT c.id, c.request_id, c.user_id, c.comment, c.created_at, c.updated_at 
        FROM document_comments c
        WHERE c.user_id = $1 ORDER BY c.created_at DESC LIMIT 1
    `
	var comm models.RequestComment
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&comm.ID,
		&comm.RequestID,
		&comm.UserID,
		&comm.Comment,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return comm, err
	}
	comm.CreatedAt, _ = time.ParseInLocation(time.RFC3339Nano, createdAt, time.UTC)
	comm.UpdatedAt, _ = time.ParseInLocation(time.RFC3339Nano, updatedAt, time.UTC)
	return comm, nil
}
