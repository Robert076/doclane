package models

import "time"

type DocumentRequestBase struct {
	ClientID       int        `db:"client_id" json:"client_id"`
	Title          string     `db:"title" json:"title"`
	Description    *string    `db:"description,omitempty" json:"description,omitempty"`
	IsRecurring    bool       `db:"is_recurring" json:"is_recurring"`
	RecurrenceCron *string    `db:"recurrence_cron" json:"recurrence_cron"`
	IsScheduled    bool       `db:"is_scheduled" json:"is_scheduled"`
	ScheduledFor   *string    `db:"scheduled_for" json:"scheduled_for"`
	IsClosed       bool       `db:"is_closed" json:"is_closed"`
	LastUploadedAt *time.Time `db:"last_uploaded_at" json:"last_uploaded_at"`
	NextDueAt      *time.Time `db:"next_due_at" json:"next_due_at"`
	DueDate        *time.Time `db:"due_date,omitempty" json:"due_date,omitempty"`
}

type DocumentRequest struct {
	DocumentRequestBase
	ID             int       `db:"id" json:"id"`
	ProfessionalID int       `db:"professional_id" json:"professional_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type DocumentRequestDTORead struct {
	DocumentRequest
	ClientEmail     string `json:"client_email"`
	ClientFirstName string `json:"client_first_name"`
	ClientLastName  string `json:"client_last_name"`
	Status          string `json:"status"`
}

type DocumentRequestDTOCreate struct {
	DocumentRequestBase
}

type DocumentRequestDTOPatch struct {
	Title string `json:"title"`
}

type DocumentFile struct {
	ID                int       `db:"id" json:"id"`
	DocumentRequestID int       `db:"document_request_id" json:"document_request_id"`
	FileName          string    `db:"file_name" json:"file_name"`
	FilePath          string    `db:"file_path" json:"file_path"`
	MimeType          *string   `db:"mime_type,omitempty" json:"mime_type,omitempty"`
	FileSize          *int64    `db:"file_size,omitempty" json:"file_size,omitempty"`
	UploadedAt        time.Time `db:"uploaded_at" json:"uploaded_at"`
	S3VersionID       *string   `db:"s3_version_id" json:"s3_version_id"`
	UploadedBy        *int      `db:"uploaded_by" json:"uploaded_by"`
}

type DocumentFileDTORead struct {
	DocumentFile
	UploadedByFirstName string `json:"uploaded_by_first_name"`
	UploadedByLastName  string `json:"uploaded_by_last_name"`
}

type DocumentFileDTOExtended struct {
	DocumentFile
	AuthorRole string
}
