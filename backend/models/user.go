package models

import "time"

type User struct {
	ID             string    `db:"id" json:"id"`
	Email          string    `db:"email" json:"email"`
	PasswordHash   string    `db:"password_hash" json:"-"`
	Role           string    `db:"role" json:"role"`
	ProfessionalID *string   `db:"professional_id" json:"professional_id,omitempty"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
