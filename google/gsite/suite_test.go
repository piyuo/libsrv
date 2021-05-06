package gsite

import (
	"os"
	"testing"

	"github.com/piyuo/libsrv/google/gaccount"
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
