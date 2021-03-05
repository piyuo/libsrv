package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockTaskHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	return false, nil
}

func mockTaskErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, errors.New("myError")
}

func TestServerTaskDefaultReturn403(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	TaskCreateFunc(mockTaskHandler).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusOK, res.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskHandlerReturnError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	TaskCreateFunc(mockTaskErrorHandler).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusContinue, res.StatusCode)
	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_TASK")
	os.Setenv("DEADLINE_TASK", "20")
	defer os.Setenv("DEADLINE_TASK", backup)
	deadlineTask = -1 // remove cache

	ctx, cancel := setDeadlineTask(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	deadlineTask = -1 // remove cache
}

func TestServerTaskDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_TASK")
	os.Setenv("DEADLINE_TASK", "")
	defer os.Setenv("DEADLINE_TASK", backup)
	deadlineTask = -1 // remove cache

	ctx, cancel := setDeadlineTask(ctx)
	defer cancel()

	ms := deadlineTask.Milliseconds()
	assert.Equal(int64(840000), ms)
	deadlineTask = -1 // remove cache
}
