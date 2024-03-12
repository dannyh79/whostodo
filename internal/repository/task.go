package repository

import entity "github.com/dannyh79/whostodo/internal/tasks/entities"

type TaskSchema struct {
	Id int
	Name string
	Status int
}

type InMemoryTaskRepository struct {
	data []TaskSchema
}

func (r *InMemoryTaskRepository) ListAll() []entity.Task {
	var tasks []entity.Task
	for _, row := range r.data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func InitInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		data: []TaskSchema{},
	}
}
