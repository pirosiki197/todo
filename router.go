package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/pirosiki197/todo/handler"
)

func setUpRoutes(e *echo.Echo, db *sqlx.DB) {
	withAuth := e.Group("")
	withAuth.Use(userAuthMidlleware)
	ah := handler.NewAuthHandler(db)
	th := handler.NewTodoHandler(db)

	e.POST("/signup", ah.SignUp)
	e.POST("/login", ah.Login)

	withAuth.GET("/check", ah.Check)
	withAuth.POST("/todo", th.CreateNewTodo)
	// TODO:分かりづらいから変えるべき
	withAuth.PATCH("/todo/:id", th.FinishTodo)
}
