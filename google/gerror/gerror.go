package gerror

import (
	"context"

	"cloud.google.com/go/errorreporting"
	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/google/gaccount"
	"github.com/piyuo/libsrv/log/logger"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

// NewClient return google error reporting client
//
//	client, err := NewClient(ctx)
//
func NewClient(ctx context.Context) (*errorreporting.Client, error) {
	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}

	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName: env.AppName,
			OnError: func(err error) {
				logger.Error(ctx, errors.Wrap(err, "error report write").Error())
			},
		},
		option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// close error reporting client
//
//	defer close(ctx,client)
//
func close(ctx context.Context, client *errorreporting.Client) {
	if err := client.Close(); err != nil {
		logger.Error(ctx, errors.Wrap(err, "error report close").Error())
	}
}

// Write error to google cloud
//
//	Write(ctx, "nil pointer",  stack)
//
func Write(ctx context.Context, message, stack string) {
	if err := write(ctx, message, stack); err != nil {
		logger.Error(ctx, errors.Wrap(err, "error report create").Error())
	}
}

// write error to google cloud
//
//	write(ctx, "nil pointer", stack)
//
func write(ctx context.Context, message, stack string) error {
	client, err := NewClient(ctx)
	if err != nil {
		return err
	}
	defer close(ctx, client)

	//	stack format like
	//	at firstLine (a.js:3)
	//	at secondLine (b.js:3)

	e := errors.New(message)
	user := env.GetUserID(ctx)

	if stack == "" {
		client.Report(errorreporting.Entry{
			Error: e, User: user,
			Req: env.GetRequest(ctx),
		})
		return nil
	}

	client.Report(errorreporting.Entry{
		Error: e, User: user,
		Stack: []byte(stack),
		Req:   env.GetRequest(ctx),
	})
	return nil
}
