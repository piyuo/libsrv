package data

import (
	"context"

	firebase "firebase.google.com/go"
	app "github.com/piyuo/go-libsrv/app"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

// ProviderFirestore implement google firestore
type ProviderFirestore struct {
	Provider
	//credentials *google.Credentials
	app *firebase.App
	ctx context.Context
}

// NewProviderFirestore provide new Provider for google firestore
func NewProviderFirestore() *ProviderFirestore {
	return &ProviderFirestore{}
}

//Initialize check env variable DATA_CRED to init google credentials for firestore
func (provider *ProviderFirestore) Initialize() {
	ctx := context.Background()
	cred, err := app.GetGoogleCloudCredential(app.DB)
	if err != nil {
		libsrv.LogAlert(ctx, "database operation failed to get google credential.  %v")
		return
	}

	provider.app, err = firebase.NewApp(ctx, nil, option.WithCredentials(cred))
	if err != nil {
		libsrv.Sys().Emergency("failed to create firebase client")
		panic(err)
	}
	provider.ctx = ctx
}

//NewDB create db instance
func (provider *ProviderFirestore) NewDB() (IDB, error) {
	client, err := provider.app.Firestore(provider.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}
	return NewDBFirestore(client), nil
}
