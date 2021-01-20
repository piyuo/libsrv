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

func TestTask(t *testing.T) {
	assert := assert.New(t)
	server := &Server{
		TaskHandler: defaultTaskHandler,
	}
	port, handler := server.prepare()
	assert.Equal(":8080", port)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusOK, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)

}

func TestTaskDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("taskDeadline")
	os.Setenv("taskDeadline", "20")
	taskDeadline = -1 // remove cache

	ctx, cancel := setTaskDeadline(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	taskDeadline = -1 // remove cache
	os.Setenv("taskDeadline", backup)
}

func TestTaskDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("taskDeadline")
	os.Setenv("taskDeadline", "")
	taskDeadline = -1 // remove cache
	ctx, cancel := setTaskDeadline(ctx)
	defer cancel()

	time.Sleep(time.Duration(21) * time.Millisecond)
	assert.Nil(ctx.Err()) // default expired is in 20,000ms
	taskDeadline = -1     // remove cache
	os.Setenv("taskDeadline", backup)
}
