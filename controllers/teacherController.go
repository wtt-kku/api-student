package controllers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"student_score/config"
	"student_score/dto"
	"student_score/middleware"
	"student_score/models"
	"student_score/utils"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TeacherLogin(c echo.Context) (err error) {
	//BIND BODY
	body := new(dto.TeacherLoginDTO)
	if err = c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(body); err != nil {
		return err
	}

	teacher := models.TeacherInfo{}

	err = config.DbPostgres.Get(&teacher, `select id, firstname, lastname, gender, class, teacher_no from teacher t 
	where t.teacher_no = $1 and t.teacher_password = $2 and t.is_deleted = false
	limit 1`, body.TeacherNo, body.Password)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    3000,
				Message: "Teacher No. or Password Invalid",
			})
		} else if err != nil {
			slog.Error("TEACHER_LOGIN", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    500,
				Message: "Server Error",
			})
		}

	}

	token, _ := middleware.GenerateJWT(teacher.Id)

	res := dto.TeacherLoginResDTO{
		Token:       token,
		TeacherInfo: teacher,
	}

	return c.JSON(http.StatusOK, &utils.Response{
		Result:  true,
		Code:    2000,
		Message: "OK",
		Data:    res,
	})
}

func Punish(c echo.Context) (err error) {
	UserId := c.Get("userId").(string)

	body := new(dto.PunishDTO)
	if err = c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(body); err != nil {
		return err
	}

	teacher := models.TeacherInfo{}
	err = config.DbPostgres.Get(&teacher, `select id, firstname, lastname, gender, class, teacher_no  from teacher t  where t.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Teacher Not found",
			})
		} else if err != nil {
			slog.Error("TEACHER_PUNISH_CHECK_TID", "msg", err)
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	//CHECK UUID FORMAT
	_, err = uuid.Parse(body.RuleId)
	if err != nil {
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    9000,
			Message: "Rule Not found",
		})

	}

	//GET RULE
	rule := models.Rule{}

	err = config.DbPostgres.Get(&rule, `select id, type, title, description, point, created_at, updated_at, is_deleted from "rule" r  where r.id  = $1 and r.is_deleted = false  limit 1`, body.RuleId)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Rule Not found",
			})
		} else if err != nil {
			slog.Error("TEACHER_PUNISH_GET_RULE", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	//CREATE VARIABLE STORE SUCCES & FAIL
	CountSuccess := 0
	CountFail := 0
	ListSuccess := []string{}
	ListFail := []string{}

	for _, v := range body.StudentList {

		//GET STUDENT DATA
		student := models.StudentInfo{}
		err = config.DbPostgres.Get(&student, `select id, firstname, lastname, gender, class, student_no from student s 
	where s.student_no = $1 limit 1`, v)

		if err != nil {
			if err != nil {
				slog.Error("TEACHER_PUNISH_GET_STUDENT_DATA", "msg", err, "Student No", v)
				CountFail++
				ListFail = append(ListFail, v)
				continue
			}

		}

		//GET STUDENT SCORE
		var studentScore int
		err = config.DbPostgres.Get(&studentScore, `select s.score  from score s 
	where s.student_no = $1 limit 1`, student.Student_no)

		if err != nil {
			if err != nil {
				slog.Error("TEACHER_PUNISH_GET_STUDENT_SCORE", "msg", err, "Student No", v)
				CountFail++
				ListFail = append(ListFail, v)
				continue
			}

		}

		tx, err := config.DbPostgres.Beginx()
		if err != nil {
			slog.Error("DB TX BEGIN", "msg", err, "Student No", v)
			CountFail++
			ListFail = append(ListFail, v)
			continue
		}

		//CHANGE STUDENT SCORE (UPDATE)
		newScore := studentScore - rule.Point
		query1 := `UPDATE score SET score=:score, updated_at=:update WHERE student_no=:student_no`
		params1 := map[string]interface{}{
			"score":      newScore,
			"update":     time.Now(),
			"student_no": student.Student_no,
		}

		_, err = tx.NamedExec(query1, params1)
		if err != nil {
			tx.Rollback()
			slog.Error("DB TX ROLLBACK", "msg", err, "Student No", v)
			CountFail++
			ListFail = append(ListFail, v)
			continue
		}

		//INSERT RECORD DECREMENT (INSERT)
		query2 := `INSERT INTO record_decrease (id, student_no, rule_id, created_at, updated_at, point) VALUES(uuid_generate_v4(), :student_no, :rule_id, :now, :now, :point)`
		params2 := map[string]interface{}{
			"student_no": student.Student_no,
			"rule_id":    rule.Id,
			"now":        time.Now(),
			"point":      rule.Point,
		}
		_, err = tx.NamedExec(query2, params2)
		if err != nil {
			tx.Rollback()
			slog.Error("DB TX ROLLBACK", "msg", err, "Student No", v)
			CountFail++
			ListFail = append(ListFail, v)
			continue
		}

		//COMMIT

		err = tx.Commit()
		if err != nil {
			slog.Error("DB TX COMMIT", "msg", err)
			CountFail++
			ListFail = append(ListFail, v)
			continue
		}
		CountSuccess++
		ListSuccess = append(ListSuccess, v)
	}

	return c.JSON(http.StatusOK, utils.Response{
		Result:  true,
		Code:    utils.CommonRespCode["OK"].Code,
		Message: utils.CommonRespCode["OK"].Message,
		Data: dto.PunishResDTO{
			CountAll:     len(body.StudentList),
			CountSuccess: CountSuccess,
			CountFail:    CountFail,
			ListSuccess:  ListSuccess,
			ListFail:     ListFail,
		},
	})

}

func AddRule(c echo.Context) (err error) {

	UserId := c.Get("userId").(string)
	// UserId := "d08d1804-c43b-4406-8c4e-e8a8d1b90b0f"

	body := new(dto.AddRuleDTO)
	if err = c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(body); err != nil {
		return err
	}

	teacher := models.TeacherInfo{}
	err = config.DbPostgres.Get(&teacher, `select id, firstname, lastname, gender, class, teacher_no  from teacher t  where t.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Teacher Not found",
			})
		} else if err != nil {
			slog.Error("TEACHER_PUNISH_CHECK_TID", "msg", err)
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	var RuleType string

	if body.RuleType != 1 && body.RuleType != 2 {
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    3001,
			Message: "Rule type invalid",
		})
	}

	if body.RuleType == 1 {
		RuleType = "INCREASE"
	} else if body.RuleType == 2 {
		RuleType = "DECREASE"
	}

	result, err := config.DbPostgres.NamedExec(`INSERT INTO public."rule" (id, "type", title, description, point, created_at, updated_at, is_deleted) VALUES(uuid_generate_v4(), :rule_type, :rule_name, :rule_desc, :rule_point, :now, :now, false)`,
		map[string]interface{}{
			"rule_type":  RuleType,
			"rule_name":  body.RuleName,
			"rule_desc":  body.RuleDesc,
			"rule_point": body.RuleScore,
			"now":        time.Now(),
		})

	if err != nil {
		slog.Error("TEACHER_ADD_RULE", "msg", err)
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
			Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
		})

	}

	row, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return
	}

	if row != 1 {
		slog.Error("TEACHER_ADD_RULE", "msg", err)
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
			Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
		})
	}

	return c.JSON(http.StatusOK, &utils.Response{
		Result:  true,
		Code:    2000,
		Message: "Success",
	})

}

func DeleteRule(c echo.Context) (err error) {

	UserId := c.Get("userId").(string)
	// UserId := "d08d1804-c43b-4406-8c4e-e8a8d1b90b0f"

	RuleId := c.Param("rule_id")

	teacher := models.TeacherInfo{}
	err = config.DbPostgres.Get(&teacher, `select id, firstname, lastname, gender, class, teacher_no  from teacher t  where t.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Teacher Not found",
			})
		} else if err != nil {
			slog.Error("TEACHER_PUNISH_CHECK_TID", "msg", err)
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	result, err := config.DbPostgres.NamedExec(`UPDATE public."rule" SET updated_at=:now, is_deleted=true WHERE id=:rule_id;
	`,
		map[string]interface{}{
			"rule_id": RuleId,
			"now":     time.Now(),
		})

	if err != nil {
		slog.Error("TEACHER_DELETE_RULE", "msg", err)
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
			Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
		})
	}

	row, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return
	}

	if row != 1 {
		slog.Error("TEACHER_ADD_RULE", "msg", err)
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
			Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
		})
	}

	return c.JSON(http.StatusOK, &utils.Response{
		Result:  true,
		Code:    2000,
		Message: "Success",
	})

}
