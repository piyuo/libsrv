package gdb

import (
	"context"
	"os"
	"testing"

	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/stretchr/testify/assert"
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

func TestCleanTest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	table := g.SampleTable()
	table.Clear(ctx)

	counters := g.Counters()
	counter := counters.SampleCounter()
	counter.Clear(ctx)

	coders := g.Coders()
	coder := coders.SampleCoder()
	coder.Clear(ctx)

	g.SampleTable().Clear(ctx)
}
