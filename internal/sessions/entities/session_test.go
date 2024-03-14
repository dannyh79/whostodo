package entity_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/sessions/entities"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

func Test_NewSession(t *testing.T) {
	s1 := entity.NewSession()
	s2 := entity.NewSession()

	util.AssertNotEqual(t)(s1.Id, s2.Id)
}
