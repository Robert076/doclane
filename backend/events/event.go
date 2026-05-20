package events

import "time"

type Event struct {
	ID             int            `json:"id"`
	Type           string         `json:"event_type"`
	ActorID        int            `json:"actor_id"`
	ResourceID     int            `json:"resource_id"`
	ResourceType   string         `json:"resource_type"`
	Metadata       map[string]any `json:"metadata"`
	OccurredAt     time.Time      `json:"occurred_at"`
	ActorFirstName *string        `json:"actor_first_name"`
	ActorLastName  *string        `json:"actor_last_name"`
}

const (
	EventRequestCreated   = "request.created"
	EventRequestClaimed   = "request.claimed"
	EventRequestUnclaimed = "request.unclaimed"
	EventRequestClosed    = "request.closed"
	EventRequestReopened  = "request.reopened"
	EventRequestCancelled = "request.cancelled"
	EventRequestUpdated   = "request.updated"

	EventDocumentUploaded = "document.uploaded"
	EventDocumentApproved = "document.approved"
	EventDocumentRejected = "document.rejected"

	EventUserDeactivated = "user.deactivated"
	EventUserNotified    = "user.notified"

	EventDepartmentCreated = "department.created"
)

const (
	ResourceTypeRequest    = "request"
	ResourceTypeDocument   = "document"
	ResourceTypeDepartment = "department"
	ResourceTypeUser       = "user"
)
