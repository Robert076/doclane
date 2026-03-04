package models

import "io"

type ExpectedDocument struct {
	ID                int     `db:"id" json:"id"`
	DocumentRequestID int     `db:"document_request_id" json:"document_request_id"`
	Title             string  `db:"title" json:"title"`
	Description       string  `db:"description" json:"description"`
	Status            string  `db:"status" json:"status"`
	RejectionReason   *string `db:"rejection_reason" json:"rejection_reason"`
	ExampleFilePath   *string `db:"example_file_path" json:"example_file_path"`
	ExampleMimeType   *string `db:"example_mime_type" json:"example_mime_type"`
}

type ExpectedDocumentInput struct {
	Title           string
	Description     string
	ExampleFile     io.Reader
	ExampleFileName string
	ExampleMimeType string
	ExampleFileSize int64
}
