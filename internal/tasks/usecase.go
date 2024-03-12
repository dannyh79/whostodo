package tasks

type TaskOutput struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Status int `json:"status"`
}

func ListTasks() []TaskOutput {
	return []TaskOutput{{Id: 1, Name: "name", Status: 0}}
}
