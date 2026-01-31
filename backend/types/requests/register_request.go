package requests

type RegisterRequest struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Password       string `json:"password"`
	Role           string `json:"role"`
	ProfessionalID *int   `json:"professional_id"`
}
