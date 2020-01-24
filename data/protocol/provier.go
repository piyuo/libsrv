package protocol

import "context"

// Provider provide
type Provider interface {
	Initialize(ctx context.Context)
	NewDB() (DB, error)
}

// provider simplify datastore create
type provider struct {
}

//Initialize db provider
func (p *provider) Initialize(ctx context.Context) error {
	panic("must implement Load() function")
}
