package entity

import (
	"crypto/rand"
	"fmt"
	"time"
)

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

// Warning: Not a safe implementation.
func nextId() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	token := fmt.Sprintf("%x", bytes)
	return token
}
