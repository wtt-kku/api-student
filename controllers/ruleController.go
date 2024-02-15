package controllers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"student_score/config"
	"student_score/models"
	"student_score/utils"

	"github.com/labstack/echo/v4"
)

func GetRule(c echo.Context) (err error) {

	condition_sql := "where r.is_deleted = false"

	queryType := c.QueryParam("type")
	if queryType != "" {
		condition_sql += fmt.Sprintf(` and r.type = '%s'`, strings.ToUpper(queryType))
	}

	rule := []models.Rule{}
	err = config.DbPostgres.Select(&rule, fmt.Sprintf(`select id, type, title, description, point, created_at, updated_at from "rule" r %s order by r.created_at asc`, condition_sql))

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, &utils.Response{
				Result:  false,
				Code:    400,
				Message: "Rule Not found",
			})
		} else if err != nil {
			slog.Error("GET_RULE", "msg", err)
			return c.JSON(http.StatusInternalServerError, &utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Code,
				Message: utils.CommonRespCode["INTERNAL_SERVER_ERROR"].Message,
			})
		}

	}

	return c.JSON(http.StatusOK, &utils.Response{
		Result:  true,
		Code:    utils.CommonRespCode["OK"].Code,
		Message: utils.CommonRespCode["OK"].Message,
		Data:    rule,
	})
}
