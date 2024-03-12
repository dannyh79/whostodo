package tasks_test

import (
	"errors"
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/google/go-cmp/cmp"
)

type MockTaskRepository struct {
	data map[int]repository.TaskSchema
}

var mockNotFoundError = errors.New("not found")

func (r *MockTaskRepository) FindBy(id int) (*entity.Task, error) {
	row, ok := r.data[id]
	if !ok {
		return nil, mockNotFoundError
	}

	return &entity.Task{Id: row.Id, Name: row.Name, Status: row.Status}, nil
}

func (r *MockTaskRepository) Update(t *entity.Task) (*entity.Task, error) {
	_, ok := r.data[t.Id]
	if !ok {
		return nil, mockNotFoundError
	}

	r.data[t.Id] = repository.TaskSchema{Id: t.Id, Name: t.Name, Status: t.Status}
	return &entity.Task{Id: t.Id, Name: t.Name, Status: t.Status}, nil
}

func (r *MockTaskRepository) Save(t *entity.Task) entity.Task {
	t.Id = len(r.data) + 1
	return *t
}

func (r *MockTaskRepository) Delete(t *entity.Task) error {
	_, ok := r.data[t.Id]
	if !ok {
		return mockNotFoundError
	}

	delete(r.data, t.Id)
	return nil
}

func (r *MockTaskRepository) ListAll() []*entity.Task {
	var tasks []*entity.Task
	for _, row := range r.data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func (r *MockTaskRepository) PopulateData(row repository.TaskSchema) {
	r.data[row.Id] = row
}

func initMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		data: make(map[int]repository.TaskSchema),
	}
}

func Test_ListTasks(t *testing.T) {
	tests := []struct {
		name     string
		data     []repository.TaskSchema
		expected []tasks.TaskOutput
	}{
		{
			name:     "returns tasks",
			data:     []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			expected: []tasks.TaskOutput{{Id: 1, Name: "name", Status: 0}},
		},
		{
			name:     "returns empty tasks",
			expected: []tasks.TaskOutput{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := initMockTaskRepository()
			if data := tc.data; len(data) > 0 {
				for _, row := range data {
					repo.PopulateData(row)
				}
			}
			usecase := tasks.InitTasksUsecase(repo)
			got := usecase.ListTasks()
			if !cmp.Equal(got, tc.expected) {
				t.Errorf(cmp.Diff(got, tc.expected))
			}
		})
	}
}

func Test_CreateTask(t *testing.T) {
	tests := []struct {
		name     string
		data     tasks.CreateTaskInput
		expected tasks.TaskOutput
	}{
		{
			name:     "returns created task",
			data:     tasks.CreateTaskInput{Name: "買晚餐"},
			expected: tasks.TaskOutput{Id: 1, Name: "買晚餐", Status: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := initMockTaskRepository()
			usecase := tasks.InitTasksUsecase(repo)
			got := usecase.CreateTask(&tc.data)
			if !cmp.Equal(*got, tc.expected) {
				t.Errorf(cmp.Diff(tc.expected, *got))
			}
		})
	}
}

func Test_UpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		data        repository.TaskSchema
		param       int
		payload     tasks.UpdateTaskInput
		expected    tasks.TaskOutput
		expectError bool
		error       error
	}{
		{
			name:        "returns updated task",
			data:        repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:       1,
			payload:     tasks.UpdateTaskInput{Name: "買晚餐", Status: 1},
			expected:    tasks.TaskOutput{Id: 1, Name: "買晚餐", Status: 1},
			expectError: false,
		},
		{
			name:        "returns error",
			data:        repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:       2,
			payload:     tasks.UpdateTaskInput{Name: "買晚餐", Status: 1},
			expectError: true,
			error:       mockNotFoundError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := initMockTaskRepository()
			repo.PopulateData(tc.data)
			usecase := tasks.InitTasksUsecase(repo)
			got, err := usecase.UpdateTask(tc.param, &tc.payload)
			if tc.expectError {
				if !errors.Is(err, tc.error) {
					t.Errorf(cmp.Diff(tc.error.Error(), err.Error()))
				}
			} else {
				if err != nil {
					t.Error(err)
				}
				if want := &tc.expected; !cmp.Equal(want, got) {
					t.Errorf(cmp.Diff(want, got))
				}
			}
		})
	}
}

func Test_DeteleTask(t *testing.T) {
	tests := []struct {
		name        string
		data        repository.TaskSchema
		param       int
		expectError bool
		error       error
	}{
		{
			name:        "deletes the task",
			data:        repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:       1,
			expectError: false,
		},
		{
			name:        "returns error",
			data:        repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			param:       2,
			expectError: true,
			error:       mockNotFoundError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := initMockTaskRepository()
			repo.PopulateData(tc.data)
			usecase := tasks.InitTasksUsecase(repo)
			err := usecase.DeleteTask(tc.param)

			if tc.expectError {
				if !errors.Is(err, tc.error) {
					t.Errorf(cmp.Diff(tc.error.Error(), err.Error()))
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
