package gdb

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/piyuo/libsrv/src/log"
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

func BenchmarkGdbClean(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	client.(*ClientFirestore).deleteAll(ctx, &Sample{}, 100)

}
