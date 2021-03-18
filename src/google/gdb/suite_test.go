package gdb

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/google/gaccount"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	gaccount.UseTestCredential(true)
}

func shutdown() {
	gaccount.UseTestCredential(false)
}

func BenchmarkGdbClean(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	client.Clear(ctx, &Sample{}, 100)
}
