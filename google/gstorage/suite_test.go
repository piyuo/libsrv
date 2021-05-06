package gstorage

import (
	"os"
	"testing"

	"github.com/piyuo/libsrv/google/gaccount"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	gaccount.UseTestCredential(true)
}

func shutdown() {
	gaccount.UseTestCredential(false)
}
