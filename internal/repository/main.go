package repository

type Repository[T any] interface {
	ListAll() []*T
	Save(*T) T
	FindBy(id int) (*T, error)
	Update(*T) (*T, error)
}
