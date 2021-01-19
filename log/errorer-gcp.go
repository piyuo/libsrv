package log

import (
	"context"
	"fmt"

	"cloud.google.com/go/errorreporting"
	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/gcp"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type gcpErrorer struct {
	Errorer
	// client is google cloud platform logging client
	//
	client *errorreporting.Client
}

// NewGCPErrorer return error errorer implement by google cloud platform
//
//	ctx := context.Background()
//	logClient, _ := NewGCPLogger(ctx)
//
func NewGCPErrorer(ctx context.Context) (Errorer, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}

	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName: appName,
			OnError: func(err error) {
				fmt.Printf("failed to write error: %v\n", err)
			},
		},
		option.WithCredentials(cred))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stackdriver error client")
	}

	return &gcpErrorer{
		client: client,
	}, nil
}

// Close errorer
//
//	defer errorer.Close()
//
func (c *gcpErrorer) Close() {
	if err := c.client.Close(); err != nil {
		fmt.Printf("!!! %v\n", err)
	}
}

// Write error to google cloud
//
//	stack format like
//
//	at firstLine (a.js:3)
//
//	at secondLine (b.js:3)
//
//	Write(ctx, "nil pointer", "app", stack, "AAA")
//
func (c *gcpErrorer) Write(ctx context.Context, where, message, stack, errID string) {
	header, id := getHeader(ctx, where)
	if shouldPrint {
		fmt.Printf("%v%v (%v)\n%v\n", header, message, errID, stack)
	}

	e := errors.New(header + message + " (" + errID + ")")
	if stack == "" {
		c.client.Report(errorreporting.Entry{
			Error: e, User: id,
			Req: env.GetRequest(ctx),
		})
		return
	}
	c.client.Report(errorreporting.Entry{
		Error: e, User: id,
		Stack: []byte(stack),
		Req:   env.GetRequest(ctx),
	})
}
