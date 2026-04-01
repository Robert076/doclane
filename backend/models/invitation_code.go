package models

import "time"

type InvitationCode struct {
	ID        int        `json:"id"`
	Code      string     `json:"code"`
	CreatedBy int        `json:"created_by"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type InvitationCodeCreateDTO struct {
	ExpiresInDays int `json:"expires_in_days"`
}

type InvitationCodeValidateDTO struct {
	Code string `json:"code"`
}
