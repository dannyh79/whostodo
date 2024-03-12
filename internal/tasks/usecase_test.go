package tasks_test

import (
	"testing"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/google/go-cmp/cmp"
)

func Test_ListTasks(t *testing.T) {
	tests := []struct {
		name string
		expected []tasks.TaskOutput
	}{
		{
			name: "returns tasks",
			expected: []tasks.TaskOutput{{Id: 1, Name: "name", Status: 0}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := repository.InitInMemoryRepo()
			usecase := tasks.InitTasksUsecase(*repo)
			got := usecase.ListTasks()
			if !cmp.Equal(got, tc.expected) {
				t.Errorf(cmp.Diff(got, tc.expected))
			}
		})
	}
}
