package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/piyuo/libsrv/command/mock"
	"github.com/stretchr/testify/assert"
)

func TestNilBodyWillReturnBadRequest(t *testing.T) {
	assert := assert.New(t)
	server := &Server{
		Map: &mock.MapXXX{},
	}
	port := server.prepare()
	assert.Equal(":8080", port)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	server.createCmdHandler().ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusBadRequest, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestEmptyRequestWillReturnBadRequest(t *testing.T) {
	assert := assert.New(t)
	server := &Server{
		Map: &mock.MapXXX{},
	}
	port := server.prepare()
	assert.Equal(":8080", port)

	req1, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	resp1 := httptest.NewRecorder()
	server.createCmdHandler().ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusBadRequest, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestCmdDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE")
	os.Setenv("DEADLINE", "20")
	defer os.Setenv("DEADLINE", backup)
	deadline = -1 // remove cache

	ctx, cancel := setDeadline(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	deadline = -1 // remove cache
}

func TestCmdDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE")
	os.Setenv("DEADLINE", "")
	defer os.Setenv("DEADLINE", backup)
	deadline = -1 // remove cache

	ctx, cancel := setDeadline(ctx)
	defer cancel()

	time.Sleep(time.Duration(21) * time.Millisecond)
	assert.Nil(ctx.Err()) // default expired is in 20,000ms
	deadline = -1         // remove cache
}
