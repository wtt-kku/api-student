package models

import "time"

type Rule struct {
	Id          string    `db:"id" json:"id"`
	Type        string    `db:"type" json:"type"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Point       int       `db:"point" json:"point"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted"`
}
