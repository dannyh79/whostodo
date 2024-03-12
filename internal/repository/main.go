package repository

type Repository[T any] interface {
	ListAll() []T
}
