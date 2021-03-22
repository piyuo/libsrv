package server

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/log"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	log.TestModeAlwaySuccess()
}

func shutdown() {
	log.TestModeBackNormal()
}

func BenchmarkClean(b *testing.B) {
	ctx := context.Background()
	client, _ := newClient(ctx)
	defer client.Close()
	client.DeleteAll(ctx, &TaskLock{}, 100)
}
