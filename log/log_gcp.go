package log

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
func gcpCreateErrorClient(ctx context.Context) (*errorreporting.Client, error) {
	cred, err := gcp.LogCredential(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get google credential, check /keys/log.key exist")
	}
	serviceName := app.PiyuoID()
	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName: serviceName,
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

//gcpLogOpen create logger to write log
//
//	logger, close, err := gcpLogOpen(ctx)
func gcpLogOpen(ctx context.Context) (*logging.Logger, func(), error) {
	client, err := gcpCreateLogClient(ctx)
	if err != nil {
		fmt.Printf("!!! %v\n", err)
		return nil, nil, err
	}
	logger := client.Logger(app.PiyuoID())
	return logger, func() {
		if err := client.Close(); err != nil {
			fmt.Printf("!!! %v\n", err)
			return
		}
	}, nil
}

//gcpLogWrite message and level to google cloud platform
//
//	gcpLogWrite(logger,time.Now(), "my error","piyuo-t-sys",'"user-store","log_gcp",WARNING)
func gcpLogWrite(logger *logging.Logger, logtime time.Time, message, application, identity, where string, level int32) {
	if message == "" {
		return
	}
	h := head(application, identity, where)
	severity := logging.Info
	switch level {
	case LevelWarning:
		severity = logging.Warning
	case LevelAlert:
		severity = logging.Critical
	case LevelDebug:
		severity = logging.Debug
	}

	entry := logging.Entry{
		Timestamp: logtime,
		Payload:   h + message,
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

//gcpErrorOpen log error and stack to google cloud
//
//	client, close, err := gcpErrorOpen(ctx)
func gcpErrorOpen(ctx context.Context) (*errorreporting.Client, func(), error) {
	client, err := gcpCreateErrorClient(ctx)
	if err != nil {
		fmt.Printf("!!! %v\n", err)
		return nil, nil, err
	}
	return client, func() {
		client.Close()
	}, nil
}

//gcpError log error to google cloud
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	gcpErrorWrite(client, message, application, identity, here, stack, id, request)
func gcpErrorWrite(client *errorreporting.Client, message, application, identity, where, stack, errID string, r *http.Request) {
	h := head(application, identity, where)
	e := errors.New(h + message + " (" + errID + ")")
	if stack == "" {
		client.Report(errorreporting.Entry{
			Error: e, User: identity,
			Req: r,
		})
		return
	}
	client.Report(errorreporting.Entry{
		Error: e, User: identity,
		Stack: []byte(stack),
		Req:   r,
	})

}
