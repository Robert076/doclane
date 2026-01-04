package requests

import "time"

type AddDocumentRequestRequest struct {
	ProfessionalID int        `json:"professional_id"`
	ClientID       int        `json:"client_id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
}
