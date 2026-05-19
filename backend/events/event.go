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
