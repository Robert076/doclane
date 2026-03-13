package repositories

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/Robert076/doclane/backend/models"
)

type RequestRepo struct {
	db *sql.DB
}

func NewRequestRepo(db *sql.DB) *RequestRepo {
	return &RequestRepo{db: db}
}

func (repo *RequestRepo) AddRequest(ctx context.Context, req models.Request) (int, error) {
	var id int
	query := `
		INSERT INTO document_requests (professional_id, client_id, title, description, is_recurring, recurrence_cron, is_scheduled, scheduled_for, is_closed, last_uploaded_at, due_date, next_due_at, template_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id
	`
	err := repo.db.QueryRowContext(ctx,
		query,
		req.ProfessionalID, req.ClientID, req.Title, req.Description, req.IsRecurring, req.RecurrenceCron, req.IsScheduled, req.ScheduledFor, req.IsClosed, req.LastUploadedAt, req.DueDate, req.NextDueAt, req.RequestTemplateID,
	).Scan(&id)
	return id, err
}

func (repo *RequestRepo) AddRequestWithTx(ctx context.Context, req models.Request, transaction *sql.Tx) (int, error) {
	var id int
	query := `
		INSERT INTO document_requests (professional_id, client_id, title, description, is_recurring, recurrence_cron, is_scheduled, scheduled_for, next_due_at, due_date, created_at, updated_at, template_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id
	`
	err := transaction.QueryRowContext(ctx, query, req.ProfessionalID, req.ClientID, req.Title, req.Description, req.IsRecurring, req.RecurrenceCron, req.IsScheduled, req.ScheduledFor, req.NextDueAt, req.DueDate, req.CreatedAt, req.UpdatedAt, req.RequestTemplateID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RequestRepo) GetRequestByID(ctx context.Context, id int) (models.RequestDTORead, error) {
	var req models.RequestDTORead
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.is_scheduled, dr.scheduled_for, dr.is_closed, dr.last_uploaded_at, u.email as client_email, u.first_name, u.last_name, 
		dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at, dr.template_id
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
		&req.RequestTemplateID,
	)
	return req, err
}

