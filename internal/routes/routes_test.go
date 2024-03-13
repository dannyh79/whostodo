package routes_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

type MockTestSuite struct {
	engine      *gin.Engine
	taskRepo    *MockTaskRepository
	sessionRepo *MockSessionsRepository
}

type MockTaskRepository struct {
	data map[int]repository.TaskSchema
}

func (r *MockTaskRepository) FindBy(id any) (*entity.Task, error) {
	row, ok := r.data[id.(int)]
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

func (r *MockTaskRepository) Delete(t *entity.Task) error {
	return nil
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

type Session = repository.Session

type MockSessionsRepository struct {
	data map[string]Session
}

var mockSessionNotFoundError = errors.New("session not found")

func (r *MockSessionsRepository) FindBy(id any) (*Session, error) {
	if id == nil {
		return nil, mockSessionNotFoundError
	}

	row, ok := r.data[id.(string)]
	if !ok {
		return nil, mockSessionNotFoundError
	}

	return &Session{Id: row.Id, CreatedAt: row.CreatedAt}, nil
}

func (r *MockSessionsRepository) Update(s *Session) (*Session, error) {
	panic("not implemented")
}

func (r *MockSessionsRepository) Save(s *Session) Session {
	r.data[s.Id] = *s
	return *s
}

func (r *MockSessionsRepository) Delete(t *Session) error {
	panic("not implemented")
}

func (r *MockSessionsRepository) ListAll() []*Session {
	panic("not implemented")
}

func (r *MockSessionsRepository) PopulateData(row Session) {
	r.data[row.Id] = row
}

func initMockSessionsRepository() *MockSessionsRepository {
	return &MockSessionsRepository{
		data: make(map[string]Session),
	}
}

var stubbedSession = Session{Id: "stubbed_token"}

func newTestSuite() *MockTestSuite {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	taskRepo := &MockTaskRepository{
		data: make(map[int]repository.TaskSchema),
	}
	sessionRepo := &MockSessionsRepository{
		data: make(map[string]Session),
	}
	tasksUsecase := tasks.InitTasksUsecase(taskRepo)
	sessionsUsecase := sessions.InitSessionsUsecase(sessionRepo)

	routes.AddRoutes(engine, tasksUsecase, sessionsUsecase)

	return &MockTestSuite{
		engine:      engine,
		taskRepo:    taskRepo,
		sessionRepo: sessionRepo,
	}
}

func Test_GETTasks(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		data       []repository.TaskSchema
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 200 with result",
			authroized: true,
			data:       []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			statusCode: http.StatusOK,
			expected:   `{"result":[{"id":1,"name":"name","status":0}]}`,
		},
		{
			name:       "returns status code 200 with empty result",
			authroized: true,
			statusCode: http.StatusOK,
			expected:   `{"result":[]}`,
		},
		{
			name:       "returns status code 403",
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := newTestSuite()
			if data := tc.data; len(data) > 0 {
				for _, row := range data {
					suite.taskRepo.PopulateData(row)
				}
			}
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
			if tc.authroized {
				suite.sessionRepo.PopulateData(stubbedSession)
				setRequestTokenHeader(t)(req, stubbedSession.Id)
			}

			suite.engine.ServeHTTP(rr, req)

			assertJsonHeader(t, rr)
			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
		})
	}
}

func Test_POSTTask(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		data       string
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 201 with result",
			authroized: true,
			data:       `{"name":"買晚餐"}`,
			statusCode: http.StatusCreated,
			expected:   `{"result":{"name":"買晚餐","status":0,"id":1}}`,
		},
		{
			name:       "returns status code 403",
			authroized: false,
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			suite := newTestSuite()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBufferString(tc.data))
			req.Header.Add("Content-Type", "application/json")
			if tc.authroized {
				suite.sessionRepo.PopulateData(stubbedSession)
				setRequestTokenHeader(t)(req, stubbedSession.Id)
			}

			suite.engine.ServeHTTP(rr, req)

			assertJsonHeader(t, rr)
			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
		})
	}
}

func Test_PUTTask(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		data       repository.TaskSchema
		param      int
		payload    string
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 200 with result",
			authroized: true,
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      1,
			payload:    `{"name":"買晚餐","status":1}`,
			statusCode: http.StatusCreated,
			expected:   `{"result":{"name":"買晚餐","status":1,"id":1}}`,
		},
		{
			name:       "returns status code 404 with empty result",
			authroized: true,
			data:       repository.TaskSchema{},
			param:      1,
			payload:    `{"name":"買晚餐","status":1}`,
			statusCode: http.StatusNotFound,
			expected:   `{"result":{}}`,
		},
		{
			name:       "returns status code 403",
			authroized: false,
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			suite := newTestSuite()
			suite.taskRepo.PopulateData(tc.data)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/task/%d", tc.param),
				bytes.NewBufferString(tc.payload),
			)
			req.Header.Add("Content-Type", "application/json")
			if tc.authroized {
				suite.sessionRepo.PopulateData(stubbedSession)
				setRequestTokenHeader(t)(req, stubbedSession.Id)
			}

			suite.engine.ServeHTTP(rr, req)

			assertJsonHeader(t, rr)
			assertHttpStatus(t, rr, tc.statusCode)
			assertResponseBody(t, rr.Body.String(), tc.expected)
		})
	}
}

func Test_DELETETask(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		data       repository.TaskSchema
		param      int
		statusCode int
	}{
		{
			name:       "returns status code 200",
			authroized: true,
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      1,
			statusCode: http.StatusOK,
		},
		{
			name:       "returns status code 404",
			authroized: true,
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      2,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "returns status code 403",
			authroized: false,
			statusCode: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			suite := newTestSuite()
			suite.taskRepo.PopulateData(tc.data)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/task/%d", tc.param), nil)
			if tc.authroized {
				suite.sessionRepo.PopulateData(stubbedSession)
				setRequestTokenHeader(t)(req, stubbedSession.Id)
			}

			suite.engine.ServeHTTP(rr, req)

			assertJsonHeader(t, rr)
			assertHttpStatus(t, rr, tc.statusCode)
		})
	}
}

func Test_POSTAuth(t *testing.T) {
	t.Run("returns status code 304 with empty result", func(t *testing.T) {
		t.Parallel()

		suite := newTestSuite()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth", nil)
		suite.sessionRepo.PopulateData(stubbedSession)
		setRequestTokenHeader(t)(req, stubbedSession.Id)

		suite.engine.ServeHTTP(rr, req)

		assertJsonHeader(t, rr)
		assertHttpStatus(t, rr, http.StatusNotModified)
		assertResponseBody(t, rr.Body.String(), "")
	})

	t.Run("returns status code 201 with token", func(t *testing.T) {
		t.Parallel()

		suite := newTestSuite()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth", nil)

		suite.engine.ServeHTTP(rr, req)

		assertJsonHeader(t, rr)
		assertHttpStatus(t, rr, http.StatusCreated)

		if emptyResult := `{"result":""}`; cmp.Equal(rr.Body.String(), emptyResult) {
			t.Error()
		}
	})
}

func assertJsonHeader(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	want := "application/json"
	if got := rr.Header().Get("Content-Type"); !strings.Contains(got, want) {
		t.Errorf("missed Content-Type %v in %v", want, got)
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

func setRequestTokenHeader(t testing.TB) func(r *http.Request, token string) {
	return func(r *http.Request, token string) {
		t.Helper()
		r.Header.Set("Authorization", "Bearer "+token)
	}
}
