package repositories

import (
	"context"
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type ExpectedDocumentTemplateRepository struct {
	db *sql.DB
}

func NewExpectedDocumentTemplateRepository(db *sql.DB) *ExpectedDocumentTemplateRepository {
	return &ExpectedDocumentTemplateRepository{db: db}
}

func (r *ExpectedDocumentTemplateRepository) GetByTemplateID(ctx context.Context, templateID int) ([]models.ExpectedDocumentTemplate, error) {
	query := `
		SELECT id, document_request_template_id, title, description, example_file_path, example_mime_type
		FROM expected_document_templates
		WHERE document_request_template_id = $1
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]models.ExpectedDocumentTemplate, 0)
	for rows.Next() {
		var t models.ExpectedDocumentTemplate
		if err := rows.Scan(
			&t.ID,
			&t.DocumentRequestTemplateID,
			&t.Title,
			&t.Description,
			&t.ExampleFilePath,
			&t.ExampleMimeType,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *ExpectedDocumentTemplateRepository) Add(ctx context.Context, t models.ExpectedDocumentTemplate) (int, error) {
	var id int
	query := `
        INSERT INTO expected_document_templates (document_request_template_id, title, description, example_file_path, example_mime_type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query,
		t.DocumentRequestTemplateID,
		t.Title,
		t.Description,
		t.ExampleFilePath,
		t.ExampleMimeType,
	).Scan(&id)
	return id, err
}

func (r *ExpectedDocumentTemplateRepository) DeleteByID(ctx context.Context, id int) error {
	query := `
        DELETE FROM expected_document_templates WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ExpectedDocumentTemplateRepository) GetByDocumentID(ctx context.Context, id int) (models.ExpectedDocumentTemplate, error) {
	query := `
		SELECT id, document_request_template_id, title, description, example_file_path, example_mime_type
		FROM expected_document_templates
		WHERE id = $1
	`

	var t models.ExpectedDocumentTemplate

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.DocumentRequestTemplateID,
		&t.Title,
		&t.Description,
		&t.ExampleFilePath,
		&t.ExampleMimeType,
	)

	if err != nil {
		return models.ExpectedDocumentTemplate{}, err
	}

	return t, nil
}
