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
	gaccount.ForceTestCredential(true)
	log.TestModeAlwaySuccess()
}

func shutdown() {
	gaccount.ForceTestCredential(false)
	log.TestModeBackNormal()
}

func BenchmarkClean(b *testing.B) {
	gaccount.ForceTestCredential(true)
	defer gaccount.ForceTestCredential(false)

	ctx := context.Background()
	client := sampleClient()
	client.Truncate(ctx, "Sample")
	client.Truncate(ctx, "Count")
	client.Truncate(ctx, "Code")
	client.Truncate(ctx, "Serial")
}
