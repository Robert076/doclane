package events

import "time"

type Event struct {
	Type         string         // "request.claimed", "user.deactivated", "department.created"
	ActorID      int            // who did it
	ResourceID   int            // what was affected
	ResourceType string         // "request", "user", "department"
	Metadata     map[string]any // anything extra — old value, new value, reason etc.
	OccurredAt   time.Time
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
	ResourceTypeRequest = "request"
)
