package testutil_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertEqual(t *testing.T) func(got any, want any) {
	return func(got any, want any) {
		t.Helper()
		if !cmp.Equal(got, want) {
			t.Errorf(cmp.Diff(want, got))
		}
	}
}

func AssertErrorEqual(t *testing.T) func(got error, want error) {
	return func(got error, want error) {
		t.Helper()
		if want != nil && !errors.Is(got, want) {
			t.Errorf(cmp.Diff(got.Error(), want.Error()))
		}
	}
}

func AssertJsonHeader(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	want := "application/json"
	if got := rr.Header().Get("Content-Type"); !strings.Contains(got, want) {
		t.Errorf("missed Content-Type %v in %v", want, got)
	}
}

func AssertHttpStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got := rr.Result().StatusCode; got != want {
		t.Errorf("got HTTP status %v, want %v", got, want)
	}
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}
