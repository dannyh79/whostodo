package repository_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

func Test_InMemoryTaskRepositoryListAll(t *testing.T) {
	t.Parallel()

	repo := repository.InitInMemoryTaskRepository()

	tasks := repo.ListAll()

	util.AssertEqual(t)(len(tasks), 0)
}

func Test_InMemoryTaskRepositorySave(t *testing.T) {
	t.Parallel()

	repo := repository.InitInMemoryTaskRepository()
	repo.Save(&entity.Task{Id: 1, Name: "買晚餐", Status: 0})

	tasks := repo.ListAll()

	util.AssertEqual(t)(len(tasks), 1)
}

func Test_InMemoryTaskRepositoryNextId(t *testing.T) {
	t.Parallel()

	repo := repository.InitInMemoryTaskRepository()

	util.AssertEqual(t)(repo.NextId(), 1)
}

func Test_InMemoryTaskRepositoryFindBy(t *testing.T) {
	tests := []struct {
		name        string
		data        entity.Task
		param       int
		expectError bool
		error       error
	}{
		{
			name:        "returns a task",
			data:        entity.Task{Id: 1, Name: "買早餐", Status: 0},
			param:       1,
			expectError: false,
		},
		{
			name:        "returns error when not found",
			data:        entity.Task{Id: 1, Name: "買早餐", Status: 0},
			param:       2,
			expectError: true,
			error:       repository.ErrorNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.InitInMemoryTaskRepository()
			repo.Save(&tc.data)

			got, err := repo.FindBy(tc.param)

			if tc.expectError {
				util.AssertErrorEqual(t)(err, tc.error)
			} else {
				util.AssertEqual(t)(*got, tc.data)
			}
		})
	}
}

func Test_InMemoryTaskRepositoryUpdate(t *testing.T) {
	tests := []struct {
		name        string
		data        entity.Task
		param       entity.Task
		expectError bool
		error       error
	}{
		{
			name:        "returns updated task",
			data:        entity.Task{Id: 1, Name: "買早餐", Status: 0},
			param:       entity.Task{Id: 1, Name: "買晚餐", Status: 1},
			expectError: false,
		},
		{
			name:        "returns error when not found",
			data:        entity.Task{Id: 1, Name: "買早餐", Status: 0},
			param:       entity.Task{Id: 2, Name: "買晚餐", Status: 1},
			expectError: true,
			error:       repository.ErrorNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.InitInMemoryTaskRepository()
			repo.Save(&tc.data)

			got, err := repo.Update(&tc.param)

			if tc.expectError {
				util.AssertErrorEqual(t)(err, tc.error)
			} else {
				util.AssertNotEqual(t)(*got, tc.data)
			}
		})
	}
}

func Test_InMemoryTaskRepositoryDelete(t *testing.T) {
	tests := []struct {
		name  string
		data  entity.Task
		param entity.Task
		error error
		state []*entity.Task
	}{
		{
			name:  "deletes the task",
			data:  entity.Task{Id: 1, Name: "買早餐", Status: 0},
			param: entity.Task{Id: 1, Name: "買早餐", Status: 0},
			state: []*entity.Task{},
		},
		{
			name:  "returns error",
			data:  entity.Task{Id: 1, Name: "買早餐", Status: 0},
			param: entity.Task{Id: 2, Name: "買早餐", Status: 0},
			error: repository.ErrorNotFound,
			state: []*entity.Task{{Id: 1, Name: "買早餐", Status: 0}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository.InitInMemoryTaskRepository()
			repo.Save(&tc.data)

			err := repo.Delete(&tc.param)
			tasks := repo.ListAll()

			util.AssertErrorEqual(t)(err, tc.error)
			util.AssertEqual(t)(len(tasks), len(tc.state))
		})
	}
}
