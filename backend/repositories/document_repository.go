package repositories

import (
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (repo *DocumentRepository) CreateDocumentRequest(req models.DocumentRequest) (int, error) {
	var id int
	err := repo.db.QueryRow(
		`INSERT INTO document_requests (professional_id, client_id, title, description, due_date, status)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		req.ProfessionalID, req.ClientID, req.Title, req.Description, req.DueDate, req.Status,
	).Scan(&id)
	return id, err
}

func (r *DocumentRepository) GetDocumentRequestByID(id int) (models.DocumentRequest, error) {
	var req models.DocumentRequest
	row := r.db.QueryRow(`
        SELECT id, professional_id, client_id, title, description, due_date, status, created_at, updated_at
        FROM document_requests WHERE id=$1
    `, id)

	err := row.Scan(
		&req.ID,
		&req.ProfessionalID,
		&req.ClientID,
		&req.Title,
		&req.Description,
		&req.DueDate,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	return req, err
}

func (r *DocumentRepository) GetDocumentRequestsByClient(clientID int) ([]models.DocumentRequest, error) {
	rows, err := r.db.Query(`
        SELECT id, professional_id, client_id, title, description, due_date, status, created_at, updated_at
        FROM document_requests
        WHERE client_id=$1
        ORDER BY created_at DESC
    `, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.DocumentRequest
	for rows.Next() {
		var req models.DocumentRequest
		if err := rows.Scan(
			&req.ID,
			&req.ProfessionalID,
			&req.ClientID,
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

func (r *DocumentRepository) GetDocumentRequestsByProfessional(professionalID int) ([]models.DocumentRequest, error) {
	rows, err := r.db.Query(`
        SELECT id, professional_id, client_id, title, description, due_date, status, created_at, updated_at
        FROM document_requests
        WHERE professional_id=$1
        ORDER BY created_at DESC
    `, professionalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.DocumentRequest
	for rows.Next() {
		var req models.DocumentRequest
		err := rows.Scan(
			&req.ID,
			&req.ProfessionalID,
			&req.ClientID,
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

func (r *DocumentRepository) UpdateDocumentRequestStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE document_requests SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

func (r *DocumentRepository) AddDocumentFile(file models.DocumentFile) (int, error) {
	var id int
	err := r.db.QueryRow(`
        INSERT INTO document_files (document_request_id, file_name, file_path, mime_type, file_size)
        VALUES ($1,$2,$3,$4,$5)
        RETURNING id
    `,
		file.DocumentRequestID,
		file.FileName,
		file.FilePath,
		file.MimeType,
		file.FileSize,
	).Scan(&id)

	return id, err
}

func (r *DocumentRepository) GetFilesByRequest(requestID int) ([]models.DocumentFile, error) {
	rows, err := r.db.Query(`
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
