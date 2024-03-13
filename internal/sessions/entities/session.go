package entity

import "time"

type Session struct {
	Id        string
	CreatedAt time.Time
}

func NewSession() *Session {
	return &Session{
		Id:        nextId(),
		CreatedAt: time.Now(),
	}
}

func nextId() string {
	panic("TODO")
}
