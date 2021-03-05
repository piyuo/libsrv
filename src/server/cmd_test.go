package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/command/mock"
	"github.com/stretchr/testify/assert"
)

func TestServerNilBodyWillReturnBadRequest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	CMDCreateFunc(&mock.MapXXX{}).ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusBadRequest, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerEmptyRequestWillReturnBadRequest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req1, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	resp1 := httptest.NewRecorder()
	CMDCreateFunc(&mock.MapXXX{}).ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusBadRequest, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestCmdDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_CMD")
	os.Setenv("DEADLINE_CMD", "20")
	defer os.Setenv("DEADLINE_CMD", backup)
	deadlineCMD = -1 // remove cache

	ctx, cancel := setDeadlineCMD(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	deadlineCMD = -1 // remove cache
}

func TestCmdDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_CMD")
	os.Setenv("DEADLINE_CMD", "")
	defer os.Setenv("DEADLINE_CMD", backup)
	deadlineCMD = -1 // remove cache

	ctx, cancel := setDeadlineCMD(ctx)
	defer cancel()

	ms := deadlineCMD.Milliseconds()
	assert.Equal(int64(20000), ms)
	deadlineCMD = -1 // remove cache
}
