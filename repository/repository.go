package repository

type Repository[TEntity any, TKey stringer] interface {
	findAll() []TEntity
	getById(id TKey) (entity TEntity, err error)
	save(entity TEntity) (error error)
	delete(entity TEntity) (error error)
	deleteById(id TKey) (err error)
}

type stringer interface {
	String() string
}
