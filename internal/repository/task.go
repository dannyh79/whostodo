package repository

import entity "github.com/dannyh79/whostodo/internal/tasks/entities"

type TaskSchema struct {
	Id     int
	Name   string
	Status int
}

type InMemoryTaskRepository struct {
	position int
	data []TaskSchema
}

func (r *InMemoryTaskRepository) ListAll() []*entity.Task {
	var tasks []*entity.Task
	for _, row := range r.data {
		tasks = append(tasks, toTask(row))
	}
	return tasks
}

func (r *InMemoryTaskRepository) NextId() int {
	r.position += 1
	return r.position
}

func (r *InMemoryTaskRepository) Save(t *entity.Task) entity.Task {
	t.Id = r.NextId()
	row := *toSchema(t)
	r.data = append(r.data, row)
	return *toTask(row)
}

func InitInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		data: []TaskSchema{},
	}
}

func toTask(row TaskSchema) *entity.Task {
	return entity.NewTask(row.Id, row.Name, row.Status)
}

func toSchema(t *entity.Task) *TaskSchema {
	return &TaskSchema{
		Id:     t.Id,
		Name:   t.Name,
		Status: t.Status,
	}
}
