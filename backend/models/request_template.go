package models

import "time"

type RequestTemplate struct {
	ID             int       `db:"id" json:"id"`
	Title          string    `db:"title" json:"title"`
	Description    *string   `db:"description,omitempty" json:"description,omitempty"`
	DepartmentID   int       `db:"department_id" json:"department_id"`
	IsRecurring    bool      `db:"is_recurring" json:"is_recurring"`
	CreatedBy      int       `db:"created_by" json:"created_by"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	RecurrenceCron *string   `db:"recurrence_cron" json:"recurrence_cron"`
	IsClosed       bool      `db:"is_closed" json:"is_closed"`
}

type RequestTemplateDTORead struct {
	RequestTemplate
	AuthorFirstName *string `json:"author_first_name"`
	AuthorLastName  *string `json:"author_last_name"`
	DepartmentName  string  `json:"department_name"`
}

type RequestTemplateDTOPatch struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	IsRecurring    *bool   `json:"is_recurring"`
	RecurrenceCron *string `json:"recurrence_cron"`
}
