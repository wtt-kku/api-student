package models

import "time"

type TeacherEntity struct {
	Id              string    `db:"id" json:"id"`
	Firstname       string    `db:"firstname" json:"firstname"`
	Lastname        string    `db:"lastname" json:"lastname"`
	Gender          string    `db:"gender" json:"gender"`
	Class           string    `db:"class" json:"class"`
	TeacherNo       string    `db:"teacher_no" json:"teacher_no"`
	TeacherPassword string    `db:"teacher_password" json:"teacher_password"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted       bool      `db:"is_deleted" json:"is_deleted"`
}

type TeacherInfo struct {
	Id        string `db:"id" json:"id"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Gender    string `db:"gender" json:"gender"`
	Class     string `db:"class" json:"class"`
	TeacherNo string `db:"teacher_no" json:"teacher_no"`
}
