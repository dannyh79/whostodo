package sessions_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/google/go-cmp/cmp"
)

type Session = repository.Session

type MockSessionsRepository struct {
	data map[string]Session
}

var mockNotFoundError = errors.New("not found")

func (r *MockSessionsRepository) FindBy(id any) (*Session, error) {
	row, ok := r.data[id.(string)]
	if !ok {
		return nil, mockNotFoundError
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
			repo := initMockSessionsRepository()
			repo.PopulateData(tc.data)
			usecase := sessions.InitSessionsUsecase(repo)

			got := usecase.Validate(tc.token)

			assertEqual(t)(got, tc.expected)
		})
	}
}

func assertEqual(t *testing.T) func(got any, want any) {
	return func(got any, want any) {
		t.Helper()
		if !cmp.Equal(got, want) {
			t.Errorf(cmp.Diff(want, got))
		}
	}
}
