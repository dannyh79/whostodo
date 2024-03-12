package repository

type TaskSchema struct {
	Id int
	Name string
	Status int
}

type InMemoryRepo struct {
	data []TaskSchema
}

func (r *InMemoryRepo) ListAll() []TaskSchema {
	var tasks []TaskSchema
	for _, task := range r.data {
		tasks = append(tasks, task)
	}
	return tasks
}

func InitInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		data: []TaskSchema{},
	}
}
