package tasks_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/google/go-cmp/cmp"
)

type MockTaskRepository struct {
	data map[int]repository.TaskSchema
}

func (r *MockTaskRepository) Save(t *entity.Task) entity.Task {
	t.Id = len(r.data) + 1
	return *t
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
				t.Errorf(cmp.Diff(got, tc.expected))
			}
		})
	}
}

func Test_UpdateTask(t *testing.T) {
	tests := []struct {
		name     string
		data     repository.TaskSchema
		payload  tasks.UpdateTaskInput
		expected tasks.TaskOutput
	}{
		{
			name:     "returns updated task",
			data:     repository.TaskSchema{Id: 1, Name: "買早餐", Status: 0},
			payload:  tasks.UpdateTaskInput{Id: 1, Name: "買晚餐", Status: 1},
			expected: tasks.TaskOutput{Id: 1, Name: "買晚餐", Status: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := initMockTaskRepository()
			usecase := tasks.InitTasksUsecase(repo)
			got := usecase.UpdateTask(tc.data.Id, &tc.payload)
			if !cmp.Equal(*got, tc.expected) {
				t.Errorf(cmp.Diff(got, tc.expected))
			}
		})
	}
}
