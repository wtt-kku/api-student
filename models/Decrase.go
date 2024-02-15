package models

import "time"

type DecreaseRecordJoinRule struct {
	Id          string    `json:"id" db:"id"`
	Title       string    `json:"title"  db:"title"`
	Description string    `json:"description" db:"description"`
	Point       int       `json:"point" db:"point"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
