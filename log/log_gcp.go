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

func createErrorClient(ctx context.Context) (*errorreporting.Client, error) {
	cred, err := gcp.LogCredential(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get google credential, check /keys/log.key exist")
	}

	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName: app.PiyuoID(),
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
//	Log(ctx,"hello",WARNING,true)
func logToGcp(ctx context.Context, message, application, identity string, level int32, fromClient bool) {
	if message == "" {
		return
	}
	head := logHeadFromAI(application, identity, fromClient)
	fmt.Printf("%v%v (logged)\n", head, message)
	severity := logging.Notice
	switch level {
	case WARNING:
		severity = logging.Warning
	case CRITICAL:
		severity = logging.Critical
	case ALERT:
		severity = logging.Alert
	case EMERGENCY:
		severity = logging.Emergency
	}

	client, err := createLogClient(ctx)
	if err != nil {
		Error(ctx, errors.Wrap(err, "failed to create log client"), nil)
		return
	}

	file := client.Logger(app.PiyuoID())
	if err != nil {
		Error(ctx, errors.Wrap(err, "failed to create log file"), nil)
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
		},
	}
	if identity != "" {
		entry.Labels["identity"] = identity
	}

	file.Log(entry)

	if err := client.Close(); err != nil {
		Error(ctx, errors.Wrap(err, "failed to close client"), nil)
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
func errorToGcp(ctx context.Context, message, application, identity, stack, errID string, fromClient bool, r *http.Request) {
	head := logHeadFromAI(application, identity, fromClient)
	client, err := createErrorClient(ctx)
	if err != nil {
		fmt.Printf("[not logged]: failed to create error client\n%v\n", err)
		return
	}
	defer client.Close()

	e := errors.New(head + message + " (" + errID + ")")
	if stack == "" {
		client.Report(errorreporting.Entry{
			Error: e, User: identity,
		})

	} else {
		client.Report(errorreporting.Entry{
			Error: e, User: identity,
			Stack: []byte(stack),
			Req:   r,
		})

	}
}
