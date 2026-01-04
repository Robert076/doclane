package requests

import "time"

type AddDocumentRequestRequest struct {
	ClientID    int        `json:"client_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}
