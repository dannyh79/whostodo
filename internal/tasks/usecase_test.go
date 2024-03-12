package tasks_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
	"github.com/google/go-cmp/cmp"
)

type MockTaskRepository struct {
	data []repository.TaskSchema
}

func (r *MockTaskRepository) ListAll() []*entity.Task {
	var tasks []*entity.Task
	for _, row := range r.data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func (r *MockTaskRepository) PopulateData(data []repository.TaskSchema) {
	for _, row := range(data) {
		r.data = append(r.data, row)
	}
}

func initMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		data: []repository.TaskSchema{},
	}
}

func Test_ListTasks(t *testing.T) {
	tests := []struct {
		name string
		data []repository.TaskSchema
		expected []tasks.TaskOutput
	}{
		{
			name: "returns tasks",
			data: []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			expected: []tasks.TaskOutput{{Id: 1, Name: "name", Status: 0}},
		},
		{
			name: "returns empty tasks",
			expected: []tasks.TaskOutput{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := initMockTaskRepository()
			if (len(tc.data) > 0) {
				repo.PopulateData(tc.data)
			}
			usecase := tasks.InitTasksUsecase(repo)
			got := usecase.ListTasks()
			if !cmp.Equal(got, tc.expected) {
				t.Errorf(cmp.Diff(got, tc.expected))
			}
		})
	}
}
