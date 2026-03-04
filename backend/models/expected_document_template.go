package models

type ExpectedDocumentTemplate struct {
	ID                        int     `db:"id" json:"id"`
	DocumentRequestTemplateID int     `db:"document_request_template_id" json:"document_request_template_id"`
	Title                     string  `db:"title" json:"title"`
	Description               string  `db:"description" json:"description"`
	ExampleFilePath           *string `db:"example_file_path" json:"example_file_path"`
	ExampleMimeType           *string `db:"example_mime_type" json:"example_mime_type"`
}
