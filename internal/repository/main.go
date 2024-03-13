package repository

type Repository[T any] interface {
	ListAll() []*T
	Save(*T) T
	FindBy(id any) (*T, error)
	Update(*T) (*T, error)
	Delete(*T) error
}
