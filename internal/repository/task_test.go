package repository_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/google/go-cmp/cmp"
)

func Test_InMemoryTaskRepositoryListAll(t *testing.T) {
	repo := repository.InitInMemoryTaskRepository()
	tasks := repo.ListAll()

	if !cmp.Equal(0, len(tasks)) {
		t.Fail()
	}
}

func Test_InMemoryTaskRepositorySave(t *testing.T) {
	repo := repository.InitInMemoryTaskRepository()
	repo.Save(&entity.Task{Id: 1, Name: "買晚餐", Status: 0})

	tasks := repo.ListAll()
	if !cmp.Equal(1, len(tasks)) {
		t.Fail()
	}
}

func Test_InMemoryTaskRepositoryNextId(t *testing.T) {
	repo := repository.InitInMemoryTaskRepository()
	if !cmp.Equal(1, repo.NextId()) {
		t.Fail()
	}
}
