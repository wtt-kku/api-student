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

func GetCard(c echo.Context) (err error) {

	allow_status := []string{"0", "1", "2"}

	condition_sql := ""

	queryType := c.QueryParam("status")

	if queryType != "" && utils.ContainsString(allow_status, queryType) {
		condition_sql = fmt.Sprintf(`where c.status = %s`, strings.ToUpper(queryType))
	}

	card := []models.Card{}
	err = config.DbPostgres.Select(&card, fmt.Sprintf(`select id, card_code, rule_id, status, created_at, updated_at from "card" c %s order by c.created_at asc`, condition_sql))

	if err != nil {

		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, &utils.Response{
				Result:  false,
				Code:    400,
				Message: "Card Not found",
			})
		} else if err != nil {
			slog.Error("GET_CARD", "msg", err)
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
		Data:    card,
	})
}
