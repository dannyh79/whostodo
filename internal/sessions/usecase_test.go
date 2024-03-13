package sessions_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/google/go-cmp/cmp"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "returns true",
			token:    "someToken",
			expected: true,
		},
		{
			name:     "returns false",
			token:    "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			usecase := sessions.InitSessionsUsecase()

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
