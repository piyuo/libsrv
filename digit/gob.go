package digit

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
)

// encode   293559         3696 ns/op      1504 B/op       22 allocs/op
// decode    72984        14990 ns/op      6928 B/op      179 allocs/op

// Encode any object to bytes
//
//	encoded, err := Encode(sample)
//
func Encode(p interface{}) ([]byte, error) {
	if p == nil {
		return nil, nil
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, errors.Wrap(err, "encode")
	}
	result := buf.Bytes()
	fmt.Printf("encode to %v bytes\n", len(result))
	return result, nil
}

// Decode bytes into object
//
//	err = Decode(encoded, sample2)
//
func Decode(encodedBytes []byte, obj interface{}) error {
	if obj == nil {
		return errors.New("obj must not nil")
	}
	if encodedBytes == nil {
		return nil
	}

	dec := gob.NewDecoder(bytes.NewReader(encodedBytes))
	err := dec.Decode(obj)
	if err != nil {
		return errors.Wrap(err, "decode")
	}
	return nil
}
