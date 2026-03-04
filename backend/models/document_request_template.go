package models

import "time"

type DocumentRequestTemplate struct {
	ID             int       `db:"id" json:"id"`
	Title          string    `db:"title" json:"title"`
	Description    *string   `db:"description,omitempty" json:"description,omitempty"`
	IsRecurring    bool      `db:"is_recurring" json:"is_recurring"`
	CreatedBy      int       `db:"created_by" json:"created_by"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	RecurrenceCron *string   `db:"recurrence_cron" json:"recurrence_cron"`
	IsClosed       bool      `db:"is_closed" json:"is_closed"`
}

type DocumentRequestTemplateDTORead struct {
	DocumentRequestTemplate
	AuthorName *string
}
