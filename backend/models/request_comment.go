package models

import "time"

type RequestComment struct {
	ID        int       `db:"id" json:"id"`
	RequestID int       `db:"request_id" json:"request_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Comment   string    `db:"comment" json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type RequestCommentDTO struct {
	RequestComment
	UserFirstName string `json:"user_first_name"`
	UserLastName  string `json:"user_last_name"`
}
