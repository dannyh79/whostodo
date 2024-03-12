package routes_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

type MockTestSuite struct {
	engine *gin.Engine
	repo   *MockTaskRepository
}

type MockTaskRepository struct {
	data map[int]repository.TaskSchema
}

func (r *MockTaskRepository) FindBy(id int) (*entity.Task, error) {
	row, ok := r.data[id]
	if !ok {
		return nil, errors.New("")
	}
	return &entity.Task{Id: row.Id, Name: row.Name, Status: row.Status}, nil
}

func (r *MockTaskRepository) Update(t *entity.Task) (*entity.Task, error) {
	r.data[t.Id] = repository.TaskSchema{Id: t.Id, Name: t.Name, Status: t.Status}
	return &entity.Task{Id: t.Id, Name: t.Name, Status: t.Status}, nil
}

func (r *MockTaskRepository) Save(t *entity.Task) entity.Task {
	t.Id = len(r.data) + 1
	return *t
}

func (r *MockTaskRepository) ListAll() []*entity.Task {
	var tasks []*entity.Task
	for _, row := range r.data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func (r *MockTaskRepository) PopulateData(row repository.TaskSchema) {
	r.data[row.Id] = row
}

func newTestSuite() *MockTestSuite {
	engine := gin.Default()

	repo := &MockTaskRepository{
		data: make(map[int]repository.TaskSchema),
	}
	usecase := tasks.InitTasksUsecase(repo)

	routes.AddRoutes(engine, usecase)

	return &MockTestSuite{
		engine: engine,
		repo:   repo,
	}
}

func Test_GETTasks(t *testing.T) {
	tests := []struct {
		name       string
		data       []repository.TaskSchema
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 200 with result",
			data:       []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			statusCode: http.StatusOK,
			expected:   `{"result":[{"id":1,"name":"name","status":0}]}`,
		},
		{
			name:       "returns status code 200 with empty result",
			statusCode: http.StatusOK,
			expected:   `{"result":[]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

			suite := newTestSuite()
			if data := tc.data; len(data) > 0 {
				for _, row := range data {
					suite.repo.PopulateData(row)
				}
			}
			suite.engine.ServeHTTP(rr, req)

			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
		})
	}
}

func Test_POSTTask(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 201 with result",
			data:       `{"name":"買晚餐"}`,
			statusCode: http.StatusCreated,
			expected:   `{"result":{"name":"買晚餐","status":0,"id":1}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBufferString(tc.data))
			req.Header.Add("Content-Type", "application/json")

			suite := newTestSuite()
			suite.engine.ServeHTTP(rr, req)

			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
		})
	}
}

func Test_PUTTask(t *testing.T) {
	tests := []struct {
		name       string
		data       repository.TaskSchema
		param      int
		payload    string
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 200 with result",
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      1,
			payload:    `{"name":"買晚餐","status":1}`,
			statusCode: http.StatusCreated,
			expected:   `{"result":{"name":"買晚餐","status":1,"id":1}}`,
		},
		{
			name:       "returns status code 404 with empty result",
			data:       repository.TaskSchema{},
			param:      1,
			payload:    `{"name":"買晚餐","status":1}`,
			statusCode: http.StatusNotFound,
			expected:   `{"result":{}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/task/%d", tc.param),
				bytes.NewBufferString(tc.payload),
			)
			req.Header.Add("Content-Type", "application/json")

			suite := newTestSuite()
			suite.repo.PopulateData(tc.data)
			suite.engine.ServeHTTP(rr, req)

			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
		})
	}
}

func Test_DELETETask(t *testing.T) {
	tests := []struct {
		name       string
		data       repository.TaskSchema
		param      int
		statusCode int
	}{
		{
			name:       "returns status code 200",
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      1,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/task/%d", tc.param), nil)

			suite := newTestSuite()
			suite.repo.PopulateData(tc.data)
			suite.engine.ServeHTTP(rr, req)

			assertHttpStatus(t, rr, tc.statusCode)
		})
	}
}

func assertHttpStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got := rr.Result().StatusCode; got != want {
		t.Errorf("got HTTP status %v, want %v", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}
