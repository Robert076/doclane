package types

type RequestEventType string

const (
	EventRequestClaimed   RequestEventType = "claimed"
	EventRequestUnclaimed RequestEventType = "unclaimed"
	EventRequestClosed    RequestEventType = "closed"
	EventRequestReopened  RequestEventType = "reopened"
	EventRequestCancelled RequestEventType = "cancelled"
)

type RequestEvent struct {
	Type       RequestEventType
	RequestID  int
	ActorID    int // who triggered it (claims.UserID)
	AssigneeID int // who owns the request (to notify them)
}
