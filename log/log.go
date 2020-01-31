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
	info    int32 = 1 //Normal but significant events, such as start up, shut down, or a configuration change.
	warning int32 = 2 //Warning events might cause problems.
	alert   int32 = 3 //A person must take an action immediately
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

// head get log head from  application, identity
//
// piyuo-m-us-sys.user-store.where:
//
//	h,identity := head("piyuo-m-us-sys","user-store","where")
func head(application, identity, where string) string {
	return fmt.Sprintf("%v.%v.%v: ", application, identity, where)
}

//Debug as Routine information, such as ongoing status or performance.
//
//	HERE := "log_test"
//	Debug(ctx,HERE,"hello")
func Debug(ctx context.Context, where, message string) {
	application, identity := aiFromContext(ctx)
	h := head(application, identity, where)
	fmt.Printf("%v%v\n", h, message)
}

//Info as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	HERE := "log_test"
//	Info(ctx,HERE,"hi")
func Info(ctx context.Context, where, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, where, info)
}

//Warning as Warning events might cause problems.
//
//	HERE := "log_test"
//	Warning(ctx,HERE,"hi")
func Warning(ctx context.Context, where, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, where, warning)
}

//Alert A person must take an action immediately
//
//	HERE := "log_test"
//	Critical(ctx,HERE,"hi")
func Alert(ctx context.Context, where, message string) {
	application, identity := aiFromContext(ctx)
	CustomLog(ctx, message, application, identity, where, alert)
}

//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	HERE := "log_test"
//	LogErr(ctx,HERE, err)
func Error(ctx context.Context, where string, err error, r *http.Request) string {
	errID := tools.UUID()
	application, identity := aiFromContext(ctx)
	h := head(application, identity, where)
	fmt.Printf("%v%v (%v)\n", h, err, errID)
	message := err.Error()
	CustomError(ctx, message, application, identity, where, "", errID, r)
	return errID
}

//CustomLog log message and level to server
//
//	HERE := "log_test"
//	CustomLog(ctx, "hello", "piyuo-m-us-sys", "user-store",HERE, WARNING)
func CustomLog(ctx context.Context, message, application, identity, where string, level int32) {
	logToGcp(ctx, message, application, identity, where, level)
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
//	HERE := "log_test"
//	LogError(ctx, "hi error", "piyuo-m-us-sys", "user-store",HERE, stack, errID)
func CustomError(ctx context.Context, message, application, identity, where, stack, errID string, r *http.Request) {
	errorToGcp(ctx, message, application, identity, where, stack, errID, r)
}
