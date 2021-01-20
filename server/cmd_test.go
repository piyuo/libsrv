package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/piyuo/libsrv/command/mock"
	"github.com/stretchr/testify/assert"
)

func TestEmptyRequestWillReturnBadRequest(t *testing.T) {
	assert := assert.New(t)
	server := &Server{
		Map: &mock.MapXXX{},
	}
	port, handler := server.prepare()
	assert.Equal(":8080", port)

	req1, _ := http.NewRequest("GET", "/", nil)
	req1.Header.Set("Accept-Encoding", "gzip")
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusBadRequest, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestCmdDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("cmdDeadline")
	os.Setenv("cmdDeadline", "20")
	cmdDeadline = -1 // remove cache

	ctx, cancel := setCmdDeadline(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	cmdDeadline = -1 // remove cache
	os.Setenv("cmdDeadline", backup)
}

func TestCmdDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("cmdDeadline")
	os.Setenv("cmdDeadline", "")
	cmdDeadline = -1 // remove cache
	ctx, cancel := setCmdDeadline(ctx)
	defer cancel()

	time.Sleep(time.Duration(21) * time.Millisecond)
	assert.Nil(ctx.Err()) // default expired is in 20,000ms
	cmdDeadline = -1      // remove cache
	os.Setenv("cmdDeadline", backup)
}
