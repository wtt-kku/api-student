package dto

import "student_score/models"

type StudentLoginDTO struct {
	StudentNo string `json:"student_no" validate:"required" `
	Password  string `json:"password" validate:"required" `
}

type StudentLoginResDTO struct {
	Token       string             `json:"token"`
	StudentInfo models.StudentInfo `json:"student_data"`
}

type StudentCheckScoreResDTO struct {
	Score int  `json:"score"`
	Pass  bool `json:"pass"`
}

type StudentUseCardDTO struct {
	Code string `json:"code" validate:"required" `
}
