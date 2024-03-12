package tasks

import (
	"github.com/dannyh79/whostodo/internal/repository"
	entity "github.com/dannyh79/whostodo/internal/tasks/entities"
)

type TaskOutput struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Status int `json:"status"`
}

type TaskRepository repository.Repository[entity.Task]

type TasksUsecase struct {
	repo TaskRepository
}

func (u *TasksUsecase) ListTasks() []TaskOutput {
	var output = make([]TaskOutput, 0)

	tasks := u.repo.ListAll()
	for _, task := range tasks {
		output = append(output, TaskOutput{
			Id: task.Id,
			Name: task.Name,
			Status: task.Status,
		})
	}

	output = append(output, TaskOutput{
		Id: 1,
		Name: "name",
		Status: 0,
	})

	return output
}

func InitTasksUsecase(r TaskRepository) *TasksUsecase {
	return &TasksUsecase{repo: r}
}
