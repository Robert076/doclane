package models

type ExpectedDocument struct {
	ID                int    `db:"id" json:"id"`
	DocumentRequestID int    `db:"document_request_id" json:"document_request_id"`
	Title             string `db:"title" json:"title"`
	Description       string `db:"description" json:"description"`
	IsUploaded        bool   `db:"is_uploaded" json:"is_uploaded"`
}
