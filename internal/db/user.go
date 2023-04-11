package db

import (
	"database/sql"
	"fmt"
	"github.com/harsha-aqfer/todo/pkg"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserDB interface {
	CreateUser(ui *pkg.User) error
	GetUser(email string) (*pkg.User, error)
	GetUserID(email string) (int64, error)
}

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserDB {
	return &userStore{db: db}
}

func (us *userStore) CreateUser(ui *pkg.User) error {
	_, err := us.db.Exec("INSERT user SET email = ?, user_name = ?, password = ?", ui.Email, ui.Username, ui.Password)
	return err
}

func (us *userStore) GetUser(email string) (*pkg.User, error) {
	row := us.db.QueryRow("SELECT email, user_name, password FROM user WHERE email = ?", email)

	r := pkg.User{}
	err := row.Scan(&r.Email, &r.Username, &r.Password)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no such user: %s", email)
	} else if err != nil {
		return nil, err
	}
	return &r, nil
}

func (us *userStore) GetUserID(email string) (int64, error) {
	row := us.db.QueryRow("SELECT id FROM user WHERE email = ?", email)

	var r int64
	err := row.Scan(&r)

	if err == sql.ErrNoRows {
		return 0, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("no such user: %s", email))
	} else if err != nil {
		return 0, err
	}
	return r, nil
}
