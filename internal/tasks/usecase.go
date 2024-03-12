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
		output = append(output, *toTaskOutput(task))
	}

	return output
}

func (u *TasksUsecase) CreateTask(i *CreateTaskInput) *TaskOutput {
	task := u.repo.Save(&entity.Task{Name: i.Name})
	return toTaskOutput(&task)
}

func (u *TasksUsecase) UpdateTask(id int, i *UpdateTaskInput) (*TaskOutput, error) {
	task, err := u.repo.FindBy(id)
	if err != nil {
		return nil, err
	}

	task.Name = i.Name
	task.Status = i.Status

	updated, err := u.repo.Update(task)
	if err != nil {
		return nil, err
	}

	return toTaskOutput(updated), nil
}

func InitTasksUsecase(repo TaskRepository) *TasksUsecase {
	return &TasksUsecase{repo}
}

func toTaskOutput(t *entity.Task) *TaskOutput {
	return &TaskOutput{
		Id:     t.Id,
		Name:   t.Name,
		Status: t.Status,
	}
}
