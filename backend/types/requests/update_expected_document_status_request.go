package requests

type UpdateExpectedDocumentStatusRequest struct {
	Status          string  `json:"status" validate:"required"` // "pending", "uploaded", "approved", "rejected"
	RejectionReason *string `json:"rejection_reason"`
}
