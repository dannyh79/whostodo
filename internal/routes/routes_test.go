package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

type MockTestSuite struct {
	engine *gin.Engine
}

func newTestSuite() *MockTestSuite {
	engine := gin.Default()

	repo := repository.InitInMemoryTaskRepository()
	usecase := tasks.InitTasksUsecase(repo)

	routes.AddRoutes(engine, usecase)

	return &MockTestSuite{
		engine: engine,
	}
}

func Test_GETTasks(t *testing.T) {
	tests := []struct{
		name string
		statusCode int
		expected string
	}{
		{
			name: "returns status code 200 with result",
			statusCode: http.StatusOK,
			expected: `{"result":[{"id":1,"name":"name","status":0}]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

			suite := newTestSuite()
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
