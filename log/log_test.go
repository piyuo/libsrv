package log

import (
	"context"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAiFromContext(t *testing.T) {
	Convey("should get application from context'", t, func() {
		ctx := context.Background()
		backupPiyuoApp := os.Getenv("PIYUO_APP")
		os.Setenv("PIYUO_APP", "dev")
		application, identity := aiFromContext(ctx)
		So(application, ShouldEqual, "dev")
		So(identity, ShouldEqual, "")
		os.Setenv("PIYUO_APP", backupPiyuoApp)
	})
}

func TestLogHead(t *testing.T) {
	Convey("should get head from application and identity'", t, func() {
		HERE := "log_test"
		h := head("piyuo-m-us-sys", "user-store", HERE)
		So(h, ShouldEqual, "piyuo-m-us-sys.user-store.log_test: ")
	})
}
