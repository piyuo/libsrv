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

//gcpCreateLogClient return stackdriver log client using credential from log-gcp.key
//
//	ctx := context.Background()
//	logClient, _ := gcpCreateLogClient(ctx)
func gcpCreateLogClient(ctx context.Context) (*logging.Client, error) {
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

//gcpCreateErrorClient return stackdriver error client using credential from log-gcp.key
//
//	ctx := context.Background()
//	errClient, _ := gcpCreateErrorClient(ctx)
func gcpCreateErrorClient(ctx context.Context, serviceName, serviceVersion string) (*errorreporting.Client, error) {
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

//gcpLog log message and level to server
//
//	HERE := "log_test"
//	gcpLog(ctx, "hello", "piyuo-m-us-sys", "user-store",HERE, WARNING)
func gcpLog(ctx context.Context, message, application, identity, where string, level int32) {
	client, err := gcpCreateLogClient(ctx)
	if err != nil {
		Error(ctx, where, err, nil)
		return
	}
	logger := client.Logger(app.PiyuoID())
	gcpWriteByLogger(ctx, logger, message, application, identity, where, level)
	if err := client.Close(); err != nil {
		Error(ctx, where, errors.Wrap(err, "failed to close client"), nil)
		return
	}
}

//gcpWriteByLogger custom message and level to google cloud platform
//
//	const HERE = "log_gcp"
//	gcpWriteByLogger(ctx,logger, "my error","piyuo-t-sys",'"user-store",HERE,WARNING)
func gcpWriteByLogger(ctx context.Context, logger *logging.Logger, message, application, identity, where string, level int32) {
	if message == "" {
		return
	}
	h := head(application, identity, where)
	fmt.Printf("%v%v (logged)\n", h, message)
	severity := logging.Info
	switch level {
	case warning:
		severity = logging.Warning
	case alert:
		severity = logging.Critical
	}

	entry := logging.Entry{
		Payload: h + message,
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
	logger.Log(entry)
}

//gcpError log error and stack to google cloud
func gcpError(ctx context.Context, message, application, identity, where, stack, errID string, r *http.Request) {

	client, err := gcpCreateErrorClient(ctx, application, where)
	if err != nil {
		fmt.Printf("!!! %v\n", err)
		return
	}
	defer client.Close()

	gcpErrorByClient(ctx, client, message, application, identity, where, stack, errID, r)
}

//gcpError log error to google cloud
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	gcpErrorByClient(ctx, message,application,identity,where, stack, id, request)
func gcpErrorByClient(ctx context.Context, client *errorreporting.Client, message, application, identity, where, stack, errID string, r *http.Request) {
	h := head(application, identity, where)

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
