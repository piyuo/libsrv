package libsrv

import (
	"context"
	"fmt"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

//Logger interface
// server: [piyuo-m-us-sys] store-user: hello
// client: <piyuo-m-us-web-page> store-user: hello

//Log level
const (
	NOTICE    int32 = 0 //Normal but significant events, such as start up, shut down, or a configuration change.
	WARNING   int32 = 1 //Warning events might cause problems.
	CRITICAL  int32 = 2 //Critical events cause more severe problems or outages.
	ALERT     int32 = 3 //A person must take an action immediately.
	EMERGENCY int32 = 4 //One or more systems are unusable.
)

func createLogClient(ctx context.Context) (*logging.Client, error) {
	cred, err := EnvGoogleCredential(ctx, LOG)
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
	cred, err := EnvGoogleCredential(ctx, LOG)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get google credential, check /keys/log.key exist")
	}

	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName: EnvPiyuoApp(),
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

// GetLogHead use UPPER case for server, lower for client app
// piyuo id should start with P on server and p on client
// server: [piyuo-m-us-sys] store-user: hello
// client: <piyuo-m-us-web-page> store-user: hello
func generateLogHead(ctx context.Context, fromClient bool) (string, string, error) {
	application := EnvPiyuoApp()
	identity := ""
	identityText := ""
	token, err := TokenFromContext(ctx)
	if err == nil {
		identity = token.Identity()
		if identity != "" {
			identityText = " " + identity
		}
	}

	if fromClient {
		return fmt.Sprintf("<%v>%v: ", application, identityText), identity, nil
	}
	return fmt.Sprintf("[%v]%v: ", application, identityText), identity, nil
}

//LogInfo as Routine information, such as ongoing status or performance.
//
//	LogInfo(ctx,"hello")
func LogInfo(ctx context.Context, message string) {
	head, _, _ := generateLogHead(ctx, false)
	fmt.Printf("%v%v\n", head, message)
}

//LogNotice as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	LogNotice(ctx,"hello")
func LogNotice(ctx context.Context, message string) {
	Log(ctx, message, NOTICE, false)
}

//LogWarning as Warning events might cause problems.
//
//	LogWarning(ctx,"hello")
func LogWarning(ctx context.Context, message string) {
	Log(ctx, message, WARNING, false)
}

//LogCritical as Critical events cause more severe problems or outages.
//
//	LogCritical(ctx,"hello")
func LogCritical(ctx context.Context, message string) {
	Log(ctx, message, CRITICAL, false)
}

//LogAlert as A person must take an action immediately.
//
//	LogAlert(ctx,"hello")
func LogAlert(ctx context.Context, message string) {
	Log(ctx, message, ALERT, false)
}

//LogEmergency as One or more systems are unusable.
//
//	LogEmergency(ctx,"hello")
func LogEmergency(ctx context.Context, message string) {
	Log(ctx, message, EMERGENCY, false)
}

//Log custom message and level to server
//
//	Log(ctx,"hello",WARNING,true)
func Log(ctx context.Context, message string, level int32, fromClient bool) {
	if message == "" {
		return
	}
	head, identity, _ := generateLogHead(ctx, fromClient)
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
		Error(ctx, errors.Wrap(err, "failed to create log client"))
		return
	}

	file := client.Logger(EnvPiyuoApp())
	if err != nil {
		Error(ctx, errors.Wrap(err, "failed to create log file"))
		return
	}

	entry := logging.Entry{
		Payload: message,
		Resource: &mrpb.MonitoredResource{
			Type: "project",
		},
		Severity: severity,
		Labels: map[string]string{
			"application": EnvPiyuoApp(),
			"identity":    identity,
		},
	}
	file.Log(entry)

	if err := client.Close(); err != nil {
		Error(ctx, errors.Wrap(err, "failed to close client"))
		return
	}
}

//LogError log error to google cloud
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	LogError(ctx, message, stack, id, true)
func LogError(ctx context.Context, message, stack, errID string, fromClient bool) {
	head, identity, _ := generateLogHead(ctx, fromClient)
	client, err := createErrorClient(ctx)
	if err != nil {
		fmt.Printf("failed to create error client\n%v\n", err)
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
		})

	}
}

//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
func Error(ctx context.Context, err error) string {
	errID := UUID()
	head, _, _ := generateLogHead(ctx, false)
	fmt.Printf("%v%v (%v)\n", head, err, errID)
	message := err.Error()
	LogError(ctx, message, "", errID, false)
	return errID
}
