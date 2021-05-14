package log

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	Debug(ctx, "my debug")
	Info(ctx, "my info")
	Warn(ctx, "my warn")

	ForceStopLog(true)
	Debug(ctx, "my debug")
	Info(ctx, "my info")
	Warn(ctx, "my warn")
}

func TestLogContextCanceled(t *testing.T) {
	t.Parallel()
	dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), dateline)
	defer cancel()
	time.Sleep(time.Duration(2) * time.Millisecond)

	Debug(ctx, "cancel debug")
	Info(ctx, "cancel info")
	Warn(ctx, "cancel warn")
	Error(ctx, errors.New("my error"))
}

func TestLogBeautyStack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	err := errors.New("beauty stack")
	stack := beautyStack(err)
	assert.NotEmpty(stack)
}

func TestLogExtractFilename(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	path := "/convey/doc.go:75"
	filename := extractFilename(path)
	assert.Equal("doc.go:75", filename)
	path = "doc.go:75"
	filename = extractFilename(path)
	assert.Equal("doc.go:75", filename)
}

func TestLogIsLineUsable(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	line := "/jtolds/doc.go:75"
	assert.False(isLineUsable(line))
}

func TestLogIsLineDuplicate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	list := []string{"/doc.go:75", "/doc.go:75"}
	assert.False(isLineDuplicate(list, 0))
	assert.True(isLineDuplicate(list, 1))
}

func TestLogError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	err := errors.New("my error")
	Error(ctx, err)

	// nil error
	Error(ctx, nil)

	ForceStopLog(true)
	Error(ctx, err)
}

func TestLogErrorToStr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	err := errors.New("my error")
	str := ErrorToStr(err)
	assert.NotEmpty(str)
}

func TestLogErrorWithRequest(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	err := errors.New("my error with request")
	Error(ctx, err)
}

func TestLogCustomError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	message := "mock error happening in flutter"
	stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
	CustomError(ctx, message, stack)
}

func TestLogHistory(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	assert := assert.New(t)

	KeepHistory(true)
	Info(ctx, "hi")
	assert.True(forceStopLog || strings.Contains(History(), "hi"))

	ResetHistory()
	assert.NotContains(History(), "hi")

	KeepHistory(false)
	Info(ctx, "hi")
	assert.Empty(History())
}

func TestLogErrorWrap(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	err := bar(ctx, 1)
	Error(ctx, err)
	Info(ctx, "server start")
}

func bar(ctx context.Context, n int) (err error) {
	err = errors.New("network error")
	err = errors.Wrapf(err, "wrap(%d)", n)
	if err != nil {
		return err
	}
	return nil
}
