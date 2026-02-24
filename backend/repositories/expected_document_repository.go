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
		INSERT INTO expected_documents(document_request_id, title, description, is_uploaded) 
		VALUES($1, $2, $3, $4) RETURNING id;
	`

	var insertedID int

	err := r.db.QueryRowContext(
		ctx,
		query,
		requestID,
		expectedDocument.Title,
		expectedDocument.Description,
		expectedDocument.IsUploaded,
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (repo *ExpectedDocumentRepository) AddExpectedDocumentToRequestWithTx(ctx context.Context, tx *sql.Tx, ed models.ExpectedDocument) error {
	query := "INSERT INTO expected_documents (document_request_id, title, description, is_uploaded) VALUES ($1, $2, $3, $4)"
	_, err := tx.ExecContext(ctx, query, ed.DocumentRequestID, ed.Title, ed.Description, ed.IsUploaded)
	return err
}

func (r *ExpectedDocumentRepository) GetExpectedDocumentsByRequest(ctx context.Context, requestID int) ([]models.ExpectedDocument, error) {
	query := `
		SELECT id, document_request_id, title, description, is_uploaded FROM expected_documents WHERE document_request_id = $1
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
			&expectedDoc.IsUploaded,
		)

		if err != nil {
			return nil, err
		}

		expectedDocuments = append(expectedDocuments, expectedDoc)
	}

	return expectedDocuments, rows.Err()
}

func (r *ExpectedDocumentRepository) MarkAsUploaded(ctx context.Context, expectedDocumentID int) error {
	query := `
		UPDATE expected_documents SET is_uploaded = true WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, expectedDocumentID)
	return err
}

func (r *ExpectedDocumentRepository) DeleteExpectedDocumentFromRequest(ctx context.Context, requestId int, expectedDocumentId int) error {
	query := `
		DELETE FROM expected_documents WHERE id=$1 AND document_request_id=$2
	`

	_, err := r.db.ExecContext(ctx, query, expectedDocumentId, requestId)

	return err
}
