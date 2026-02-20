package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/Robert076/doclane/backend/models"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (repo *DocumentRepository) AddDocumentRequest(ctx context.Context, req models.DocumentRequest) (int, error) {
	var id int
	err := repo.db.QueryRowContext(ctx,
		`INSERT INTO document_requests (professional_id, client_id, title, description, is_recurring, recurrence_cron, is_scheduled, scheduled_for, is_closed, last_uploaded_at, due_date, next_due_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		req.ProfessionalID, req.ClientID, req.Title, req.Description, req.IsRecurring, req.RecurrenceCron, req.IsScheduled, req.ScheduledFor, req.IsClosed, req.LastUploadedAt, req.DueDate, req.NextDueAt,
	).Scan(&id)
	return id, err
}

func (r *DocumentRepository) GetDocumentRequestByID(ctx context.Context, id int) (models.DocumentRequestDTORead, error) {
	var req models.DocumentRequestDTORead
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.is_scheduled, dr.scheduled_for, dr.is_closed, dr.last_uploaded_at, u.email as client_email, u.first_name, u.last_name, 
		dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at
		FROM document_requests dr
		JOIN users u ON dr.client_id = u.id
		WHERE dr.id=$1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&req.ID,
		&req.ProfessionalID,
		&req.ClientID,
		&req.IsRecurring,
		&req.RecurrenceCron,
		&req.IsScheduled,
		&req.ScheduledFor,
		&req.IsClosed,
		&req.LastUploadedAt,
		&req.ClientEmail,
		&req.ClientFirstName,
		&req.ClientLastName,
		&req.Title,
		&req.Description,
		&req.DueDate,
		&req.NextDueAt,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	return req, err
}

func (r *DocumentRepository) GetDocumentRequestsByProfessional(
	ctx context.Context,
	professionalID int,
	search *string,
) ([]models.DocumentRequestDTORead, error) {
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.is_scheduled, dr.scheduled_for, dr.is_closed,
		dr.last_uploaded_at, u.email, u.first_name, u.last_name, 
		dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at
		FROM document_requests dr
		JOIN users u ON dr.client_id = u.id
		WHERE dr.professional_id=$1
	`

	args := []interface{}{professionalID}
	argIndex := 2

	if search != nil && *search != "" {
		searchPattern := "%" + strings.ToLower(*search) + "%"
		query += ` AND (
            LOWER(dr.title) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(dr.description) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(u.email) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(u.first_name) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(u.last_name) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(u.first_name || ' ' || u.last_name) LIKE $` + strconv.Itoa(argIndex) + `
        )`
		args = append(args, searchPattern)
	}

	query += " ORDER BY dr.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.DocumentRequestDTORead
	for rows.Next() {
		var req models.DocumentRequestDTORead
		err := rows.Scan(
			&req.ID,
			&req.ProfessionalID,
			&req.ClientID,
			&req.IsRecurring,
			&req.RecurrenceCron,
			&req.IsScheduled,
			&req.ScheduledFor,
			&req.IsClosed,
			&req.LastUploadedAt,
			&req.ClientEmail,
			&req.ClientFirstName,
			&req.ClientLastName,
			&req.Title,
			&req.Description,
			&req.DueDate,
			&req.NextDueAt,
			&req.CreatedAt,
			&req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *DocumentRepository) GetDocumentRequestsByClient(
	ctx context.Context,
	clientID int,
	search *string,
) ([]models.DocumentRequestDTORead, error) {
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.is_scheduled, dr.scheduled_for, dr.is_closed,
			dr.last_uploaded_at, u.email, u.first_name, u.last_name, 
			dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at
		FROM document_requests dr
		JOIN users u ON dr.client_id = u.id
		WHERE dr.client_id=$1 
		AND (dr.is_scheduled = false OR dr.scheduled_for <= NOW())
	`

	args := []interface{}{clientID}
	argIndex := 2

	if search != nil && *search != "" {
		searchPattern := "%" + strings.ToLower(*search) + "%"
		query += ` AND (
            LOWER(dr.title) LIKE $` + strconv.Itoa(argIndex) + ` OR
            LOWER(dr.description) LIKE $` + strconv.Itoa(argIndex) + `
        )`
		args = append(args, searchPattern)
	}

	query += " ORDER BY dr.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.DocumentRequestDTORead
	for rows.Next() {
		var req models.DocumentRequestDTORead
		if err := rows.Scan(
			&req.ID,
			&req.ProfessionalID,
			&req.ClientID,
			&req.IsRecurring,
			&req.RecurrenceCron,
			&req.IsScheduled,
			&req.ScheduledFor,
			&req.IsClosed,
			&req.LastUploadedAt,
			&req.ClientEmail,
			&req.ClientFirstName,
			&req.ClientLastName,
			&req.Title,
			&req.Description,
			&req.DueDate,
			&req.NextDueAt,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *DocumentRepository) CloseDocumentRequest(ctx context.Context, id int) error {
	query := `
		UPDATE document_requests SET is_closed=true WHERE id=$1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *DocumentRepository) UpdateDocumentRequestTitle(ctx context.Context, id int, newTitle string) error {
	query := `
		UPDATE document_requests SET title=$1 WHERE id=$2
	`
	_, err := r.db.ExecContext(ctx, query, newTitle, id)
	return err
}

func (r *DocumentRepository) GetFileByID(ctx context.Context, id int) (models.DocumentFile, error) {
	var f models.DocumentFile
	query := `
        SELECT id, document_request_id, file_name, file_path, mime_type, file_size, uploaded_at, s3_version_id, uploaded_by
        FROM document_files
        WHERE id=$1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.DocumentRequestID,
		&f.FileName,
		&f.FilePath,
		&f.MimeType,
		&f.FileSize,
		&f.UploadedAt,
		&f.S3VersionID,
		&f.UploadedBy,
	)
	return f, err
}

func (r *DocumentRepository) GetFileByIDExtended(ctx context.Context, id int) (models.DocumentFileDTOExtended, error) {
	var f models.DocumentFileDTOExtended
	query := `
        SELECT df.id, df.document_request_id, df.file_name, df.file_path, df.mime_type, df.file_size, df.uploaded_at, df.s3_version_id, df.uploaded_by, u.role
        FROM document_files df 
		JOIN users u ON u.id = df.uploaded_by
        WHERE df.id=$1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.DocumentRequestID,
		&f.FileName,
		&f.FilePath,
		&f.MimeType,
		&f.FileSize,
		&f.UploadedAt,
		&f.S3VersionID,
		&f.UploadedBy,
		&f.AuthorRole,
	)
	return f, err
}

func (r *DocumentRepository) GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentFileDTORead, error) {
	query := `
        SELECT df.id, df.document_request_id, df.file_name, df.file_path, df.mime_type, df.file_size, df.uploaded_at, df.s3_version_id, df.uploaded_by, u.first_name, u.last_name
        FROM document_files df
		JOIN users u ON u.id = df.uploaded_by
        WHERE document_request_id=$1
        ORDER BY uploaded_at ASC
    `
	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.DocumentFileDTORead
	for rows.Next() {
		var f models.DocumentFileDTORead
		if err := rows.Scan(
			&f.ID,
			&f.DocumentRequestID,
			&f.FileName,
			&f.FilePath,
			&f.MimeType,
			&f.FileSize,
			&f.UploadedAt,
			&f.S3VersionID,
			&f.UploadedBy,
			&f.UploadedByFirstName,
			&f.UploadedByLastName,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}

func (r *DocumentRepository) AddDocumentFile(ctx context.Context, file models.DocumentFile) (int, error) {
	var id int

	query := `
        INSERT INTO document_files (
            document_request_id, 
            file_name, 
            file_path, 
            mime_type, 
            file_size, 
            s3_version_id, 
            uploaded_at,
			uploaded_by
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
		`
	// ON CONFLICT (document_request_id, file_name)
	// DO UPDATE SET
	//     file_path = EXCLUDED.file_path,
	//     mime_type = EXCLUDED.mime_type,
	//     file_size = EXCLUDED.file_size,
	//     s3_version_id = EXCLUDED.s3_version_id,
	//     uploaded_at = EXCLUDED.uploaded_at
	// RETURNING id

	err := r.db.QueryRowContext(ctx,
		query,
		file.DocumentRequestID,
		file.FileName,
		file.FilePath,
		file.MimeType,
		file.FileSize,
		file.S3VersionID,
		file.UploadedAt,
		file.UploadedBy,
	).Scan(&id)

	return id, err
}

func (r *DocumentRepository) SetFileUploaded(ctx context.Context, id int) error {
	query := `
		UPDATE document_requests
		SET last_uploaded_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
