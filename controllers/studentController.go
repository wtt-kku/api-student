package controllers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"student_score/config"
	"student_score/dto"
	"student_score/middleware"
	"student_score/models"
	"student_score/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func StudentLogin(c echo.Context) (err error) {
	//BIND BODY
	body := new(dto.StudentLoginDTO)
	if err = c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(body); err != nil {
		return err
	}

	student := models.StudentInfo{}

	err = config.DbPostgres.Get(&student, `select id, firstname, lastname, gender, class, student_no from student s 
	where s.student_no = $1 and s.student_password = $2 and s.is_deleted = false
	limit 1`, body.StudentNo, body.Password)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    3000,
				Message: "Student No. or Password Invalid",
			})
		} else if err != nil {
			slog.Error("STUDENT_LOGIN", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	token, _ := middleware.GenerateJWT(student.Id)

	res := dto.StudentLoginResDTO{
		Token:       token,
		StudentInfo: student,
	}

	return c.JSON(http.StatusOK, &utils.Response{
		Result:  true,
		Code:    2000,
		Message: "OK",
		Data:    res,
	})
}

func StudentCheckScore(c echo.Context) (err error) {
	UserId := c.Get("userId").(string)
	// UserId := "f3ac15f8-d17e-4489-afed-6e8c60293585"

	student := models.StudentInfo{}

	err = config.DbPostgres.Get(&student, `select id, firstname, lastname, gender, class, student_no from student s 
	where s.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Student Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_CHECKSCORE_#1", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	var studentScore int

	err = config.DbPostgres.Get(&studentScore, `select s.score  from score s 
	where s.student_no = $1 limit 1`, student.Student_no)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Score Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_CHECKSCORE_STEP#2", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	var pass = false
	var res dto.StudentCheckScoreResDTO
	res.Score = studentScore

	if studentScore >= 50 {
		pass = true
	}
	res.Pass = pass

	return c.JSON(http.StatusOK, utils.Response{
		Result:  true,
		Code:    utils.CommonRespCode["OK"].Code,
		Message: utils.CommonRespCode["OK"].Message,
		Data:    res,
	})
}

func StudentUseCard(c echo.Context) (err error) {
	UserId := c.Get("userId").(string)
	// UserId := "f3ac15f8-d17e-4489-afed-6e8c60293585"

	body := new(dto.StudentUseCardDTO)
	if err = c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(body); err != nil {
		return err
	}

	//GET STUDENT DATA
	student := models.StudentInfo{}

	err = config.DbPostgres.Get(&student, `select id, firstname, lastname, gender, class, student_no from student s 
	where s.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Student Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_USECARD_#1", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	//GET STUDENT SCORE
	var studentScore int
	err = config.DbPostgres.Get(&studentScore, `select s.score  from score s 
	where s.student_no = $1 limit 1`, student.Student_no)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    9000,
				Message: "Score Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_USECARD_#2", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	//CHECK CARD & RULE DATA

	rule := models.StudentGetRuleInfoByCard{}
	err = config.DbPostgres.Get(&rule, `select  r.id as "rule_id" , c.id  as "card_id" , r.point , c.status as "card_status" from card c 
	left join "rule" r 
	on r.id  = c.rule_id 
	where c.card_code  = $1`, body.Code)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusOK, &utils.Response{
				Result:  false,
				Code:    3200,
				Message: "Card Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_USECARD_#3", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	if rule.CardStatus == 1 {
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    3300,
			Message: "This card has already been used",
		})
	}

	if rule.CardStatus == 2 {
		return c.JSON(http.StatusOK, &utils.Response{
			Result:  false,
			Code:    3400,
			Message: "This card has been cancelled",
		})
	}

	tx, err := config.DbPostgres.Beginx()
	if err != nil {
		slog.Error("DB TX BEGIN", "msg", err)
		return c.JSON(http.StatusInternalServerError, &utils.Response{
			Result:  false,
			Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
			Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
		})
	}

	//CHANGE STUDENT SCORE (UPDATE)
	newScore := studentScore + rule.Point
	query1 := `UPDATE score SET score=:score, updated_at=:update WHERE student_no=:student_no`
	params1 := map[string]interface{}{
		"score":      newScore,
		"update":     time.Now(),
		"student_no": student.Student_no,
	}
	_, err = tx.NamedExec(query1, params1)
	if err != nil {
		tx.Rollback()
	}

	//INSERT RECORD INCREMENT (INSERT)

	query2 := `INSERT INTO record_increase (id, student_no, rule_id, created_at, updated_at, point, card_id) VALUES(uuid_generate_v4(), :student_no, :rule_id, :now, :now, :point, :card_id)`
	params2 := map[string]interface{}{
		"student_no": student.Student_no,
		"rule_id":    rule.RuleId,
		"now":        time.Now(),
		"point":      rule.Point,
		"card_id":    rule.CardId,
	}
	_, err = tx.NamedExec(query2, params2)
	if err != nil {
		tx.Rollback()
	}

	//STAMP CARD TO USED (UPDATE)
	query3 := `UPDATE card SET status=:status , updated_at=:now  WHERE id=:card_id`
	params3 := map[string]interface{}{
		"status":  1,
		"card_id": rule.CardId,
		"now":     time.Now(),
	}

	_, err = tx.NamedExec(query3, params3)
	if err != nil {
		tx.Rollback()
	}

	//COMMIT
	err = tx.Commit()
	if err != nil {
		slog.Error("DB TX COMMIT", "msg", err)
		return c.JSON(http.StatusInternalServerError, &utils.Response{
			Result:  false,
			Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
			Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
		})
	}

	return c.JSON(http.StatusOK, utils.Response{
		Result:  true,
		Code:    utils.CommonRespCode["OK"].Code,
		Message: utils.CommonRespCode["OK"].Message,
		Data: map[string]interface{}{
			"score_increase": rule.Point,
		},
	})
}

func StudentIncreaseHistory(c echo.Context) (err error) {
	UserId := c.Get("userId").(string)
	// UserId := "f3ac15f8-d17e-4489-afed-6e8c60293585"

	//GET STUDENT DATA
	student := models.StudentInfo{}

	err = config.DbPostgres.Get(&student, `select id, firstname, lastname, gender, class, student_no from student s 
	where s.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, &utils.Response{
				Result:  false,
				Code:    3000,
				Message: "Student Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_INCREASE_HISTORY_#1", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	//GET HISTORY
	//GET STUDENT DATA
	history := []models.IncreaseRecordJoinRule{}

	err = config.DbPostgres.Select(&history, `select ri.id , r.title  , r.description , ri.point , ri.created_at  from record_increase ri 
	left join "rule" r 
	on r.id  = ri.rule_id 
	where student_no = $1
	order by ri.created_at asc`, student.Student_no)

	if err != nil {

		if err != nil {
			slog.Error("STUDENT_INCREASE_HISTORY_#2", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	return c.JSON(http.StatusOK, utils.Response{
		Result:  true,
		Code:    utils.CommonRespCode["OK"].Code,
		Message: utils.CommonRespCode["OK"].Message,
		Data:    history,
	})
}

func StudentDecreaseHistory(c echo.Context) (err error) {
	UserId := c.Get("userId").(string)
	// UserId := "f3ac15f8-d17e-4489-afed-6e8c60293585"

	//GET STUDENT DATA
	student := models.StudentInfo{}

	err = config.DbPostgres.Get(&student, `select id, firstname, lastname, gender, class, student_no from student s 
	where s.id = $1 limit 1`, UserId)

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, &utils.Response{
				Result:  false,
				Code:    3000,
				Message: "Student Not found",
			})
		} else if err != nil {
			slog.Error("STUDENT_DECREASE_HISTORY_#1", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	//GET HISTORY
	//GET STUDENT DATA
	history := []models.DecreaseRecordJoinRule{}

	err = config.DbPostgres.Select(&history, `select rd.id , r.title  , r.description , rd.point , rd.created_at  from record_decrease rd
	left join "rule" r 
	on r.id  = rd.rule_id 
	where student_no = $1
	order by rd.created_at asc`, student.Student_no)

	if err != nil {

		if err != nil {
			slog.Error("STUDENT_DECREASE_HISTORY_#2", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	return c.JSON(http.StatusOK, utils.Response{
		Result:  true,
		Code:    utils.CommonRespCode["OK"].Code,
		Message: utils.CommonRespCode["OK"].Message,
		Data:    history,
	})
}
