package models

import "time"

type StudentEntity struct {
	Id               string    `db:"id" json:"id"`
	Firstname        string    `db:"firstname" json:"firstname"`
	Lastname         string    `db:"lastname" json:"lastname"`
	Gender           string    `db:"gender" json:"gender"`
	Class            string    `db:"class" json:"class"`
	Student_no       string    `db:"student_no" json:"student_no"`
	Student_password string    `db:"student_password" json:"student_password"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted        bool      `db:"is_deleted" json:"is_deleted"`
}

type StudentInfo struct {
	Id         string `db:"id" json:"id"`
	Firstname  string `db:"firstname" json:"firstname"`
	Lastname   string `db:"lastname" json:"lastname"`
	Gender     string `db:"gender" json:"gender"`
	Class      string `db:"class" json:"class"`
	Student_no string `db:"student_no" json:"student_no"`
}

type StudentGetRuleInfoByCard struct {
	RuleId     string `db:"rule_id" json:"rule_id"`
	CardId     string `db:"card_id" json:"card_id"`
	Point      int    `db:"point" json:"point"`
	CardStatus int    `db:"card_status" json:"card_status"`
}
