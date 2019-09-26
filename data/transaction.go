package data

// ITransaction is query interface
type ITransaction interface {
	Get(obj IObject) error
	Put(obj IObject) error
	Delete(obj IObject) error
}

// Transaction represent db transaction
type Transaction struct {
}
