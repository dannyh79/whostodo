package repository

type Repository[T any] interface {
	ListAll() []*T
	Save(*T) T
}
