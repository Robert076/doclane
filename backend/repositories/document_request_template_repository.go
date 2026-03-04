package repositories

import (
	"context"
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type DocumentRequestTemplateRepository struct {
	db *sql.DB
}

func NewDocumentRequestTemplateRepository(db *sql.DB) *DocumentRequestTemplateRepository {
	return &DocumentRequestTemplateRepository{db: db}
}

func (r *DocumentRequestTemplateRepository) GetDocumentRequestTemplatesByProfessionalID(ctx context.Context, professionalID int) ([]models.DocumentRequestTemplate, error) {
	query := `
        SELECT id, title, description, is_recurring, recurrence_cron, created_by, created_at, updated_at, is_closed
        FROM document_request_templates
        WHERE created_by = $1
        ORDER BY created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, professionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]models.DocumentRequestTemplate, 0)
	for rows.Next() {
		var t models.DocumentRequestTemplate
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.IsRecurring,
			&t.RecurrenceCron,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.IsClosed,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *DocumentRequestTemplateRepository) GetDocumentRequestTemplateByID(ctx context.Context, id int) (models.DocumentRequestTemplate, error) {
	var t models.DocumentRequestTemplate
	query := `
        SELECT id, title, description, is_recurring, recurrence_cron, created_by, created_at, updated_at, is_closed
        FROM document_request_templates
        WHERE id = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.IsRecurring,
		&t.RecurrenceCron,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.IsClosed,
	)
	return t, err
}

func (r *DocumentRequestTemplateRepository) AddDocumentRequestTemplate(ctx context.Context, tmp models.DocumentRequestTemplate) (int, error) {
	var id int
	query := `
		INSERT INTO document_request_templates (title, description, is_recurring, recurrence_cron, created_by, is_closed)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		tmp.Title,
		tmp.Description,
		tmp.IsRecurring,
		tmp.RecurrenceCron,
		tmp.CreatedBy,
		tmp.IsClosed,
	).Scan(&id)
	return id, err
}

func (r *DocumentRequestTemplateRepository) CloseDocumentRequestTemplate(ctx context.Context, id int) error {
	query := `
		UPDATE document_request_templates SET is_closed=true WHERE id=$1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *DocumentRequestTemplateRepository) ReopenDocumentRequestTemplate(ctx context.Context, id int) error {
	query := `
		UPDATE document_request_templates SET is_closed=false WHERE id=$1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
