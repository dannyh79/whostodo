package testutil_test

import (
	"errors"
	"time"

	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/tasks/entities"
)

type TaskSchema = repository.TaskSchema

type Task = entity.Task

type MockTaskRepository struct {
	Data map[int]TaskSchema
}

var MockNotFoundError = errors.New("not found")

func (r *MockTaskRepository) FindBy(id any) (*Task, error) {
	row, ok := r.Data[id.(int)]
	if !ok {
		return nil, MockNotFoundError
	}
	return &Task{Id: row.Id, Name: row.Name, Status: row.Status}, nil
}

func (r *MockTaskRepository) Update(t *Task) (*Task, error) {
	r.Data[t.Id] = TaskSchema{Id: t.Id, Name: t.Name, Status: t.Status}
	return &Task{Id: t.Id, Name: t.Name, Status: t.Status}, nil
}

func (r *MockTaskRepository) Save(t *Task) Task {
	t.Id = len(r.Data) + 1
	return *t
}

func (r *MockTaskRepository) Delete(t *Task) error {
	return nil
}

func (r *MockTaskRepository) ListAll() []*Task {
	var tasks []*entity.Task
	for _, row := range r.Data {
		tasks = append(tasks, entity.NewTask(row.Id, row.Name, row.Status))
	}
	return tasks
}

func (r *MockTaskRepository) PopulateData(row TaskSchema) {
	r.Data[row.Id] = row
}

type Session = repository.Session

type MockSessionsRepository struct {
	Data map[string]Session
}

func InitMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		Data: make(map[int]repository.TaskSchema),
	}
}

func (r *MockSessionsRepository) FindBy(id any) (*Session, error) {
	if id == nil {
		return nil, MockNotFoundError
	}

	row, ok := r.Data[id.(string)]
	if !ok {
		return nil, MockNotFoundError
	}

	return &Session{Id: row.Id, CreatedAt: row.CreatedAt}, nil
}

func (r *MockSessionsRepository) Update(s *Session) (*Session, error) {
	panic("not implemented")
}

func (r *MockSessionsRepository) Save(s *Session) Session {
	r.Data[s.Id] = *s
	return *s
}

func (r *MockSessionsRepository) Delete(t *Session) error {
	panic("not implemented")
}

func (r *MockSessionsRepository) ListAll() []*Session {
	panic("not implemented")
}

func (r *MockSessionsRepository) PopulateData(row Session) {
	r.Data[row.Id] = row
}

func InitMockSessionsRepository() *MockSessionsRepository {
	return &MockSessionsRepository{
		Data: make(map[string]Session),
	}
}

func NewSession() Session {
	return newStubSession("stubbed_token", time.Now())
}

func NewExpiredSession() Session {
	oneMinuteAgo := time.Now().Add(-(time.Minute + time.Second))
	return newStubSession("stubbed_token", oneMinuteAgo)
}

func newStubSession(id string, createdAt time.Time) Session {
	return Session{Id: id, CreatedAt: createdAt}
}
