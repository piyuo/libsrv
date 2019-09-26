package data

// IProvider provide
type IProvider interface {
	Initialize()
	NewDB() (IDB, error)
}

// Provider simplify datastore create
type Provider struct {
}

//Initialize db provider
func (provider *Provider) Initialize() error {
	panic("must implement Load() function")
}

var instance IProvider

// ProviderInstance get Provider singleton instance
func ProviderInstance() IProvider {
	if instance == nil {
		instance = NewProviderFirestore()
		instance.Initialize()
	}
	return instance
}
