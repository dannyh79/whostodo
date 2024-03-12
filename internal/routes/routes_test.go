package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/gin-gonic/gin"
)

type MockTestSuite struct {
	engine *gin.Engine
}

func newTestSuite() *MockTestSuite {
	engine := gin.Default()

	routes.AddRoutes(engine)

	return &MockTestSuite{
		engine: engine,
	}
}

func Test_GETTasks(t *testing.T) {
	t.Run("returns status code 200 with result", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

		suite := newTestSuite()
		suite.engine.ServeHTTP(rr, req)

		assertHttpStatus(t, rr, http.StatusOK)

		expectedBody := `{"result":[{"id":1,"name":"name","status":0}]}`
		assertResponseBody(t, rr.Body.String(), expectedBody)
	})
}

func assertHttpStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()

	got := rr.Result().StatusCode
	if got != want {
		t.Errorf("got HTTP status %v, want %v", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
}
