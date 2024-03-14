package sessions_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/sessions"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

type Session = repository.Session

func Test_Authenticate(t *testing.T) {
	t.Parallel()

	repo := util.InitMockSessionsRepository()
	usecase := sessions.InitSessionsUsecase(repo)

	token := usecase.Authenticate()

	util.AssertEqual(t)(token, repo.Data[token].Id)
}

func Test_Validate(t *testing.T) {
	tests := []struct {
		name     string
		data     Session
		expected bool
	}{
		{
			name:     "returns true",
			data:     util.NewSession(),
			expected: true,
		},
		{
			name:     "returns false",
			data:     util.NewExpiredSession(),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := util.InitMockSessionsRepository()
			repo.PopulateData(tc.data)
			usecase := sessions.InitSessionsUsecase(repo)

			got := usecase.Validate(tc.data.Id)

			util.AssertEqual(t)(got, tc.expected)
		})
	}
}
