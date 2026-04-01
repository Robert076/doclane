package models

import "time"

type User struct {
	ID           int        `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	FirstName    string     `db:"first_name" json:"first_name"`
	LastName     string     `db:"last_name" json:"last_name"`
	PasswordHash string     `db:"password_hash" json:"-"`
	Role         string     `db:"role" json:"role"` // 'admin' | 'member'
	DepartmentID *int       `db:"department_id" json:"department_id"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	LastNotified *time.Time `db:"last_notified" json:"last_notified"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}
