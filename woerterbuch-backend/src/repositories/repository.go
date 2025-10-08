package repositories

type Repository[T any] interface {
	GetById(uint64) (*T, error)
	Insert(*T) error
	InsertMultiple([]*T) error
	GetAll() ([]*T, error)
	existsByWordAndTranslation(*T) bool
}
