package log

import (
	"context"
	"fmt"

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

// GetLogHead use UPPER case for server, lower for client app
// piyuo id should start with P on server and p on client
// server: [piyuo-m-us-sys] store-user: hello
// client: <piyuo-m-us-web-page> store-user: hello
func generateLogHead(ctx context.Context, fromClient bool) (string, string, error) {
	application := app.PiyuoID()
	identity := ""
	identityText := ""
	token, err := app.TokenFromContext(ctx)
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

//Info as Routine information, such as ongoing status or performance.
//
//	Info(ctx,"hello")
func Info(ctx context.Context, message string) {
	head, _, _ := generateLogHead(ctx, false)
	fmt.Printf("%v%v\n", head, message)
}

//Notice as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	Notice(ctx,"hello")
func Notice(ctx context.Context, message string) {
	CustomLog(ctx, message, NOTICE, false)
}

//Warning as Warning events might cause problems.
//
//	Warning(ctx,"hello")
func Warning(ctx context.Context, message string) {
	CustomLog(ctx, message, WARNING, false)
}

//Critical as Critical events cause more severe problems or outages.
//
//	Critical(ctx,"hello")
func Critical(ctx context.Context, message string) {
	CustomLog(ctx, message, CRITICAL, false)
}

//Alert as A person must take an action immediately.
//
//	Alert(ctx,"hello")
func Alert(ctx context.Context, message string) {
	CustomLog(ctx, message, ALERT, false)
}

//Emergency as One or more systems are unusable.
//
//	Emergency(ctx,"hello")
func Emergency(ctx context.Context, message string) {
	CustomLog(ctx, message, EMERGENCY, false)
}

//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
func Error(ctx context.Context, err error) string {
	errID := tools.UUID()
	head, _, _ := generateLogHead(ctx, false)
	fmt.Printf("%v%v (%v)\n", head, err, errID)
	message := err.Error()
	CustomError(ctx, message, "", errID, false)
	return errID
}

//CustomLog log message and level to server
//
//	Log(ctx,"hello",WARNING,true)
func CustomLog(ctx context.Context, message string, level int32, fromClient bool) {
	logToGcp(ctx, message, level, fromClient)
}

//CustomError log error to google cloud
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	LogError(ctx, message, stack, id, true)
func CustomError(ctx context.Context, message, stack, errID string, fromClient bool) {
	errorToGcp(ctx, message, stack, errID, fromClient)
}
