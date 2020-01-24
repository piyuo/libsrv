package protocol

//Transaction is query interface
type Transaction interface {
	Get(obj Object) error
	Put(obj Object) error
	Delete(obj Object) error
}
