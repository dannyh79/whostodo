package tasks

import "github.com/dannyh79/whostodo/internal/repository"

type TaskOutput struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Status int `json:"status"`
}

type TasksUsecase struct {
	repo repository.InMemoryRepo
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

func InitTasksUsecase(r repository.InMemoryRepo) *TasksUsecase {
	return &TasksUsecase{repo: r}
}
