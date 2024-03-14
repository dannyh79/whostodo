package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/rest/v1"
	"github.com/dannyh79/whostodo/internal/sessions/entities"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

type Session = entity.Session

func Test_GETTasks(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		session    Session
		data       []repository.TaskSchema
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 200 with result",
			authroized: true,
			session:    util.NewSession(),
			data:       []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			statusCode: http.StatusOK,
			expected:   `{"result":[{"id":1,"name":"name","status":0}]}`,
		},
		{
			name:       "returns status code 200 with empty result",
			authroized: true,
			session:    util.NewSession(),
			statusCode: http.StatusOK,
			expected:   `{"result":[]}`,
		},
		{
			name:       "without session token returns status code 403",
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
		{
			name:       "with expired session returns status code 403",
			authroized: true,
			session:    util.NewExpiredSession(),
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := util.NewTestSuite()
			if data := tc.data; len(data) > 0 {
				for _, row := range data {
					suite.TaskRepo.PopulateData(row)
				}
			}
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/v1/tasks", nil)
			if tc.authroized {
				suite.SessionRepo.PopulateData(tc.session)
				setRequestTokenHeader(t)(req, tc.session.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t)(rr)
			util.AssertHttpStatus(t)(rr, tc.statusCode)
			util.AssertEqual(t)(rr.Body.String(), tc.expected)
		})
	}
}

func Test_POSTTask(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		session    Session
		data       string
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 201 with result",
			authroized: true,
			session:    util.NewSession(),
			data:       `{"name":"買晚餐"}`,
			statusCode: http.StatusCreated,
			expected:   `{"result":{"name":"買晚餐","status":0,"id":1}}`,
		},
		{
			name:       "without session token returns status code 403",
			authroized: false,
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
		{
			name:       "with expired session returns status code 403",
			authroized: true,
			session:    util.NewExpiredSession(),
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := util.NewTestSuite()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/v1/task", bytes.NewBufferString(tc.data))
			req.Header.Add("Content-Type", "application/json")
			if tc.authroized {
				suite.SessionRepo.PopulateData(tc.session)
				setRequestTokenHeader(t)(req, tc.session.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t)(rr)
			util.AssertHttpStatus(t)(rr, tc.statusCode)
			util.AssertEqual(t)(rr.Body.String(), tc.expected)
		})
	}
}

func Test_PUTTask(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		session    Session
		data       repository.TaskSchema
		param      int
		payload    string
		statusCode int
		expected   string
	}{
		{
			name:       "returns status code 200 with result",
			authroized: true,
			session:    util.NewSession(),
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      1,
			payload:    `{"name":"買晚餐","status":1}`,
			statusCode: http.StatusCreated,
			expected:   `{"result":{"name":"買晚餐","status":1,"id":1}}`,
		},
		{
			name:       "returns status code 404 with empty result",
			authroized: true,
			session:    util.NewSession(),
			data:       repository.TaskSchema{},
			param:      1,
			payload:    `{"name":"買晚餐","status":1}`,
			statusCode: http.StatusNotFound,
			expected:   `{"result":{}}`,
		},
		{
			name:       "without session token returns status code 403",
			authroized: false,
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
		{
			name:       "with expired session returns status code 403",
			authroized: true,
			session:    util.NewExpiredSession(),
			statusCode: http.StatusForbidden,
			expected:   `{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := util.NewTestSuite()
			suite.TaskRepo.PopulateData(tc.data)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/v1/task/%d", tc.param),
				bytes.NewBufferString(tc.payload),
			)
			req.Header.Add("Content-Type", "application/json")
			if tc.authroized {
				suite.SessionRepo.PopulateData(tc.session)
				setRequestTokenHeader(t)(req, tc.session.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t)(rr)
			util.AssertHttpStatus(t)(rr, tc.statusCode)
			util.AssertEqual(t)(rr.Body.String(), tc.expected)
		})
	}
}

func Test_DELETETask(t *testing.T) {
	tests := []struct {
		name       string
		authroized bool
		session    Session
		data       repository.TaskSchema
		param      int
		statusCode int
	}{
		{
			name:       "returns status code 200",
			authroized: true,
			session:    util.NewSession(),
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      1,
			statusCode: http.StatusOK,
		},
		{
			name:       "returns status code 404",
			authroized: true,
			session:    util.NewSession(),
			data:       repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:      2,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "without session token returns status code 403",
			authroized: false,
			statusCode: http.StatusForbidden,
		},
		{
			name:       "with expired session returns status code 403",
			authroized: true,
			session:    util.NewExpiredSession(),
			statusCode: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := util.NewTestSuite()
			suite.TaskRepo.PopulateData(tc.data)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/task/%d", tc.param), nil)
			if tc.authroized {
				suite.SessionRepo.PopulateData(tc.session)
				setRequestTokenHeader(t)(req, tc.session.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t)(rr)
			util.AssertHttpStatus(t)(rr, tc.statusCode)
		})
	}
}

func Test_POSTAuth(t *testing.T) {
	tests := []struct {
		name             string
		authroized       bool
		session          Session
		statusCode       int
		expectNewSession bool
	}{
		{
			name:             "returns status code 304 with empty result",
			authroized:       true,
			session:          util.NewSession(),
			statusCode:       http.StatusNotModified,
			expectNewSession: false,
		},
		{
			name:             "returns status code 201 with token",
			authroized:       false,
			session:          util.NewSession(),
			statusCode:       http.StatusCreated,
			expectNewSession: true,
		},
		{
			name:             "returns status code 201 with new token after 1 minute",
			authroized:       true,
			session:          util.NewExpiredSession(),
			statusCode:       http.StatusCreated,
			expectNewSession: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			suite := util.NewTestSuite()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/v1/auth", nil)
			if tc.authroized {
				suite.SessionRepo.PopulateData(tc.session)
				setRequestTokenHeader(t)(req, tc.session.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t)(rr)
			util.AssertHttpStatus(t)(rr, tc.statusCode)
			if tc.expectNewSession {
				token := getTokenFromResponse(rr)
				util.AssertNotEqual(t)(token, tc.session.Id)
			} else {
				util.AssertEqual(t)(rr.Body.String(), "")
			}
		})
	}
}

func setRequestTokenHeader(t testing.TB) func(r *http.Request, token string) {
	return func(r *http.Request, token string) {
		t.Helper()
		r.Header.Set("Authorization", "Bearer "+token)
	}
}

func getTokenFromResponse(r *httptest.ResponseRecorder) string {
	var body routes.PostAuthSuccessOutput
	_ = json.Unmarshal(r.Body.Bytes(), &body)
	return body.Token
}
