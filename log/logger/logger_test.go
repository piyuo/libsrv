package logger

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/env"
	"github.com/stretchr/testify/assert"
)

func TestLoggerInitial(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	l := initial()
	assert.NotNil(l)
}

func TestLoggerInitialDebug(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	bak := env.Debug
	env.Debug = true
	l := initial()
	assert.NotNil(l)
	env.Debug = bak
}

func TestLoggerAddContextInformation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	ctx = env.SetUserID(ctx, "user1")
	fields := addContextInformation(ctx)
	assert.NotEmpty(fields)
}

func TestLoggerDebug(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	Debug(ctx, "logger debug")
}

func TestLoggerInfo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	Info(ctx, "logger info")
}

func TestLoggerWarn(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	Warn(ctx, "logger warn")
}

func TestLoggerError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	Error(ctx, "logger error")
}
