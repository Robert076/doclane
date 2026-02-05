package repositories

import (
	"context"
	"database/sql"

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
		`INSERT INTO document_requests (professional_id, client_id, title, description, is_recurring, recurrence_cron, last_uploaded_at, due_date, next_due_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		req.ProfessionalID, req.ClientID, req.Title, req.Description, req.IsRecurring, req.RecurrenceCron, req.LastUploadedAt, req.DueDate, req.NextDueAt,
	).Scan(&id)
	return id, err
}

func (r *DocumentRepository) GetDocumentRequestByID(ctx context.Context, id int) (models.DocumentRequestDTORead, error) {
	var req models.DocumentRequestDTORead
	query := `
        SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.last_uploaded_at, u.email as client_email, 
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
		&req.LastUploadedAt,
		&req.ClientEmail,
		&req.Title,
		&req.Description,
		&req.DueDate,
		&req.NextDueAt,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	return req, err
}

func (r *DocumentRepository) GetDocumentRequestsByProfessional(ctx context.Context, professionalID int) ([]models.DocumentRequestDTORead, error) {
	query := `
        SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.last_uploaded_at, u.email as client_email, 
               dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at
        FROM document_requests dr
        JOIN users u ON dr.client_id = u.id
        WHERE dr.professional_id=$1
        ORDER BY dr.created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, professionalID)
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
			&req.LastUploadedAt,
			&req.ClientEmail,
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

func (r *DocumentRepository) GetDocumentRequestsByClient(ctx context.Context, clientID int) ([]models.DocumentRequestDTORead, error) {
	query := `
        SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.last_uploaded_at, u.email as client_email, 
               dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at
        FROM document_requests dr
        JOIN users u ON dr.client_id = u.id
        WHERE dr.client_id=$1
        ORDER BY dr.created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, clientID)
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
			&req.LastUploadedAt,
			&req.ClientEmail,
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

func (r *DocumentRepository) GetFileByID(ctx context.Context, id int) (models.DocumentFile, error) {
	var f models.DocumentFile
	query := `
        SELECT id, document_request_id, file_name, file_path, mime_type, file_size, uploaded_at, s3_version_id
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
	)
	return f, err
}

func (r *DocumentRepository) GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentFile, error) {
	query := `
        SELECT id, document_request_id, file_name, file_path, mime_type, file_size, uploaded_at, s3_version_id
        FROM document_files
        WHERE document_request_id=$1
        ORDER BY uploaded_at ASC
    `
	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.DocumentFile
	for rows.Next() {
		var f models.DocumentFile
		if err := rows.Scan(
			&f.ID,
			&f.DocumentRequestID,
			&f.FileName,
			&f.FilePath,
			&f.MimeType,
			&f.FileSize,
			&f.UploadedAt,
			&f.S3VersionID,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}

func (r *DocumentRepository) UpdateDocumentRequestStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE document_requests SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
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
            uploaded_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
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
	if err != nil {
		return err
	}

	return nil
}