func (r *RequestRepo) GetRequestsByProfessional(
	ctx context.Context,
	professionalID int,
	search *string,
) ([]models.RequestDTORead, error) {
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.is_scheduled, dr.scheduled_for, dr.is_closed,
		dr.last_uploaded_at, u.email, u.first_name, u.last_name, 
		dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at, dr.template_id
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

	requests := make([]models.RequestDTORead, 0)
	for rows.Next() {
		var req models.RequestDTORead
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
			&req.RequestTemplateID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, rows.Err()
}

func (r *RequestRepo) GetRequestsByClient(
	ctx context.Context,
	clientID int,
	search *string,
) ([]models.RequestDTORead, error) {
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, dr.is_scheduled, dr.scheduled_for, dr.is_closed,
			dr.last_uploaded_at, u.email, u.first_name, u.last_name, 
			dr.title, dr.description, dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at, dr.template_id
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

	var requests []models.RequestDTORead
	for rows.Next() {
		var req models.RequestDTORead
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
			&req.RequestTemplateID,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *RequestRepo) ReopenRequest(ctx context.Context, id int) error {
	query := `
		UPDATE document_requests SET is_closed=false WHERE id=$1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RequestRepo) CloseRequest(ctx context.Context, id int) error {
	query := `
		UPDATE document_requests SET is_closed=true WHERE id=$1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RequestRepo) UpdateRequestTitle(ctx context.Context, id int, newTitle string) error {
	query := `
		UPDATE document_requests SET title=$1 WHERE id=$2
	`
	_, err := r.db.ExecContext(ctx, query, newTitle, id)
	return err
}

func (r *RequestRepo) GetFileByID(ctx context.Context, id int) (models.Document, error) {
	var f models.Document
	query := `
        SELECT id, document_request_id, expected_document_id, file_name, file_path, mime_type, file_size, uploaded_at, s3_version_id, uploaded_by
        FROM document_files
        WHERE id=$1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.RequestID,
		&f.ExpectedDocumentID,
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

func (r *RequestRepo) GetFileByIDExtended(ctx context.Context, id int) (models.DocumentDTOExtended, error) {
	var f models.DocumentDTOExtended
	query := `
        SELECT df.id, df.document_request_id, df.expected_document_id, df.file_name, df.file_path, df.mime_type, df.file_size, df.uploaded_at, df.s3_version_id, df.uploaded_by, u.role
        FROM document_files df 
		JOIN users u ON u.id = df.uploaded_by
        WHERE df.id=$1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.RequestID,
		&f.ExpectedDocumentID,
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

func (r *RequestRepo) GetFilesByRequest(ctx context.Context, requestID int) ([]models.DocumentDTORead, error) {
	query := `
        SELECT df.id, df.document_request_id, df.expected_document_id, df.file_name, df.file_path, df.mime_type, df.file_size, df.uploaded_at, df.s3_version_id, df.uploaded_by, u.first_name, u.last_name
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

	var files []models.DocumentDTORead
	for rows.Next() {
		var f models.DocumentDTORead
		if err := rows.Scan(
			&f.ID,
			&f.RequestID,
			&f.ExpectedDocumentID,
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

func (r *RequestRepo) AddDocument(ctx context.Context, file models.Document) (int, error) {
	var id int
	query := `
        INSERT INTO document_files (
            document_request_id,
            expected_document_id,
            file_name,
            file_path,
            mime_type,
            file_size,
            s3_version_id,
            uploaded_at,
            uploaded_by
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query,
		file.RequestID,
		file.ExpectedDocumentID,
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

func (r *RequestRepo) SetFileUploaded(ctx context.Context, id int) error {
	query := `
		UPDATE document_requests
		SET last_uploaded_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *RequestRepo) GetRequestsByProfessionalWithExpectedDocs(
	ctx context.Context,
	professionalID int,
	search *string,
) ([]models.RequestDTORead, error) {
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, 
			   dr.is_scheduled, dr.scheduled_for, dr.is_closed, dr.last_uploaded_at, 
			   u.email, u.first_name, u.last_name, dr.title, dr.description, 
			   dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at, dr.template_id,
			   ed.id, ed.document_request_id, ed.title, ed.description, ed.status, ed.rejection_reason, ed.example_file_path, ed.example_mime_type
		FROM document_requests dr
		JOIN users u ON dr.client_id = u.id
		LEFT JOIN expected_documents ed ON ed.document_request_id = dr.id
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

	query += " ORDER BY dr.created_at DESC, ed.id ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requestMap := make(map[int]*models.RequestDTORead)
	requestOrder := make([]int, 0)

	for rows.Next() {
		var req models.RequestDTORead
		var edID *int
		var edRequestID *int
		var edTitle *string
		var edDescription *string
		var edStatus *string
		var edRejectionReason *string
		var edExampleFilePath *string
		var edExampleMimeType *string

		err := rows.Scan(
			&req.ID, &req.ProfessionalID, &req.ClientID, &req.IsRecurring, &req.RecurrenceCron,
			&req.IsScheduled, &req.ScheduledFor, &req.IsClosed, &req.LastUploadedAt,
			&req.ClientEmail, &req.ClientFirstName, &req.ClientLastName,
			&req.Title, &req.Description, &req.DueDate, &req.NextDueAt,
			&req.CreatedAt, &req.UpdatedAt, &req.RequestTemplateID,
			&edID, &edRequestID, &edTitle, &edDescription, &edStatus, &edRejectionReason, &edExampleFilePath, &edExampleMimeType,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := requestMap[req.ID]; !exists {
			req.ExpectedDocuments = make([]models.ExpectedDocument, 0)
			requestMap[req.ID] = &req
			requestOrder = append(requestOrder, req.ID)
		}

		if edID != nil {
			ed := models.ExpectedDocument{
				ID:              *edID,
				RequestID:       *edRequestID,
				Title:           *edTitle,
				Description:     *edDescription,
				Status:          *edStatus,
				RejectionReason: edRejectionReason,
				ExampleFilePath: edExampleFilePath,
				ExampleMimeType: edExampleMimeType,
			}
			requestMap[req.ID].ExpectedDocuments = append(requestMap[req.ID].ExpectedDocuments, ed)
		}
	}

	requests := make([]models.RequestDTORead, 0, len(requestOrder))
	for _, id := range requestOrder {
		requests = append(requests, *requestMap[id])
	}

	return requests, rows.Err()
}

func (r *RequestRepo) GetRequestsByClientWithExpectedDocs(
	ctx context.Context,
	clientID int,
	search *string,
) ([]models.RequestDTORead, error) {
	query := `
		SELECT dr.id, dr.professional_id, dr.client_id, dr.is_recurring, dr.recurrence_cron, 
			   dr.is_scheduled, dr.scheduled_for, dr.is_closed, dr.last_uploaded_at, 
			   u.email, u.first_name, u.last_name, dr.title, dr.description, 
			   dr.due_date, dr.next_due_at, dr.created_at, dr.updated_at, dr.template_id,
			   ed.id, ed.document_request_id, ed.title, ed.description, ed.status, ed.rejection_reason, ed.example_file_path, ed.example_mime_type
		FROM document_requests dr
		JOIN users u ON dr.client_id = u.id
		LEFT JOIN expected_documents ed ON ed.document_request_id = dr.id
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

	query += " ORDER BY dr.created_at DESC, ed.id ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requestMap := make(map[int]*models.RequestDTORead)
	requestOrder := make([]int, 0)

	for rows.Next() {
		var req models.RequestDTORead
		var edID *int
		var edRequestID *int
		var edTitle *string
		var edDescription *string
		var edStatus *string
		var edRejectionReason *string
		var edExampleFilePath *string
		var edExampleMimeType *string

		err := rows.Scan(
			&req.ID, &req.ProfessionalID, &req.ClientID, &req.IsRecurring, &req.RecurrenceCron,
			&req.IsScheduled, &req.ScheduledFor, &req.IsClosed, &req.LastUploadedAt,
			&req.ClientEmail, &req.ClientFirstName, &req.ClientLastName,
			&req.Title, &req.Description, &req.DueDate, &req.NextDueAt,
			&req.CreatedAt, &req.UpdatedAt, &req.RequestTemplateID,
			&edID, &edRequestID, &edTitle, &edDescription, &edStatus, &edRejectionReason, &edExampleFilePath, &edExampleMimeType,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := requestMap[req.ID]; !exists {
			req.ExpectedDocuments = make([]models.ExpectedDocument, 0)
			requestMap[req.ID] = &req
			requestOrder = append(requestOrder, req.ID)
		}

		if edID != nil {
			ed := models.ExpectedDocument{
				ID:              *edID,
				RequestID:       *edRequestID,
				Title:           *edTitle,
				Description:     *edDescription,
				Status:          *edStatus,
				RejectionReason: edRejectionReason,
				ExampleFilePath: edExampleFilePath,
				ExampleMimeType: edExampleMimeType,
			}
			requestMap[req.ID].ExpectedDocuments = append(requestMap[req.ID].ExpectedDocuments, ed)
		}
	}

	requests := make([]models.RequestDTORead, 0, len(requestOrder))
	for _, id := range requestOrder {
		requests = append(requests, *requestMap[id])
	}

	return requests, rows.Err()
}
