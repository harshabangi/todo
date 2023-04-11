package service_echo

import (
	"fmt"
	"github.com/harsha-aqfer/todo/internal/db"
	"github.com/labstack/echo/v4"
)

type Config struct {
	UserName   string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Host       string `json:"host"`
	ListenAddr string `json:"listen_addr"`
	SigningKey string `json:"signing_key"`
}

func NewConfig() *Config {
	return &Config{}
}

type Service struct {
	conf *Config
	db   *db.DB
}

func NewService(c *Config) (*Service, error) {
	store, err := db.NewDB(c.UserName, c.Password, c.Host, c.Database)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Service{
		conf: c,
		db:   store,
	}, nil
}

func (s *Service) Run() {
	e := echo.New()

	// Register app (*App) to be injected into all HTTP handlers.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("service", s)
			return next(c)
		}
	})

	e.POST("/v1/sign_up", signUp)
	e.POST("/v1/sign_in", signIn)

	todoGrp := e.Group("")
	todoGrp.Use(IsAuthorized)

	todoGrp.POST("/v1/todos", createTodo)
	todoGrp.GET("/v1/todos", listTodos)

	todoGrp.GET("/v1/todos/:id", getTodo)
	todoGrp.PUT("/v1/todos/:id", updateTodo)
	todoGrp.DELETE("/v1/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(s.conf.ListenAddr))
}
