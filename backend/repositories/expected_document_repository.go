package repositories

import (
	"context"
	"database/sql"

	"github.com/Robert076/doclane/backend/models"
)

type ExpectedDocumentRepository struct {
	db *sql.DB
}

func NewExpectedDocRepository(db *sql.DB) *ExpectedDocumentRepository {
	return &ExpectedDocumentRepository{
		db: db,
	}
}

func (r *ExpectedDocumentRepository) AddExpectedDocumentToRequest(ctx context.Context, requestID int, expectedDocument models.ExpectedDocument) (int, error) {
	query := `
		INSERT INTO expected_documents(document_request_id, title, description, status, rejection_reason, example_file_path, example_mime_type) 
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	var insertedID int
	err := r.db.QueryRowContext(
		ctx,
		query,
		requestID,
		expectedDocument.Title,
		expectedDocument.Description,
		expectedDocument.Status,
		expectedDocument.RejectionReason,
		expectedDocument.ExampleFilePath,
		expectedDocument.ExampleMimeType,
	).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (r *ExpectedDocumentRepository) AddExpectedDocumentToRequestWithTx(ctx context.Context, tx *sql.Tx, ed models.ExpectedDocument) error {
	query := `
		INSERT INTO expected_documents (document_request_id, title, description, status, rejection_reason, example_file_path, example_mime_type) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.ExecContext(ctx, query,
		ed.DocumentRequestID,
		ed.Title,
		ed.Description,
		ed.Status,
		ed.RejectionReason,
		ed.ExampleFilePath,
		ed.ExampleMimeType,
	)
	return err
}

func (r *ExpectedDocumentRepository) GetExpectedDocumentsByRequest(ctx context.Context, requestID int) ([]models.ExpectedDocument, error) {
	query := `
		SELECT id, document_request_id, title, description, status, rejection_reason, example_file_path, example_mime_type 
		FROM expected_documents 
		WHERE document_request_id = $1 
		ORDER BY id
	`
	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expectedDocuments []models.ExpectedDocument
	for rows.Next() {
		var expectedDoc models.ExpectedDocument
		err := rows.Scan(
			&expectedDoc.ID,
			&expectedDoc.DocumentRequestID,
			&expectedDoc.Title,
			&expectedDoc.Description,
			&expectedDoc.Status,
			&expectedDoc.RejectionReason,
			&expectedDoc.ExampleFilePath,
			&expectedDoc.ExampleMimeType,
		)
		if err != nil {
			return nil, err
		}
		expectedDocuments = append(expectedDocuments, expectedDoc)
	}
	return expectedDocuments, rows.Err()
}

func (r *ExpectedDocumentRepository) GetExpectedDocumentByID(ctx context.Context, id int) (models.ExpectedDocument, error) {
	query := `
		SELECT id, document_request_id, title, description, status, rejection_reason, example_file_path, example_mime_type 
		FROM expected_documents 
		WHERE id = $1
	`
	var ed models.ExpectedDocument
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ed.ID,
		&ed.DocumentRequestID,
		&ed.Title,
		&ed.Description,
		&ed.Status,
		&ed.RejectionReason,
		&ed.ExampleFilePath,
		&ed.ExampleMimeType,
	)
	if err != nil {
		return models.ExpectedDocument{}, err
	}
	return ed, nil
}

func (r *ExpectedDocumentRepository) UpdateExpectedDocumentStatus(ctx context.Context, documentID int, status string, rejectionReason *string) error {
	query := `
		UPDATE expected_documents 
		SET status = $1, rejection_reason = $2 
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, rejectionReason, documentID)
	return err
}

func (r *ExpectedDocumentRepository) DeleteExpectedDocumentFromRequest(ctx context.Context, requestId int, expectedDocumentId int) error {
	query := `
		DELETE FROM expected_documents WHERE id=$1 AND document_request_id=$2
	`
	_, err := r.db.ExecContext(ctx, query, expectedDocumentId, requestId)
	return err
}
