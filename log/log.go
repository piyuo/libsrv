package log

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	app "github.com/piyuo/go-libsrv/app"
	tools "github.com/piyuo/go-libsrv/tools"
)

//Logger interface
// server: [piyuo-m-us-sys] store-user: hello
// client: <piyuo-m-us-web-page> store-user: hello

//Log level
const (
	LevelDebug   int32 = 0 //debug info
	LevelInfo    int32 = 1 //Normal but significant events, such as start up, shut down, or a configuration change.
	LevelWarning int32 = 2 //Warning events might cause problems.
	LevelAlert   int32 = 3 //A person must take an action immediately
)

//Debug as Routine information, such as ongoing status or performance.
//
//	HERE := "log_test"
//	Debug(ctx,HERE,"hello")
func Debug(ctx context.Context, where, message string) {
	application, identity := aiFromContext(ctx)
	h := head(application, identity, where)
	if app.IsDebug() {
		fmt.Printf("\u001b[34m%v\u001b[0m%v\n", h, message)
	} else {
		fmt.Printf("%v%v\n", h, message)
	}
}

//Info as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	HERE := "log_test"
//	Info(ctx,HERE,"hi")
func Info(ctx context.Context, where, message string) {
	if ctx.Err() != nil {
		return
	}
	application, identity := aiFromContext(ctx)
	Log(ctx, message, application, identity, where, LevelInfo)
}

//Warning as Warning events might cause problems.
//
//	HERE := "log_test"
//	Warning(ctx,HERE,"hi")
func Warning(ctx context.Context, where, message string) {
	if ctx.Err() != nil {
		return
	}
	application, identity := aiFromContext(ctx)
	Log(ctx, message, application, identity, where, LevelWarning)
}

//Alert A person must take an action immediately
//
//	HERE := "log_test"
//	Critical(ctx,HERE,"hi")
func Alert(ctx context.Context, where, message string) {
	if ctx.Err() != nil {
		return
	}
	application, identity := aiFromContext(ctx)
	Log(ctx, message, application, identity, where, LevelAlert)
}

//Error log error to google cloud and return error id
//
//	err := errors.New("my error1")
//	LogErr(ctx, err)
//Error log error to google cloud and return error id, return empty if error not logged
//
//	err := errors.New("my error1")
//	HERE := "log_test"
//	LogErr(ctx,HERE, err)
func Error(ctx context.Context, where string, err error, r *http.Request) string {
	if ctx.Err() != nil {
		return ""
	}
	errID := tools.UUID()
	application, identity := aiFromContext(ctx)
	h := head(application, identity, where)
	message := err.Error()
	stack := beautyStack(err)
	if app.IsDebug() {
		fmt.Printf("\u001b[34m%v\u001b[31m%v \u001b[35m(%v)\n\u001b[33m%v\n", h, err.Error(), errID, stack)
	} else {
		fmt.Printf("%v%v (%v)\n%v\n", h, err.Error(), errID, stack)
	}
	timer := tools.NewTimer()
	timer.Start()
	ErrorLog(ctx, message, application, identity, where, stack, errID, r)
	ms := int(timer.Stop())
	Debug(ctx, where, fmt.Sprintf("error log took %v ms", int(ms)))
	return errID
}

//Log message and level to server
//
//	here := "log_test"
//	CustomLog(ctx, "hello", "piyuo-m-us-sys", "user-store",here, WARNING)
func Log(ctx context.Context, message, application, identity, where string, level int32) {
	if ctx.Err() != nil {
		return
	}
	logger, close, err := Open(ctx)
	if err != nil {
		return
	}
	defer close()
	Write(ctx, logger, time.Now(), message, application, identity, where, level)
}

//Open log client to do batch log
//
//	logger, close, err := Open(ctx)
func Open(ctx context.Context) (*logging.Logger, func(), error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}
	return gcpLogOpen(ctx)
}

