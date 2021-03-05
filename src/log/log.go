package log

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/piyuo/libsrv/src/env"
	"github.com/piyuo/libsrv/src/identifier"
)

// history keep all printed log. !Be careful this is global history include all thread log
//
var history *strings.Builder

// testMode is true should return success, false return error, otherwise behave normal
//
var testMode *bool

// TestModeAlwaySuccess will let every function success
//
func TestModeAlwaySuccess() {
	t := true
	testMode = &t
}

// TestModeAlwayFail will let every function fail
//
func TestModeAlwayFail() {
	f := false
	testMode = &f
}

// TestModeBackNormal stop test mode and back to normal
//
func TestModeBackNormal() {
	testMode = nil
}

// Level define log level
//
type Level int

const (
	// DEBUG print to console
	//
	DEBUG Level = iota

	// INFO Normal but significant events, such as start up, shut down, or a configuration change.
	//
	INFO

	// WARNING events might cause problems.
	//
	WARNING

	// ALERT A person must take an action immediately
	//
	ALERT
)

var appName = os.Getenv("NAME")

// getHeader return log header and user id for log
//
//	header,id := getHeader(ctx,"mail") // user-store@piyuo-m-us-sys/mail:,user-store
//
func getHeader(ctx context.Context, where string) (string, string) {
	header := getLocation(where)
	id := env.GetUserID(ctx)
	if id != "" {
		header = id + "@" + header
	}
	return header + ": ", id
}

// getLocation return the location where log happen
//
//	loc := getLocation("mail") // piyuo-m-us-sys/mail
//
func getLocation(where string) string {
	return appName + "/" + where
}

// Print to local console, only check this log in application console log
//
//	here := "log_test"
//	Print(ctx,here,"hello")
//
func Print(ctx context.Context, where, format string, a ...interface{}) {
	header, _ := getHeader(ctx, where)
	message := fmt.Sprintf(format, a...)
	fmt.Printf("%v%v\n", header, message)
	if history != nil {
		history.WriteString(fmt.Sprintf("%v%v\n", header, message))
	}
}

// KeepHistory keep all printed log into history
//
//	KeepHistory(true)
//
func KeepHistory(flag bool) {
	if flag {
		history = &strings.Builder{}
		return
	}
	history = nil
}

// ResetHistory reset history
//
//	ResetHistory()
//
func ResetHistory() {
	history.Reset()
}

// History get log history in string
//
//	History()
//
func History() string {
	if history != nil {
		return history.String()
	}
	return ""
}

// Info as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	HERE := "log_test"
//	Info(ctx,HERE,"hi")
//
func Info(ctx context.Context, where, message string) {
	Log(ctx, INFO, where, message)
}

// Warning as Warning events might cause problems.
//
//	HERE := "log_test"
//	Warning(ctx,HERE,"hi")
//
func Warning(ctx context.Context, where, message string) {
	Log(ctx, WARNING, where, message)
}

// Alert A person must take an action immediately
//
//	HERE := "log_test"
//	Critical(ctx,HERE,"hi")
//
func Alert(ctx context.Context, where, message string) {
	Log(ctx, ALERT, where, message)
}

// Log message and level to server
//
//	here := "log_test"
//	Log(ctx, "hello", here, WARNING)
//
func Log(ctx context.Context, level Level, where, message string) {
	if testMode != nil {
		Print(ctx, where, "%v - [TestMode]", message)
		return
	}
	if ctx.Err() != nil { // deadline error
		return
	}
	logger, err := NewLogger(ctx)
	if err != nil {
		return
	}
	defer logger.Close()
	logger.Write(ctx, level, where, message)
}

// NewLogger return logger
//
//	logger, err := NewLogger(ctx)
//	if err != nil {
//		return
//	}
//	defer logger.Close()
//
func NewLogger(ctx context.Context) (Logger, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return NewGCPLogger(ctx)
}

// WriteError write error and stack to server
//
//	stack format like
//
//	at firstLine (a.js:3)
//
//	at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	errID := identifier.UUID()
//	here := "log_test"
//	LogError(ctx,here, "hi error", stack, errID)
//
func WriteError(ctx context.Context, where, message, stack, errID string) {
	errorer, err := NewErrorer(ctx)
	if err != nil {
		return
	}
	defer errorer.Close()
	errorer.Write(ctx, where, message, stack, errID)
}

// Error log error to google cloud and return error id, return empty if error not logged
//
//	stack format like
//
//	at firstLine (a.js:3)
//
//	at secondLine (b.js:3)
//
//	err := errors.New("my error1")
//	errID := identifier.UUID()
//	HERE := "log_test"
//	LogErr(ctx,HERE, err)
//
func Error(ctx context.Context, where string, err error) {
	if ctx.Err() != nil { // deadline error
		return
	}
	if err == nil {
		return
	}
	if testMode != nil {
		Print(ctx, where, "%v - [TestMode]", err.Error())
		return
	}
	message := err.Error()
	stack := beautyStack(err)
	errID := identifier.UUID()
	WriteError(ctx, where, message, stack, errID)
}

// NewErrorer return errorer
//
//	errorer, err := NewErrorer(ctx)
//
func NewErrorer(ctx context.Context) (Errorer, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return NewGCPErrorer(ctx)
}

// beautyStack return simple format stack trace
//
//	formatedStackFromError(err)
//
func beautyStack(err error) string {
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

// isLineUsable check line to see if we need it for debug
//
//	line := "/convey/doc.go:75"
//	So(isLineUsable(line), ShouldBeFalse)
//
func isLineUsable(line string) bool {
	notUsableKeywords := []string{"smartystreets", "jtolds", "log.go", "log_gcp.go", "net/http", "runtime.goexit", "testing.tRunner"}
	for _, keyword := range notUsableKeywords {
		if strings.Contains(line, keyword) {
			return false
		}
	}
	return true
}

// isLineDuplicate check current line to see if it duplicate in list
//
//	list := []string{"/doc.go:75", "/doc.go:75"}
//	So(isLineDuplicate(list, 1), ShouldBeTrue)
//
func isLineDuplicate(list []string, currentIndex int) bool {
	line := list[currentIndex]
	for index := currentIndex - 1; index >= 0; index-- {
		if line == list[index] {
			return true
		}
	}
	return false
}

// extractFilename extract filename from path
//
// 	path := "/@v1.6.4/doc.go:75"
//	filename := extractFileName(path)
//	So(filename, ShouldEqual, "doc.go:75")
//
func extractFilename(path string) string {
	index := strings.LastIndex(path, "/")
	if index != -1 {
		return path[index+1:]
	}
	return path
}
