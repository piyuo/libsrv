package server

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/piyuo/libsrv/src/google/gdb"
	"github.com/piyuo/libsrv/src/log"
)

var sampleClientInstance db.Client

func sampleClient() db.Client {
	if sampleClientInstance != nil {
		return sampleClientInstance
	}
	ctx := context.Background()
	cred, _ := gaccount.GlobalCredential(ctx)
	client, _ := gdb.NewClient(ctx, cred)
	sampleClientInstance = client
	return sampleClientInstance
}

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
	ctx := context.Background()
	client := sampleClient()
	client.Truncate(ctx, "Task", 100)
}
