package requests

type UpdateExpectedDocumentStatusRequest struct {
	Status          string  `json:"status" validate:"required"` // "pending", "uploaded", "accepted", "rejected"
	RejectionReason *string `json:"rejection_reason"`
}
