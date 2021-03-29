package gtask

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/piyuo/libsrv/src/google/gdb"
	"github.com/piyuo/libsrv/src/log"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	//	gaccount.UseTestCredential(true)
	log.TestModeAlwaySuccess()
}

func shutdown() {
	//	gaccount.UseTestCredential(false)
	log.TestModeBackNormal()
}

func BenchmarkClean(b *testing.B) {
	ctx := context.Background()
	cred, _ := gaccount.GlobalCredential(ctx)
	client, _ := gdb.NewClient(ctx, cred)
	defer client.Close()
	client.Truncate(ctx, "Task", 100)
}
