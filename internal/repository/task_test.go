package repository_test

import (
	"errors"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/google/go-cmp/cmp"
)

func Test_InMemoryTaskRepositoryListAll(t *testing.T) {
	repo := repository.InitInMemoryTaskRepository()

	tasks := repo.ListAll()

	assertEqual(t)(len(tasks), 0)
}

func Test_InMemoryTaskRepositorySave(t *testing.T) {
	repo := repository.InitInMemoryTaskRepository()
	repo.Save(&entity.Task{Id: 1, Name: "買晚餐", Status: 0})

	tasks := repo.ListAll()

	assertEqual(t)(len(tasks), 1)
}

func Test_InMemoryTaskRepositoryNextId(t *testing.T) {
	repo := repository.InitInMemoryTaskRepository()

	assertEqual(t)(repo.NextId(), 1)
}

func Test_InMemoryTaskRepositoryFindBy(t *testing.T) {
	t.Run("returns a task", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()
		task := entity.Task{Id: 1, Name: "買早餐", Status: 0}
		repo.Save(&task)

		got, err := repo.FindBy(task.Id)

		assertEqual(t)(*got, task)
		assertErrorEqual(t)(err, nil)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()

		got, err := repo.FindBy(1)

		if got != nil {
			assertEqual(t)(got, nil)
		}
		assertErrorEqual(t)(err, repository.ErrorNotFound)
	})
}

func Test_InMemoryTaskRepositoryUpdate(t *testing.T) {
	t.Run("returns updated task", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()
		task := entity.Task{Id: 1, Name: "買早餐", Status: 0}
		repo.Save(&task)

		got, err := repo.Update(&task)

		assertEqual(t)(*got, task)
		assertErrorEqual(t)(err, nil)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()
		task := entity.Task{Id: 1, Name: "買早餐", Status: 0}

		got, err := repo.Update(&task)

		if got != nil {
			assertEqual(t)(got, nil)
		}
		assertErrorEqual(t)(err, repository.ErrorNotFound)
	})
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

			assertErrorEqual(t)(err, tc.error)
			assertEqual(t)(len(tasks), len(tc.state))
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

func assertErrorEqual(t *testing.T) func(got error, want error) {
	return func(got error, want error) {
		t.Helper()
		if want != nil && !errors.Is(got, want) {
			t.Errorf(cmp.Diff(got.Error(), want.Error()))
		}
	}
}
