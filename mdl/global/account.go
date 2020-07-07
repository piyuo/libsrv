package gmdl

import (
	data "github.com/piyuo/libsrv/data"
)

// Account represent single account
//
type Account struct {
	data.DocObject `firestore:"-"`
}
