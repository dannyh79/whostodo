package routes_test

import (
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
	repo *MockTaskRepository
}

type MockTaskRepository struct {
	data []repository.TaskSchema
}

func (r *MockTaskRepository) ListAll() []*entity.Task {
	var tasks []*entity.Task
	for _, row := range r.data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func (r *MockTaskRepository) PopulateData(data []repository.TaskSchema) {
	for _, row := range(data) {
		r.data = append(r.data, row)
	}
}

func newTestSuite() *MockTestSuite {
	engine := gin.Default()

	repo := &MockTaskRepository{
		data: []repository.TaskSchema{},
	}
	usecase := tasks.InitTasksUsecase(repo)

	routes.AddRoutes(engine, usecase)

	return &MockTestSuite{
		engine: engine,
		repo: repo,
	}
}

func Test_GETTasks(t *testing.T) {
	tests := []struct{
		name string
		data []repository.TaskSchema
		statusCode int
		expected string
	}{
		{
			name: "returns status code 200 with result",
			data: []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			statusCode: http.StatusOK,
			expected: `{"result":[{"id":1,"name":"name","status":0}]}`,
		},
		{
			name: "returns status code 200 with empty result",
			statusCode: http.StatusOK,
			expected: `{"result":[]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

			suite := newTestSuite()
			if (len(tc.data) > 0) {
				suite.repo.PopulateData(tc.data)
			}
			suite.engine.ServeHTTP(rr, req)

			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
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
