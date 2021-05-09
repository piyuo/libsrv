package log

import (
	"context"
	"fmt"
	"strings"

	"github.com/piyuo/libsrv/gerror"
	"github.com/piyuo/libsrv/log/logger"
	"github.com/pkg/errors"
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

// prepare write message to history and return information logger need
//
//	message, fields := prepare(ctx, format, a...)
//
func initMessage(ctx context.Context, format string, a ...interface{}) string {
	message := fmt.Sprintf(format, a...)
	if history != nil {
		history.WriteString(message + "\n")
	}
	return message
}

// Debug only print message when os.Getenv("DEBUG") is define
//
//	Debug(ctx,"server start")
//
func Debug(ctx context.Context, format string, a ...interface{}) {
	if testMode != nil {
		return
	}
	if ctx.Err() != nil { // deadline error
		return
	}
	logger.Debug(ctx, initMessage(ctx, format, a...))
}

// Info as Normal but significant events, such as start up, shut down, or a configuration change.
//
//	Info(ctx,"server start")
//
func Info(ctx context.Context, format string, a ...interface{}) {
	if testMode != nil {
		return
	}
	if ctx.Err() != nil { // deadline error
		return
	}
	logger.Info(ctx, initMessage(ctx, format, a...))
}

// Warn as Warning events might cause problems.
//
//	Warning(ctx,"hi")
//
func Warn(ctx context.Context, format string, a ...interface{}) {
	if testMode != nil {
		return
	}
	if ctx.Err() != nil { // deadline error
		return
	}
	logger.Warn(ctx, initMessage(ctx, format, a...))
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
	if history != nil {
		history.Reset()
	}
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

// Error log error to google cloud and return error id, return empty if error not logged
//
//	stack format like
//	at firstLine (a.js:3)
//	at secondLine (b.js:3)
//
//	Error(ctx, err)
//
func Error(ctx context.Context, err error) {
	if ctx.Err() != nil { // deadline error
		return
	}
	if err == nil {
		return
	}
	if testMode != nil {
		return
	}
	stack := beautyStack(err)
	logger.Error(ctx, stack)

	message := beautyMessage(err)
	gerror.Write(ctx, message, stack)
}

// Error log error to google cloud and return error id, return empty if error not logged
//
//	stack format like
//	at firstLine (a.js:3)
//	at secondLine (b.js:3)
//
//	Error(ctx, err)
//
func ErrorToStr(err error) string {
	stack := beautyStack(err)
	return stack
}

// CustomError write error and stack direct to database
//
//	stack format like
//	at firstLine (a.js:3)
//	at secondLine (b.js:3)
//
//	CustomError(ctx, "hi error", stack)
//
func CustomError(ctx context.Context, message, stack string) {
	gerror.Write(ctx, message, stack)
}

// beautyMessage return simple message
//
//	message := beautyMessage(err)
//
func beautyMessage(err error) string {
	cause := errors.Cause(err)
	if cause != nil {
		return cause.Error()
	}
	return err.Error()
}

// beautyStack return simple format stack trace
//
//	stack := beautyStack(err)
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
				sb.WriteString(line + "\n") //this is message
			}
		}
	}
	return strings.Trim(sb.String(), "\n")
}

// isLineUsable check line to see if we need it for debug
//
//	line := "/jtolds/doc.go:75"
//	usable = isLineUsable(line) //false
//
func isLineUsable(line string) bool {
	notUsableKeywords := []string{"jtolds", "log.go", "log_gcp.go", "net/http", "runtime.goexit", "testing.tRunner"}
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
//	duplicate := isLineDuplicate(list, 1) // true
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
//	filename := extractFileName(path) // "doc.go:75"
//
func extractFilename(path string) string {
	index := strings.LastIndex(path, "/")
	if index != -1 {
		return path[index+1:]
	}
	return path
}
