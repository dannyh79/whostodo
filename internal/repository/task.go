package repository

import (
	"errors"

	"github.com/dannyh79/whostodo/internal/tasks/entities"
)

type TaskSchema struct {
	Id     int
	Name   string
	Status int
}

type InMemoryTaskRepository struct {
	position int
	data     map[int]TaskSchema
}

var ErrorNotFound = errors.New("Task not found")

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
	r.data[row.Id] = row
	return *toTask(row)
}

func (r *InMemoryTaskRepository) FindBy(id int) (*entity.Task, error) {
	if (id >= len(r.data)) {
		return &entity.Task{}, ErrorNotFound
	}

	row := r.data[id]
	return toTask(row), nil
}

func (r *InMemoryTaskRepository) Update(t *entity.Task) (*entity.Task, error) {
	if (t.Id >= len(r.data)) {
		return &entity.Task{}, ErrorNotFound
	}

	r.data[t.Id] = *toSchema(t)
	return toTask(r.data[t.Id]), nil
}

func InitInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		data: map[int]TaskSchema{},
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
