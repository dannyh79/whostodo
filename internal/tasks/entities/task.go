package entity

type Task struct {
	Id int
	Name string
	Status int
}

func NewTask(id int, name string, status int) *Task {
	return &Task{
		Id: id,
		Name: name,
		Status: status,
	}
}
