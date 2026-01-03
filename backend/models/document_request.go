package models

import "time"

type DocumentRequest struct {
	ID             int        `db:"id" json:"id"`
	ProfessionalID int        `db:"professional_id" json:"professional_id"`
	ClientID       int        `db:"client_id" json:"client_id"`
	Title          string     `db:"title" json:"title"`
	Description    *string    `db:"description,omitempty" json:"description,omitempty"`
	DueDate        *time.Time `db:"due_date,omitempty" json:"due_date,omitempty"`
	Status         string     `db:"status" json:"status"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
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
}
