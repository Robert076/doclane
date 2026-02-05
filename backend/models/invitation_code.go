package models

import "time"

type InvitationCode struct {
	ID             int        `json:"id"`
	Code           string     `json:"code"`
	ProfessionalID int        `json:"professional_id"`
	UsedByUserID   *int       `json:"used_by_user_id,omitempty"`
	UsedAt         *time.Time `json:"used_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type InvitationCodeCreateDTO struct {
	ExpiresInDays int `json:"expires_in_days"`
}

type InvitationCodeValidateDTO struct {
	Code string `json:"code"`
}
