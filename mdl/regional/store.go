package rmdl

import (
	data "github.com/piyuo/libsrv/data"
)

// Store represent single store
//
type Store struct {
	data.DocObject `firestore:"-"`
}
