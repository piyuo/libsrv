package server

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/src/log"
)

// deadlineHTTP cache os env DEADLINE_HTTP value
//
var deadlineHTTP time.Duration = -1

// setDeadlineHTTP set context deadline using os.Getenv("DEADLINE_HTTP"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setDeadlineHTTP(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadlineHTTP == -1 {
		text := os.Getenv("DEADLINE_HTTP")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 30000
			log.Warn(ctx, "use default 30 seconds for DEADLINE_HTTP")
		}
		deadlineHTTP = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(deadlineHTTP)
	return context.WithDeadline(ctx, expired)
}

// HTTPCreateFunc create http handler function
//
func HTTPCreateFunc(httpHandler HTTPHandler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		//add deadline to context
		ctx, cancel := setDeadlineHTTP(r.Context())
		defer cancel()

		err := httpHandler(ctx, w, r)
		if err != nil {
			handleRouteException(ctx, w, err)
			return
		}
	}
	withoutArchive := http.HandlerFunc(f)
	withArchive := ArchiveHandler(withoutArchive)
	return withArchive
}
