// Code generated by proto watcher. DO NOT EDIT.
package shared

//MapXXX mapping ID to object
type MapXXX struct {
}

// NewObjectByID return new object from id
func (r *MapXXX) NewObjectByID(id uint16) (interface{}) {
	switch id {
	case 1: return new(Err)
	case 2: return new(Num)
	case 4: return new(PingAction)
	case 5: return new(Text)
	}
	return nil
}
