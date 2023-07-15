package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "naro",
		Loc:    jst,
	}
	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(store))

	setUpRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}

func userAuthMidlleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if sess.Values["userName"] == nil {
			fmt.Printf("%+v", sess)
			fmt.Println(c.Request().Header.Get("sessions"))
			return c.String(http.StatusUnauthorized, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}
