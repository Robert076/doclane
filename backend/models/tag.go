package models

import "time"

type Tag struct {
	ID        int       `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Color     string    `db:"color"      json:"color"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type TagDTOCreate struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TagDTOUpdate struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type SetTemplateTagsDTO struct {
	TagIDs []int `json:"tag_ids"`
}
