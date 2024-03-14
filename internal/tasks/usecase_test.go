package tasks_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks"
	util "github.com/dannyh79/whostodo/internal/testutil"
)

func Test_ListTasks(t *testing.T) {
	tests := []struct {
		name     string
		data     []repository.TaskSchema
		expected []*tasks.TaskOutput
	}{
		{
			name:     "returns tasks",
			data:     []repository.TaskSchema{{Id: 1, Name: "name", Status: 0}},
			expected: []*tasks.TaskOutput{{Id: 1, Name: "name", Status: 0}},
		},
		{
			name:     "returns empty tasks",
			expected: []*tasks.TaskOutput{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := util.InitMockTaskRepository()
			if data := tc.data; len(data) > 0 {
				for _, row := range data {
					repo.PopulateData(row)
				}
			}
			usecase := tasks.InitTasksUsecase(repo)
			got := usecase.ListTasks()

			util.AssertEqual(t)(got, tc.expected)
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
			repo := util.InitMockTaskRepository()
			usecase := tasks.InitTasksUsecase(repo)
			got := usecase.CreateTask(&tc.data)

			util.AssertEqual(t)(*got, tc.expected)
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
			error:       util.MockNotFoundError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := util.InitMockTaskRepository()
			repo.PopulateData(tc.data)
			usecase := tasks.InitTasksUsecase(repo)
			got, err := usecase.UpdateTask(tc.param, &tc.payload)

			if tc.expectError {
				util.AssertErrorEqual(t)(err, tc.error)
			} else {
				if err != nil {
					t.Error(err)
				}
				util.AssertEqual(t)(*got, tc.expected)
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
			error:       util.MockNotFoundError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := util.InitMockTaskRepository()
			repo.PopulateData(tc.data)
			usecase := tasks.InitTasksUsecase(repo)
			err := usecase.DeleteTask(tc.param)

			if tc.expectError {
				util.AssertErrorEqual(t)(err, tc.error)
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
