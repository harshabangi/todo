package pkg

import (
	"fmt"
	"github.com/harsha-aqfer/todo/internal/util"
	"strings"
	"time"
)

type TodoRequest struct {
	Task     string `json:"task"`
	Done     bool   `json:"done,omitempty"`
	Category string `json:"category,omitempty"`
	Priority string `json:"priority,omitempty"`
}

type TodoResponse struct {
	Id          int64      `json:"id"`
	Task        string     `json:"task"`
	Category    string     `json:"category"`
	Priority    string     `json:"priority"`
	CreatedAt   *time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func (tr *TodoRequest) IsZero() bool {
	return tr.Task == "" &&
		tr.Priority == "" &&
		tr.Category == "" &&
		tr.Done == false
}

func (tr *TodoRequest) Validate() error {
	if tr.Task == "" {
		return fmt.Errorf("inadequate input parameters. Required field: task")
	}

	category := tr.Category
	tr.Category = strings.ToLower(tr.Category)

	categories := []string{"work", "home"}

	if !util.Contains(categories, tr.Category) {
		return fmt.Errorf("unknown category value: %s", category)
	}

	pr := tr.Priority
	tr.Priority = strings.ToLower(tr.Priority)

	priorities := []string{"low", "medium", "high"}

	if !util.Contains(priorities, tr.Priority) {
		return fmt.Errorf("unknown priority value: %s", pr)
	}
	return nil
}

type User struct {
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (u User) Validate() error {
	s := []string{u.Email, u.Password, u.Username}

	if util.Contains(s, "") {
		return fmt.Errorf("inadequate input parameters. Required email, username, password")
	}
	return nil
}

type MsgResp struct {
	Message string `json:"message"`
}

func NewMsgResp(message string) *MsgResp {
	return &MsgResp{Message: message}
}

type Token struct {
	Type      string `json:"type"`
	ExpiresIn int    `json:"expires_in"`
	JWTToken  string `json:"jwt_token"`
}
