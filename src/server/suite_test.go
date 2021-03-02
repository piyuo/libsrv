package server

import (
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

func TestClean(t *testing.T) {
}
