package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"student_score/config"
	"student_score/router"
	"student_score/utils"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	var errs []string
	if err := cv.validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		fmt.Println(err.Error())
		for _, validationError := range validationErrors {
			errs = append(errs, (validationError.StructField()))
		}
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, utils.Response{
			Result:  false,
			Code:    400,
			Message: "required: " + strings.Join(errs, " , "),
		})
	}
	return nil
}

func main() {
	utils.CheckEnvReady()
	config.PostgresConn()

	e := echo.New()
	e = router.New(e)

	e.GET("/healthcheck", healthcheck)

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.CORS())

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9000"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
