package models

import "time"

type Document struct {
	ID                 int       `db:"id" json:"id"`
	RequestID          int       `db:"document_request_id" json:"document_request_id"`
	ExpectedDocumentID int       `db:"expected_document_id" json:"expected_document_id"`
	FileName           string    `db:"file_name" json:"file_name"`
	FilePath           string    `db:"file_path" json:"file_path"`
	MimeType           *string   `db:"mime_type,omitempty" json:"mime_type,omitempty"`
	FileSize           *int64    `db:"file_size,omitempty" json:"file_size,omitempty"`
	UploadedAt         time.Time `db:"uploaded_at" json:"uploaded_at"`
	S3VersionID        *string   `db:"s3_version_id" json:"s3_version_id"`
	UploadedBy         *int      `db:"uploaded_by" json:"uploaded_by"`
}

type DocumentDTORead struct {
	Document
	UploadedByFirstName string `json:"uploaded_by_first_name"`
	UploadedByLastName  string `json:"uploaded_by_last_name"`
}

type DocumentDTOExtended struct {
	Document
	AuthorRole string
}
