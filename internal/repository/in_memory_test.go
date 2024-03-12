package repository_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/google/go-cmp/cmp"
)

func Test_InMemoryListAll(t *testing.T) {
	repo := repository.InitInMemoryRepo()
	tasks := repo.ListAll()

	if !cmp.Equal(0, len(tasks)) {
		t.Fail()
	}
}
