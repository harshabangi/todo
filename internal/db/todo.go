package db

import (
	"database/sql"
	"fmt"
	"github.com/harsha-aqfer/todo/pkg"
	"strings"
	"time"
)

type TodoDB interface {
	ListTodos(userID int64, all bool) ([]pkg.TodoResponse, error)
	GetTodo(userID, todoID int64) (*pkg.TodoResponse, error)
	CreateTodo(userID int64, tr *pkg.TodoRequest) error
	UpdateTodo(userID, todoID int64, tr *pkg.TodoRequest) error
	DeleteTodo(userID, todoID int64) error
}

type todoStore struct {
	db *sql.DB
}

func NewTodoStore(db *sql.DB) TodoDB {
	return &todoStore{db: db}
}

func (ts *todoStore) ListTodos(userID int64, all bool) ([]pkg.TodoResponse, error) {
	query := "SELECT id, task, category, priority, created_at, completed_at FROM todo WHERE user_id = ?"

	if !all {
		query += " AND NOT done"
	}
	rows, err := ts.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	todos := make([]pkg.TodoResponse, 0)

	for rows.Next() {
		t := pkg.TodoResponse{}
		var ct sql.NullTime

		err = rows.Scan(&t.Id, &t.Task, &t.Category, &t.Priority, &t.CreatedAt, &ct)

		if err != nil {
			return nil, err
		}
		if ct.Valid {
			t.CompletedAt = &ct.Time
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (ts *todoStore) CreateTodo(userID int64, tr *pkg.TodoRequest) error {
	var (
		query  = "INSERT todo SET user_id = ?, task = ?"
		params = []interface{}{userID, tr.Task}
	)

	if tr.Category != "" {
		query += ", category = ?"
		params = append(params, tr.Category)
	}

	if tr.Priority != "" {
		query += ", priority = ?"
		params = append(params, tr.Priority)
	}

	_, err := ts.db.Exec(query, params...)
	return err
}

func (ts *todoStore) GetTodo(userID, todoID int64) (*pkg.TodoResponse, error) {
	query := "SELECT id, task, category, priority, created_at, completed_at FROM todo WHERE user_id = ? AND id = ?"

	rows, err := ts.db.Query(query, userID, todoID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		t := pkg.TodoResponse{}
		var ct sql.NullTime

		err = rows.Scan(&t.Id, &t.Task, &t.Category, &t.Priority, &t.CreatedAt, &ct)

		if err != nil {
			return nil, err
		}
		if ct.Valid {
			t.CompletedAt = &ct.Time
		}
		return &t, nil
	}
	return nil, nil
}

func (ts *todoStore) UpdateTodo(userID, todoID int64, tr *pkg.TodoRequest) error {
	var (
		qs     []string
		params []interface{}
	)

	if tr.Task != "" {
		qs = append(qs, "task = ?")
		params = append(params, tr.Task)
	}

	if tr.Category != "" {
		qs = append(qs, "category = ?")
		params = append(params, tr.Category)
	}

	if tr.Priority != "" {
		qs = append(qs, "priority = ?")
		params = append(params, tr.Priority)
	}

	if tr.Done {
		qs = append(qs, "done = ?")
		params = append(params, int64(1))

		qs = append(qs, "completed_at = ?")
		params = append(params, time.Now().UTC())
	}

	params = append(params, todoID, userID)
	_, err := ts.db.Exec(fmt.Sprintf("UPDATE todo SET %s WHERE id = ? AND user_id = ?", strings.Join(qs, ", ")), params...)
	return err
}

func (ts *todoStore) DeleteTodo(userID, todoID int64) error {
	_, err := ts.db.Exec("DELETE FROM todo WHERE user_id = ? AND id = ?", userID, todoID)
	return err
}
