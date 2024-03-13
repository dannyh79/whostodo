package entity_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/sessions/entities"
	"github.com/google/go-cmp/cmp"
)

func Test_NewSession(t *testing.T) {
	s1 := entity.NewSession()
	s2 := entity.NewSession()

	if cmp.Equal(s1.Id, s2.Id) {
		t.Error(s1)
		t.Error(s2)
	}
}
