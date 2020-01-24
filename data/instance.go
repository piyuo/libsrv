package data

import (
	"context"

	firestroe "github.com/piyuo/go-libsrv/data/fire"
)

var defaultProvider Provider

// NewDB create db from default provider
func NewDB(ctx context.Context) (DB, error) {
	if defaultProvider == nil {
		defaultProvider = firestroe.NewProviderFirestore()
		defaultProvider.Initialize(ctx)
	}
	return defaultProvider.NewDB()
}
