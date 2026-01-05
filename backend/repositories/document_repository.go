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
		`INSERT INTO document_requests (professional_id, client_id, title, description, due_date, status)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		req.ProfessionalID, req.ClientID, req.Title, req.Description, req.DueDate, req.Status,
	).Scan(&id)
	return id, err
}

func (r *DocumentRepository) GetDocumentRequestByID(ctx context.Context, id int) (models.DocumentRequestDTO, error) {
	var req models.DocumentRequestDTO
	query := `
        SELECT dr.id, dr.professional_id, dr.client_id, u.email as client_email, 
               dr.title, dr.description, dr.due_date, dr.status, dr.created_at, dr.updated_at
        FROM document_requests dr
        JOIN users u ON dr.client_id = u.id
        WHERE dr.id=$1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&req.ID,
		&req.ProfessionalID,
		&req.ClientID,
		&req.ClientEmail,
		&req.Title,
		&req.Description,
		&req.DueDate,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	return req, err
}

func (r *DocumentRepository) GetDocumentRequestsByProfessional(ctx context.Context, professionalID int) ([]models.DocumentRequestDTO, error) {
	query := `
        SELECT dr.id, dr.professional_id, dr.client_id, u.email as client_email, 
               dr.title, dr.description, dr.due_date, dr.status, dr.created_at, dr.updated_at
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

	var requests []models.DocumentRequestDTO
	for rows.Next() {
		var req models.DocumentRequestDTO
		err := rows.Scan(
			&req.ID,
			&req.ProfessionalID,
			&req.ClientID,
			&req.ClientEmail,
			&req.Title,
			&req.Description,
			&req.DueDate,
			&req.Status,
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

func (r *DocumentRepository) GetDocumentRequestsByClient(ctx context.Context, clientID int) ([]models.DocumentRequestDTO, error) {
	query := `
        SELECT dr.id, dr.professional_id, dr.client_id, u.email as client_email, 
               dr.title, dr.description, dr.due_date, dr.status, dr.created_at, dr.updated_at
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

	var requests []models.DocumentRequestDTO
	for rows.Next() {
		var req models.DocumentRequestDTO
		if err := rows.Scan(
			&req.ID,
			&req.ProfessionalID,
			&req.ClientID,
			&req.ClientEmail,
			&req.Title,
			&req.Description,
			&req.DueDate,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, rows.Err()
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
        ON CONFLICT (document_request_id, file_name) 
        DO UPDATE SET 
            file_path = EXCLUDED.file_path,
            mime_type = EXCLUDED.mime_type,
            file_size = EXCLUDED.file_size,
            s3_version_id = EXCLUDED.s3_version_id,
            uploaded_at = EXCLUDED.uploaded_at
        RETURNING id
    `

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
func (r *DocumentRepository) GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentFile, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, document_request_id, file_name, file_path, mime_type, file_size, uploaded_at
        FROM document_files
        WHERE document_request_id=$1
        ORDER BY uploaded_at ASC
    `, requestID)
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
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, rows.Err()
}
