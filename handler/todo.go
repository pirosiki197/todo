package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	db *sqlx.DB
}

func NewTodoHandler(db *sqlx.DB) *TodoHandler {
	return &TodoHandler{
		db: db,
	}
}

type Todo struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

func (h *TodoHandler) CreateNewTodo(c echo.Context) error {
	var todoRequest Todo
	c.Bind(&todoRequest)

	createrName := c.Get("userName").(string)
	_, err := h.db.Exec("INSERT INTO todos (name, creater, priority) VALUES (?, ?, ?)", todoRequest.Name, createrName, todoRequest.Priority)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *TodoHandler) FinishTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// todoの作成者を取得する
	var creater string
	err = h.db.Get(&creater, "SELECT creater FROM todos WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// todoの作成者であるか確かめる
	if creater != c.Get("userName").(string) {
		return c.String(http.StatusBadRequest, "you are not this todo's creater")
	}

	// 完了状態にする。すでに完了していてもエラーは発生しない
	_, err = h.db.Exec("UPDATE todos SET is_finished = 1 WHERE id = ?", id)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
