package repository

import "errors"

type Repository[T any] interface {
	ListAll() []*T
	Save(*T) T
	FindBy(id any) (*T, error)
	Update(*T) (*T, error)
	Delete(*T) error
}

var ErrorNotFound = errors.New("Task not found")
