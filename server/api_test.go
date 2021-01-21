package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAPIDefaultReturn403(t *testing.T) {
	assert := assert.New(t)
	server := &Server{
		APIHandler: defaultHandler,
	}
	port := server.prepare()
	assert.Equal(":8080", port)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	server.createAPIHandler().ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusForbidden, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)

}

func TestApiDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("apiDeadline")
	os.Setenv("apiDeadline", "20")
	apiDeadline = -1 // remove cache

	ctx, cancel := setAPIDeadline(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	apiDeadline = -1 // remove cache
	os.Setenv("apiDeadline", backup)
}

func TestApiDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("apiDeadline")
	os.Setenv("apiDeadline", "")
	apiDeadline = -1 // remove cache
	ctx, cancel := setAPIDeadline(ctx)
	defer cancel()

	time.Sleep(time.Duration(21) * time.Millisecond)
	assert.Nil(ctx.Err()) // default expired is in 20,000ms
	apiDeadline = -1      // remove cache
	os.Setenv("apiDeadline", backup)
}
