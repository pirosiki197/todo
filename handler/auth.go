package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	salt = "a1e93f67d82b4c56"
)

type AuthHandler struct {
	db *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

type LoginRequestBody struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type User struct {
	UserName       string `db:"username"`
	HashedPassword string `db:"password"`
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.UserName == "" || req.Password == "" {
		return c.String(http.StatusBadRequest, "username or password is empty")
	}

	// UserNameが同じuserがいないか確認する
	var count int
	err := h.db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username = ?", req.UserName)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if count > 0 {
		return c.String(http.StatusConflict, "Username is already used")
	}

	// パスワードのハッシュ化
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// 作成したuserをデータベースに格納
	_, err = h.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.UserName, hashedPass)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	// データベースからusernameをkeyにuserを持ってくる
	var user User
	err := h.db.Get(&user, "SELECT * FROM users WHERE username = ?", req.UserName)
	if errors.Is(err, sql.ErrNoRows) {
		return c.NoContent(http.StatusUnauthorized)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// パスワードが合っているか確認
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password+salt))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return c.NoContent(http.StatusUnauthorized)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// セッションに登録
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	fmt.Printf("%+v\n", sess)
	sess.Values["userName"] = req.UserName
	sess.Save(c.Request(), c.Response())
	fmt.Printf("%+v\n", sess)

	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) Check(c echo.Context) error {
	return c.String(200, "hello")
}
