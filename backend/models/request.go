package models

import "time"

type RequestBase struct {
	Title          string     `db:"title" json:"title"`
	Description    *string    `db:"description,omitempty" json:"description,omitempty"`
	Assignee       int        `db:"assignee" json:"assignee"`
	DepartmentID   int        `db:"department_id" json:"department_id"`
	IsRecurring    bool       `db:"is_recurring" json:"is_recurring"`
	RecurrenceCron *string    `db:"recurrence_cron" json:"recurrence_cron"`
	IsScheduled    bool       `db:"is_scheduled" json:"is_scheduled"`
	ScheduledFor   *time.Time `db:"scheduled_for" json:"scheduled_for"`
	IsClosed       bool       `db:"is_closed" json:"is_closed"`
	LastUploadedAt *time.Time `db:"last_uploaded_at" json:"last_uploaded_at"`
	NextDueAt      *time.Time `db:"next_due_at" json:"next_due_at"`
	DueDate        *time.Time `db:"due_date,omitempty" json:"due_date,omitempty"`
}
type Request struct {
	RequestBase
	ID                int       `db:"id" json:"id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	RequestTemplateID *int      `db:"template_id" json:"template_id"`
}

type RequestDTORead struct {
	Request
	AssigneeEmail     string             `json:"assignee_email"`
	AssigneeFirstName string             `json:"assignee_first_name"`
	AssigneeLastName  string             `json:"assignee_last_name"`
	Status            string             `json:"status"`
	ExpectedDocuments []ExpectedDocument `json:"expected_documents"`
}

type RequestDTOPatch struct {
	Title string `json:"title"`
}

type RequestDTOCreate struct {
	TemplateID   int        `json:"template_id"`
	IsScheduled  bool       `json:"is_scheduled"`
	ScheduledFor *time.Time `json:"scheduled_for"`
	DueDate      *time.Time `json:"due_date"`
}
