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

func Test_InMemoryTaskRepositoryFindBy(t *testing.T) {
	t.Run("returns a task", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()
		task := entity.Task{Id: 1, Name: "買早餐", Status: 0}
		repo.Save(&task)

		got, err := repo.FindBy(task.Id)
		if want := &task; !cmp.Equal(want, got) {
			t.Errorf(cmp.Diff(want, got))
		}
		if err != nil {
			t.Errorf("unexpected error %v", err.Error())
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()

		got, err := repo.FindBy(1)
		if got != nil {
			t.Errorf(cmp.Diff(nil, got))
		}
		if want := repository.ErrorNotFound; !errors.Is(err, want) {
			t.Errorf(cmp.Diff(want.Error(), err.Error()))
		}
	})
}

func Test_InMemoryTaskRepositoryUpdate(t *testing.T) {
	t.Run("returns updated task", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()
		task := entity.Task{Id: 1, Name: "買早餐", Status: 0}
		repo.Save(&task)

		got, err := repo.Update(&task)
		if want := &task; !cmp.Equal(want, got) {
			t.Errorf(cmp.Diff(want, got))
		}
		if err != nil {
			t.Errorf("unexpected error %v", err.Error())
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		t.Parallel()

		repo := repository.InitInMemoryTaskRepository()

		task := entity.Task{Id: 1, Name: "買早餐", Status: 0}
		got, err := repo.Update(&task)
		if got != nil {
			t.Errorf(cmp.Diff(nil, got))
		}
		if want := repository.ErrorNotFound; !errors.Is(err, want) {
			t.Errorf(cmp.Diff(want.Error(), err.Error()))
		}
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
			if want := tc.error; !errors.Is(err, want) {
				t.Errorf(cmp.Diff(want.Error(), err.Error()))
			}

			tasks := repo.ListAll()
			if len(tc.state) != len(tasks) {
				t.Errorf(cmp.Diff(tc.state, tasks))
			}
		})
	}
}
