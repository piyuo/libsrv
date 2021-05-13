package log

import (
	"os"
	"testing"

	"github.com/piyuo/libsrv/gaccount"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	gaccount.ForceTestCredential(true)
}

func shutdown() {
	gaccount.ForceTestCredential(false)
}
