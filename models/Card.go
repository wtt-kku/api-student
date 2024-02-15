package models

import "time"

type Card struct {
	Id        string    `db:"id" json:"id"`
	CardCode  string    `db:"card_code" json:"card_code"`
	RuleId    string    `db:"rule_id" json:"rule_id"`
	Status    int       `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
