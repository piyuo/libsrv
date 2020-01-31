package log

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	app "github.com/piyuo/go-libsrv/app"
	gcp "github.com/piyuo/go-libsrv/secure/gcp"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

//createLogClient return stackdriver log client using credential from log-gcp.key
//
//	ctx := context.Background()
//	logClient, _ := createLogClient(ctx)
func createLogClient(ctx context.Context) (*logging.Client, error) {
	cred, err := gcp.LogCredential(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get google credential, check /keys/log.key exist")
	}

	client, err := logging.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stackdriver logging client")
	}
	return client, nil
}

//createErrorClient return stackdriver error client using credential from log-gcp.key
//
//	ctx := context.Background()
//	errClient, _ := createErrorClient(ctx)
func createErrorClient(ctx context.Context, serviceName, serviceVersion string) (*errorreporting.Client, error) {
	cred, err := gcp.LogCredential(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get google credential, check /keys/log.key exist")
	}

	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName:    serviceName,
			ServiceVersion: serviceVersion,
			OnError: func(err error) {
				fmt.Printf("failed to config error reporting %v\n", err)
			},
		},
		option.WithCredentials(cred))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stackdriver error client")
	}
	return client, nil
}

//Log custom message and level to google cloud platform
//
//	const HERE = "log_gcp"
//	logToGcp(ctx,"my error","piyuo-t-sys",'"user-store",HERE,WARNING)
func logToGcp(ctx context.Context, message, application, identity, where string, level int32) {
	if message == "" {
		return
	}
	h := head(application, identity, where)
	fmt.Printf("%v%v (logged)\n", h, message)
	severity := logging.Info
	switch level {
	case warning:
		severity = logging.Warning
	case critical:
		severity = logging.Critical
	}

	client, err := createLogClient(ctx)
	if err != nil {
		Error(ctx, where, errors.Wrap(err, "failed to create log client"), nil)
		return
	}

	file := client.Logger(app.PiyuoID())
	if err != nil {
		Error(ctx, where, errors.Wrap(err, "failed to create log file"), nil)
		return
	}

	entry := logging.Entry{
		Payload: message,
		Resource: &mrpb.MonitoredResource{
			Type: "project",
		},
		Severity: severity,
		Labels: map[string]string{
			"application": app.PiyuoID(),
			"where":       where,
		},
	}
	if identity != "" {
		entry.Labels["identity"] = identity
	}

	file.Log(entry)

	if err := client.Close(); err != nil {
		Error(ctx, where, errors.Wrap(err, "failed to close client"), nil)
		return
	}
}

//errorToGcp log error to google cloud
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	LogError(ctx, message, stack, id, true)
func errorToGcp(ctx context.Context, message, application, identity, where, stack, errID string, r *http.Request) {
	h := head(application, identity, where)
	client, err := createErrorClient(ctx, application, where)
	if err != nil {
		fmt.Printf("[not logged]: failed to create error client\n%v\n", err)
		return
	}
	defer client.Close()

	e := errors.New(h + message + " (" + errID + ")")
	if stack == "" {
		client.Report(errorreporting.Entry{
			Error: e, User: identity,
			Req: r,
		})

	} else {
		client.Report(errorreporting.Entry{
			Error: e, User: identity,
			Stack: []byte(stack),
			Req:   r,
		})

	}
}
