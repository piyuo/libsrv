package rmdl

import (
	data "github.com/piyuo/libsrv/data"
)

// Location represent single location
//
type Location struct {
	data.DocObject `firestore:"-"`
}
