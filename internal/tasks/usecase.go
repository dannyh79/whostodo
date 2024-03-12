package tasks

import (
	"github.com/dannyh79/whostodo/internal/repository"
	entity "github.com/dannyh79/whostodo/internal/tasks/entities"
)

type TaskOutput struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type CreateTaskInput struct {
	Name string `json:"name"`
}

type UpdateTaskInput struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
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
			Id:     task.Id,
			Name:   task.Name,
			Status: task.Status,
		})
	}

	return output
}

func (u *TasksUsecase) CreateTask(i *CreateTaskInput) *TaskOutput {
	task := u.repo.Save(&entity.Task{Name: i.Name})
	return &TaskOutput{
		Id:     task.Id,
		Name:   task.Name,
		Status: task.Status,
	}
}

func (u *TasksUsecase) UpdateTask(id int, i *UpdateTaskInput) (*TaskOutput, error) {
	return &TaskOutput{
		Id:     1,
		Name:   "買晚餐",
		Status: 1,
	}, nil
}

func InitTasksUsecase(repo TaskRepository) *TasksUsecase {
	return &TasksUsecase{repo}
}
