package sessions_test

import (
	"testing"
	"time"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/sessions"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

type Session = repository.Session

func Test_Validate(t *testing.T) {
	tests := []struct {
		name     string
		data     Session
		token    string
		expected bool
	}{
		{
			name:     "returns true",
			data:     Session{Id: "someToken", CreatedAt: time.Now()},
			token:    "someToken",
			expected: true,
		},
		{
			name:     "returns false",
			token:    "someToken",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := util.InitMockSessionsRepository()
			repo.PopulateData(tc.data)
			usecase := sessions.InitSessionsUsecase(repo)

			got := usecase.Validate(tc.token)

			util.AssertEqual(t)(got, tc.expected)
		})
	}
}
