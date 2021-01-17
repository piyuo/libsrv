package data

import (
	"context"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	//	setup()
	code := m.Run()
	//	shutdown()
	os.Exit(code)
}

func TestCleanTest(t *testing.T) {
	Convey("should clean global database", t, func() {
		ctx := context.Background()
		g, err := NewSampleGlobalDB(ctx)
		So(err, ShouldBeNil)
		defer g.Close()

		counters := g.Counters()
		counter := counters.SampleCounter()
		counter.Clear(ctx)

		coders := g.Coders()
		coder := coders.SampleCoder()
		coder.Clear(ctx)

		g.SampleTable().Clear(ctx)
	})

}
