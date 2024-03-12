package router

import (
	"student_score/controllers"
	"student_score/middleware"

	"github.com/labstack/echo/v4"
)

func New(e *echo.Echo) *echo.Echo {

	//LOGIN
	e.POST("/api/v1/student/login", controllers.StudentLogin)
	e.POST("/api/v1/teacher/login", controllers.TeacherLogin)

	//RULE
	e.GET("/api/v1/rule", controllers.GetRule)

	//CARD
	e.GET("/api/v1/card", controllers.GetCard)

	//STUDENT
	e.GET("/api/v1/student/check-score", controllers.StudentCheckScore, middleware.JWTMiddleware)
	e.POST("/api/v1/student/use-card", controllers.StudentUseCard, middleware.JWTMiddleware)
	e.GET("/api/v1/student/increase-history", controllers.StudentIncreaseHistory, middleware.JWTMiddleware)
	e.GET("/api/v1/student/decrease-history", controllers.StudentDecreaseHistory, middleware.JWTMiddleware)

	//TEACHER
	e.POST("/api/v1/teacher/punish", controllers.Punish, middleware.JWTMiddleware)
	e.GET("/api/v1/teacher/student-list", controllers.StudentList, middleware.JWTMiddleware)

	//TEACHER RULE MANAGEMENT
	e.POST("/api/v1/rule", controllers.AddRule, middleware.JWTMiddleware)
	e.DELETE("/api/v1/rule/:rule_id", controllers.DeleteRule, middleware.JWTMiddleware)

	//TEACHER CARD MANAGEMENT
	e.DELETE("/api/v1/card/:card_id", controllers.DeleteCard, middleware.JWTMiddleware)
	e.POST("/api/v1/create-card", controllers.CreateCard, middleware.JWTMiddleware)

	return e
}
