package gsite

import (
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
	log.ForceStopLog(true)
}

func shutdown() {
	gaccount.ForceTestCredential(false)
	log.ForceStopLog(false)
}
