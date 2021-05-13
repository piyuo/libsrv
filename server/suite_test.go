package server

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/gaccount"
	"github.com/piyuo/libsrv/gdb"
	"github.com/piyuo/libsrv/log"
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
	gaccount.ForceTestCredential(true)
	log.TestModeAlwaySuccess()
}

func shutdown() {
	gaccount.ForceTestCredential(false)
	log.TestModeBackNormal()
}

func BenchmarkClean(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	client.Truncate(ctx, "Task")
}
