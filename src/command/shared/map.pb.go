// Code generated by proto watcher. DO NOT EDIT.
package shared

//MapXXX mapping ID to object
type MapXXX struct {
}

// NewObjectByID return new object from id
func (r *MapXXX) NewObjectByID(id uint16) (interface{}) {
	switch id {
	case 1: return new(PbBool)
	case 2: return new(PbError)
	case 3: return new(PbInt)
	case 4: return new(PbOK)
	case 5: return new(PbString)
	}
	return nil
}