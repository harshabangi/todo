package service

import (
	"encoding/json"
	"github.com/harsha-aqfer/todo/internal/db"
	"github.com/harsha-aqfer/todo/pkg"
	asserts "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTodoDB struct {
	mock.Mock
}

func (t *mockTodoDB) ListTodos(all bool) ([]pkg.TodoResponse, error) {
	args := t.Called(all)
	return args.Get(0).([]pkg.TodoResponse), args.Error(1)
}

func (t *mockTodoDB) CreateTodo(tr *pkg.TodoRequest) error {
	args := t.Called(tr)
	return args.Error(0)
}

func (t *mockTodoDB) GetTodo(id int64) (*pkg.TodoResponse, error) {
	args := t.Called(id)
	return args.Get(0).(*pkg.TodoResponse), args.Error(1)
}

func (t *mockTodoDB) UpdateTodo(id int64, tr *pkg.TodoRequest) error {
	args := t.Called(id, tr)
	return args.Error(0)
}

func (t *mockTodoDB) DeleteTodo(id int64) error {
	args := t.Called(id)
	return args.Error(0)
}

func Test_ListTodos(t *testing.T) {
	assert := asserts.New(t)

	var (
		md = &mockTodoDB{}
		s  = &Service{db: &db.DB{Todo: md}}

		rr = httptest.NewRecorder()
		rq = httptest.NewRequest(http.MethodGet, "http://localhost:3000/todos?all=true", nil)

		actualResponse []pkg.TodoResponse
	)

	dbOut := []pkg.TodoResponse{
		{Id: 1, Task: "task-1", Category: "work", Priority: "low"},
	}
	md.On("ListTodos", true).Return(dbOut, nil)

	err := s.listTodos(rr, rq)
	assert.Nil(err)

	response := rr.Result()
	bb, err := io.ReadAll(response.Body)
	assert.Nil(err)
	err = json.Unmarshal(bb, &actualResponse)
	assert.Nil(err)

	assert.Equal("200 OK", response.Status)
	assert.Equal(dbOut, actualResponse)
}
