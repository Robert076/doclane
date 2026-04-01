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

func (r *RequestTemplateRepo) GetRequestTemplatesByDepartment(ctx context.Context, departmentID int) ([]models.RequestTemplateDTORead, error) {
	query := `
		SELECT t.id, t.title, t.description, t.department_id, t.is_recurring, t.recurrence_cron,
			t.created_by, t.created_at, t.updated_at, t.is_closed,
			u.first_name, u.last_name
		FROM document_request_templates t
		JOIN users u ON u.id = t.created_by
		WHERE t.department_id = $1
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]models.RequestTemplateDTORead, 0)
	for rows.Next() {
		var t models.RequestTemplateDTORead
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.DepartmentID,
			&t.IsRecurring,
			&t.RecurrenceCron,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.IsClosed,
			&t.AuthorFirstName,
			&t.AuthorLastName,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *RequestTemplateRepo) GetAllRequestTemplates(ctx context.Context) ([]models.RequestTemplateDTORead, error) {
	query := `
		SELECT t.id, t.title, t.description, t.department_id, t.is_recurring, t.recurrence_cron,
			t.created_by, t.created_at, t.updated_at, t.is_closed,
			u.first_name, u.last_name
		FROM document_request_templates t
		JOIN users u ON u.id = t.created_by
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]models.RequestTemplateDTORead, 0)
	for rows.Next() {
		var t models.RequestTemplateDTORead
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.DepartmentID,
			&t.IsRecurring,
			&t.RecurrenceCron,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.IsClosed,
			&t.AuthorFirstName,
			&t.AuthorLastName,
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
		SELECT id, title, description, department_id, is_recurring, recurrence_cron, created_by, created_at, updated_at, is_closed
		FROM document_request_templates
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.DepartmentID,
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
		INSERT INTO document_request_templates (title, description, department_id, is_recurring, recurrence_cron, created_by, is_closed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		tmp.Title,
		tmp.Description,
		tmp.DepartmentID,
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
		INSERT INTO document_request_templates (title, description, department_id, is_recurring, recurrence_cron, created_by, is_closed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := tx.QueryRowContext(ctx, query,
		tmp.Title,
		tmp.Description,
		tmp.DepartmentID,
		tmp.IsRecurring,
		tmp.RecurrenceCron,
		tmp.CreatedBy,
		tmp.IsClosed,
	).Scan(&id)
	return id, err
}

func (r *RequestTemplateRepo) CloseRequestTemplate(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE document_request_templates SET is_closed=true WHERE id=$1`, id)
	return err
}

func (r *RequestTemplateRepo) ReopenRequestTemplate(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE document_request_templates SET is_closed=false WHERE id=$1`, id)
	return err
}

func (r *RequestTemplateRepo) DeleteRequestTemplate(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM document_request_templates WHERE id = $1`, id)
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
