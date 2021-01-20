package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
}
