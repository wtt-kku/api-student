package dto

import "student_score/models"

type TeacherLoginDTO struct {
	TeacherNo string `json:"teacher_no" validate:"required" `
	Password  string `json:"password" validate:"required" `
}

type TeacherLoginResDTO struct {
	Token       string             `json:"token"`
	TeacherInfo models.TeacherInfo `json:"teacher_data"`
}

type PunishDTO struct {
	RuleId      string   `json:"rule_id"`
	StudentList []string `json:"student_list"`
}

type PunishResDTO struct {
	CountAll     int      `json:"count_all"`
	CountSuccess int      `json:"count_success"`
	CountFail    int      `json:"count_fail"`
	ListSuccess  []string `json:"list_success"`
	ListFail     []string `json:"list_fail"`
}

type AddRuleDTO struct {
	RuleName  string `json:"rule_name"`
	RuleDesc  string `json:"rule_desc"`
	RuleType  int    `json:"rule_type"`
	RuleScore int    `json:"rule_point"`
}

type CreateCardDTO struct {
	RuleId     string `json:"rule_id"`
	CardAmount int    `json:"card_amount"`
}
