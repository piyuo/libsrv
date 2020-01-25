package log

import (
	"context"
	"fmt"
	"net/http"

	app "github.com/piyuo/go-libsrv/app"
	tools "github.com/piyuo/go-libsrv/tools"
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

// aiFromContext get application, identity from context
//
// application: piyuo-m-us-sys
//
// identity: user-store
//
//	application,identity := aiFromContext(ctx)
func aiFromContext(ctx context.Context) (string, string) {
	application := app.PiyuoID()
	identity := ""
	token, err := app.TokenFromContext(ctx)
	if err == nil {
		identity = token.Identity()
	}
	return application, identity
}

// logHeadFromAI get log head from  application, identity
//
// server: [piyuo-m-us-sys] user-store:
//
// client: <piyuo-m-us-web-page> user-store:
//
//	head,identity := logHeadFromAI("piyuo-m-us-sys","user-store")
func logHeadFromAI(application, identity string, fromClient bool) string {
	identityText := ""
	if identity != "" {
		identityText = " " + identity
	}
	if fromClient {
		return fmt.Sprintf("<%v>%v: ", application, identityText)
	}
	return fmt.Sprintf("[%v]%v: ", application, identityText)
}

//Info as Routine information, such as ongoing status or performance.
//
//	Info(ctx,"hello")
func Info(ctx context.Context, message string) {
	application, identity := aiFromContext(ctx)
	head := logHeadFromAI(application, identity, false)
	fmt.Printf("%v%v\n", head, message)
}

//Notice as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	Notice(ctx,"hello")
func Notice(ctx context.Context, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, NOTICE, false)
}

//Warning as Warning events might cause problems.
//
//	Warning(ctx,"hello")
func Warning(ctx context.Context, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, WARNING, false)
}

//Critical as Critical events cause more severe problems or outages.
//
//	Critical(ctx,"hello")
func Critical(ctx context.Context, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, CRITICAL, false)
}

//Alert as A person must take an action immediately.
//
//	Alert(ctx,"hello")
func Alert(ctx context.Context, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, ALERT, false)
}

//Emergency as One or more systems are unusable.
//
//	Emergency(ctx,"hello")
func Emergency(ctx context.Context, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, EMERGENCY, false)
}

//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
func Error(ctx context.Context, err error, r *http.Request) string {
	errID := tools.UUID()
	application, identity := aiFromContext(ctx)
	head := logHeadFromAI(application, identity, false)
	fmt.Printf("%v%v (%v)\n", head, err, errID)
	message := err.Error()
	CustomError(ctx, message, application, identity, "", errID, false, r)
	return errID
}

//CustomLog log message and level to server
//
//	CustomLog(ctx, "hello", "piyuo-m-us-sys", "user-store", WARNING, true)
func CustomLog(ctx context.Context, message, application, identity string, level int32, fromClient bool) {
	logToGcp(ctx, message, application, identity, level, fromClient)
}

//CustomError log error and stack to server
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	errID := tools.UUID()
//	LogError(ctx, "hi error", "piyuo-m-us-sys", "user-store",, stack, errID, true)
func CustomError(ctx context.Context, message, application, identity, stack, errID string, fromClient bool, r *http.Request) {
	errorToGcp(ctx, message, application, identity, stack, errID, fromClient, r)
}
