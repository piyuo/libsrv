package log

import (
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/gaccount"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	gaccount.TestMode(true)
}

func shutdown() {
	gaccount.TestMode(false)
}
