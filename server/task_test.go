package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/piyuo/libsrv/gtask"
	"github.com/stretchr/testify/assert"
)

func mockTaskErrorHandler(ctx context.Context, r *http.Request) error {
	return errors.New("myError")
}

func TestServerTaskHandlerOK(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	taskID, err := gtask.New(ctx, "task", "http://not-exists", []byte{}, "TaskHandlerOK", 1800, 3)
	assert.Nil(err)

	req, _ := http.NewRequest("GET", "/?TaskID="+taskID, nil)
	resp := httptest.NewRecorder()
	TaskEntry(func(ctx context.Context, r *http.Request) error {
		return nil
	}).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusOK, res.StatusCode)

	found, err := sampleClient().Exists(ctx, &gtask.Task{}, taskID)
	assert.Nil(err)
	assert.False(found)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskHandlerInProgress(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	setDeadlineTask(ctx)
	taskID, err := gtask.New(ctx, "task", "http://not-exists.com", []byte{}, "TaskHandlerInProgress", 1800, 3)
	assert.Nil(err)
	err = gtask.Lock(ctx, taskID)
	assert.Nil(err)
	defer gtask.Delete(ctx, taskID)

	req, _ := http.NewRequest("GET", "/?TaskID="+taskID, nil)
	resp := httptest.NewRecorder()
	TaskEntry(func(ctx context.Context, r *http.Request) error {
		return nil
	}).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusOK, res.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskHandlerReturnError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	taskID, err := gtask.New(ctx, "task", "http://not-exists.com", []byte{}, "TaskHandlerInProgress", 1800, 3)
	assert.Nil(err)
	defer gtask.Delete(ctx, taskID)

	req, _ := http.NewRequest("GET", "/?TaskID="+taskID, nil)
	resp := httptest.NewRecorder()
	TaskEntry(mockTaskErrorHandler).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusOK, res.StatusCode)
	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskHandlerNoTaskID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	TaskEntry(mockTaskErrorHandler).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusOK, res.StatusCode)
	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskHandlerDebug(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/?debug=1", nil)
	resp := httptest.NewRecorder()
	TaskEntry(func(ctx context.Context, r *http.Request) error {
		return nil
	}).ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(http.StatusOK, res.StatusCode)
	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerTaskDeadline(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_TASK")
	os.Setenv("DEADLINE_TASK", "")
	defer os.Setenv("DEADLINE_TASK", backup)
	deadlineTask = -1 // remove cache

	_, cancel := setDeadlineTask(ctx)
	defer cancel()

	ms := deadlineTask.Milliseconds()
	assert.Greater(ms, int64(0))
	deadlineTask = -1 // remove cache
}
