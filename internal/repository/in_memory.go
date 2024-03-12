package repository

import entity "github.com/dannyh79/whostodo/internal/tasks/entities"

type TaskSchema struct {
	Id int
	Name string
	Status int
}

type InMemoryRepo struct {
	data []TaskSchema
}

func (r *InMemoryRepo) ListAll() []entity.Task {
	var tasks []entity.Task
	for _, row := range r.data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func InitInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		data: []TaskSchema{},
	}
}
