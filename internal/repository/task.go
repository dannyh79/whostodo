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

func (r *InMemoryTaskRepository) ListAll() []*entity.Task {
	var tasks []*entity.Task
	for _, row := range r.data {
		tasks = append(tasks, toTask(row))
	}
	return tasks
}

func InitInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		data: []TaskSchema{},
	}
}

func toTask(row TaskSchema) *entity.Task {
	return entity.NewTask(row.Id, row.Name, row.Status)
}
