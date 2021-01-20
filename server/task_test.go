package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
}