//Write log through client
//
//	Write(ctx, logger,time.Now(), message, application, identity, here, info)
func Write(ctx context.Context, logger *logging.Logger, logtime time.Time, message, application, identity, where string, level int32) {
	if ctx.Err() != nil {
		return
	}
	h := head(application, identity, where)
	fontColor := "\u001b[0m" // reset
	switch level {
	case LevelInfo:
		fontColor = "\u001b[36m" // Cyan
	case LevelWarning:
		fontColor = "\u001b[33m" // yellow
	case LevelAlert:
		fontColor = "\u001b[31m" // red
	}
	if app.IsDebug() {
		fmt.Printf("\u001b[34m%v%v%v \u001b[35m(logged)\n", h, fontColor, message)
	} else {
		fmt.Printf("%v%v%v (logged)\n", h, fontColor, message)
	}
	gcpLogWrite(logger, logtime, message, application, identity, where, level)
}

//ErrorLog log error and stack to server
//
//stack format like
//
//at firstLine (a.js:3)
//
//at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	errID := tools.UUID()
//	here := "log_test"
//	LogError(ctx, "hi error", "piyuo-m-us-sys", "user-store",here, stack, errID)
func ErrorLog(ctx context.Context, message, application, identity, where, stack, errID string, r *http.Request) {
	if ctx.Err() != nil {
		return
	}
	client, close, err := ErrorOpen(ctx)
	if err != nil {
		return
	}
	defer close()
	ErrorWrite(ctx, client, message, application, identity, where, stack, errID, r)
}

//ErrorOpen open error client to do batch log
//
//	client, close, err := ErrorOpen(ctx)
func ErrorOpen(ctx context.Context) (*errorreporting.Client, func(), error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}
	return gcpErrorOpen(ctx)
}

//ErrorWrite log error through client
//
//	ErrorWrite(ctx,client, message, application, identity, here, stack, id, nil)
func ErrorWrite(ctx context.Context, client *errorreporting.Client, message, application, identity, where, stack, errID string, r *http.Request) {
	if ctx.Err() != nil {
		return
	}
	gcpErrorWrite(client, message, application, identity, where, stack, errID, r)
}

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
// user-store@piyuo-m-us-sys/where:
//
//	h,identity := head("piyuo-m-us-sys","user-store","where")
func head(application, identity, where string) string {
	text := application + "/" + where
	if identity != "" {
		text = identity + "@" + text
	}
	return text + ": "
}

//beautyStack return simple format stack trace
//
//	formatedStackFromError(err)
func beautyStack(err error) string {
	//debug.PrintStack()
	var sb strings.Builder
	stack := fmt.Sprintf("%+v", err)
	stackFormated := strings.ReplaceAll(stack, "\n\t", "|")
	lines := strings.Split(stackFormated, "\n")
	for index, line := range lines {
		if isLineUsable(line) && !isLineDuplicate(lines, index) {
			parts := strings.Split(line, "|")
			if len(parts) == 2 {
				filename := extractFilename(parts[1])
				newline := fmt.Sprintf("at %v (%v)\n", parts[0], filename)
				sb.WriteString(newline)
			} else {
				//this is message, just ignore it
			}
		}
	}
	return strings.Trim(sb.String(), "\n")
}

//isLineUsable check line to see if we need it for debug
//
//	line := "/convey/doc.go:75"
//	So(isLineUsable(line), ShouldBeFalse)
func isLineUsable(line string) bool {
	notUsableKeywords := []string{"smartystreets", "jtolds", "log.go", "log_gcp.go", "net/http", "runtime.goexit", "testing.tRunner"}
	for _, keyword := range notUsableKeywords {
		if strings.Contains(line, keyword) {
			return false
		}
	}
	return true
}

//isLineDuplicate check current line to see if it duplicate in list
//
//	list := []string{"/doc.go:75", "/doc.go:75"}
//	So(isLineDuplicate(list, 1), ShouldBeTrue)
func isLineDuplicate(list []string, currentIndex int) bool {
	line := list[currentIndex]
	for index := currentIndex - 1; index >= 0; index-- {
		if line == list[index] {
			return true
		}
	}
	return false
}

//extractFilename extract filename from path
//
// 	path := "/goconvey@v1.6.4/convey/doc.go:75"
//	filename := extractFileName(path)
//	So(filename, ShouldEqual, "doc.go:75")
func extractFilename(path string) string {
	index := strings.LastIndex(path, "/")
	if index != -1 {
		return path[index+1 : len(path)]
	}
	return path
}
