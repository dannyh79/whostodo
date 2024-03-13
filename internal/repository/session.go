package repository

import (
	"time"
)

type Session struct {
	Id        string
	CreatedAt time.Time
}

type SessionSchema struct {
	Id        string
	CreatedAt time.Time
}

type InMemorySessionRepository struct {
	data map[string]SessionSchema
}

func (r *InMemorySessionRepository) Delete(*Session) error {
	panic("not implemented")
}

func (r *InMemorySessionRepository) FindBy(id any) (*Session, error) {
	if id == nil {
		return nil, ErrorNotFound
	}

	row, ok := r.data[id.(string)]
	if !ok {
		return nil, ErrorNotFound
	}

	return toSession(row), nil
}

func (r *InMemorySessionRepository) ListAll() []*Session {
	panic("not implemented")
}

func (r *InMemorySessionRepository) Save(s *Session) Session {
	r.data[s.Id] = *toSessionSchema(s)
	return *s
}

func (r *InMemorySessionRepository) Update(*Session) (*Session, error) {
	panic("not implemented")
}

func InitInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{
		data: map[string]SessionSchema{},
	}
}

func toSession(s SessionSchema) *Session {
	return &Session{
		Id:        s.Id,
		CreatedAt: s.CreatedAt,
	}
}

func toSessionSchema(s *Session) *SessionSchema {
	return &SessionSchema{
		Id:        s.Id,
		CreatedAt: s.CreatedAt,
	}
}
