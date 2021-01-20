package server

import (
	"net/http"
)

// taskHandler create handler for task
//
func (s *Server) taskHandler() http.Handler {
	return http.HandlerFunc(s.taskServe)
}

// cmdServe serve command request, it filter empty and bad request and send correct one to dispatch
//
//	Cross origin access enabled
//
func (s *Server) taskServe(w http.ResponseWriter, r *http.Request) {

	//add deadline to context
	ctx, cancel := setDeadline(r.Context())
	defer cancel()

	err := s.TaskHandler(ctx, w, r)
	if err != nil {
		handleRouteException(ctx, w, err)
		return
	}
}
