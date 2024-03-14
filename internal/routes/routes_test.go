package routes_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

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

			suite := util.NewTestSuite()
			if data := tc.data; len(data) > 0 {
				for _, row := range data {
					suite.TaskRepo.PopulateData(row)
				}
			}
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
			if tc.authroized {
				suite.SessionRepo.PopulateData(util.StubbedSession)
				setRequestTokenHeader(t)(req, util.StubbedSession.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t, rr)
			util.AssertHttpStatus(t, rr, tc.statusCode)
			util.AssertResponseBody(t, rr.Body.String(), tc.expected)
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
			suite := util.NewTestSuite()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBufferString(tc.data))
			req.Header.Add("Content-Type", "application/json")
			if tc.authroized {
				suite.SessionRepo.PopulateData(util.StubbedSession)
				setRequestTokenHeader(t)(req, util.StubbedSession.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t, rr)
			util.AssertHttpStatus(t, rr, tc.statusCode)
			util.AssertResponseBody(t, rr.Body.String(), tc.expected)
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
			suite := util.NewTestSuite()
			suite.TaskRepo.PopulateData(tc.data)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/task/%d", tc.param),
				bytes.NewBufferString(tc.payload),
			)
			req.Header.Add("Content-Type", "application/json")
			if tc.authroized {
				suite.SessionRepo.PopulateData(util.StubbedSession)
				setRequestTokenHeader(t)(req, util.StubbedSession.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t, rr)
			util.AssertHttpStatus(t, rr, tc.statusCode)
			util.AssertResponseBody(t, rr.Body.String(), tc.expected)
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
			suite := util.NewTestSuite()
			suite.TaskRepo.PopulateData(tc.data)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/task/%d", tc.param), nil)
			if tc.authroized {
				suite.SessionRepo.PopulateData(util.StubbedSession)
				setRequestTokenHeader(t)(req, util.StubbedSession.Id)
			}

			suite.Engine.ServeHTTP(rr, req)

			util.AssertJsonHeader(t, rr)
			util.AssertHttpStatus(t, rr, tc.statusCode)
		})
	}
}

func Test_POSTAuth(t *testing.T) {
	t.Run("returns status code 304 with empty result", func(t *testing.T) {
		t.Parallel()

		suite := util.NewTestSuite()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth", nil)
		suite.SessionRepo.PopulateData(util.StubbedSession)
		setRequestTokenHeader(t)(req, util.StubbedSession.Id)

		suite.Engine.ServeHTTP(rr, req)

		util.AssertJsonHeader(t, rr)
		util.AssertHttpStatus(t, rr, http.StatusNotModified)
		util.AssertResponseBody(t, rr.Body.String(), "")
	})

	t.Run("returns status code 201 with token", func(t *testing.T) {
		t.Parallel()

		suite := util.NewTestSuite()
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth", nil)

		suite.Engine.ServeHTTP(rr, req)

		util.AssertJsonHeader(t, rr)
		util.AssertHttpStatus(t, rr, http.StatusCreated)
		util.AssertNotEqual(t)(rr.Body.String(), `{"result":""}`)
	})
}

func setRequestTokenHeader(t testing.TB) func(r *http.Request, token string) {
	return func(r *http.Request, token string) {
		t.Helper()
		r.Header.Set("Authorization", "Bearer "+token)
	}
}
