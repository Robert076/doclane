package types

const (
	StatusPending  = "pending"
	StatusUploaded = "uploaded"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusUploaded || status == StatusAccepted || status == StatusRejected
}
