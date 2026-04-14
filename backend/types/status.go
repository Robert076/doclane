package types

const (
	StatusPending  = "pending"
	StatusUploaded = "uploaded"
	StatusAccepted = "accepted"
	StatusDenied   = "denied"
)

func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusUploaded || status == StatusAccepted || status == StatusDenied
}
