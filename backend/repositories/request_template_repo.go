package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Robert076/doclane/backend/models"
)

type RequestTemplateRepo struct {
	db *sql.DB
}

func NewRequestTemplateRepo(db *sql.DB) *RequestTemplateRepo {
	return &RequestTemplateRepo{db: db}
}

func (r *RequestTemplateRepo) GetRequestTemplatesByProfessionalID(ctx context.Context, professionalID int) ([]models.RequestTemplate, error) {
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

	templates := make([]models.RequestTemplate, 0)
	for rows.Next() {
		var t models.RequestTemplate
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

func (r *RequestTemplateRepo) GetRequestTemplateByID(ctx context.Context, id int) (models.RequestTemplate, error) {
	var t models.RequestTemplate
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

func (r *RequestTemplateRepo) AddRequestTemplate(ctx context.Context, tmp models.RequestTemplate) (int, error) {
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

func (r *RequestTemplateRepo) AddRequestTemplateWithTx(ctx context.Context, tx *sql.Tx, tmp models.RequestTemplate) (int, error) {
	var id int
	query := `
        INSERT INTO document_request_templates (title, description, is_recurring, recurrence_cron, created_by, is_closed)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	err := tx.QueryRowContext(ctx, query,
		tmp.Title,
		tmp.Description,
		tmp.IsRecurring,
		tmp.RecurrenceCron,
		tmp.CreatedBy,
		tmp.IsClosed,
	).Scan(&id)
	return id, err
}

func (r *RequestTemplateRepo) CloseRequestTemplate(ctx context.Context, id int) error {
	query := `
		UPDATE document_request_templates SET is_closed=true WHERE id=$1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RequestTemplateRepo) ReopenRequestTemplate(ctx context.Context, id int) error {
	query := `
		UPDATE document_request_templates SET is_closed=false WHERE id=$1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RequestTemplateRepo) DeleteRequestTemplate(ctx context.Context, id int) error {
	query := `
		DELETE FROM document_request_templates WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RequestTemplateRepo) PatchRequestTemplate(ctx context.Context, templateID int, dto models.RequestTemplateDTOPatch) error {
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if dto.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *dto.Title)
		argIdx++
	}

	if dto.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *dto.Description)
		argIdx++
	}

	if dto.IsRecurring != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_recurring = $%d", argIdx))
		args = append(args, *dto.IsRecurring)
		argIdx++
	}

	if dto.RecurrenceCron != nil {
		setClauses = append(setClauses, fmt.Sprintf("recurrence_cron = $%d", argIdx))
		args = append(args, *dto.RecurrenceCron)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now())
	argIdx++

	args = append(args, templateID)
	query := fmt.Sprintf(
		"UPDATE document_request_templates SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		argIdx,
	)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
