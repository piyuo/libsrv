package gdb

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/gaccount"
	"github.com/piyuo/libsrv/log"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	gaccount.UseTestCredential(true)
	log.TestModeAlwaySuccess()
}

func shutdown() {
	gaccount.UseTestCredential(false)
	log.TestModeBackNormal()
}

func BenchmarkClean(b *testing.B) {
	gaccount.UseTestCredential(true)
	defer gaccount.UseTestCredential(false)

	ctx := context.Background()
	client := sampleClient()
	client.Truncate(ctx, "Sample")
	client.Truncate(ctx, "Count")
	client.Truncate(ctx, "Code")
	client.Truncate(ctx, "Serial")
}
