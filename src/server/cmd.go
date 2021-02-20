package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/src/command"
	"github.com/piyuo/libsrv/src/env"
	"github.com/piyuo/libsrv/src/log"
)

// deadlineCMD cache os env DEADLINE_CMD value
//
var deadlineCMD time.Duration = -1

// setDeadlineCMD set context deadline using os.Getenv("DEADLINE_CMD"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setDeadlineCMD(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadlineCMD == -1 {
		text := os.Getenv("DEADLINE_CMD")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			log.Debug(ctx, "cmd", "use default 20 seconds for DEADLINE_CMD")
		}
		deadlineCMD = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(deadlineCMD)
	return context.WithDeadline(ctx, expired)
}

// CMDCreateFunc create command handler function
//
func CMDCreateFunc(cmdMap command.IMap) http.Handler {
	dispatch := &command.Dispatch{
		Map: cmdMap,
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		//add deadline to context
		ctx, cancel := setDeadlineCMD(r.Context())
		defer cancel()

		//add request to context
		ctx = env.SetRequest(ctx, r)

		if r.Body == nil {
			WriteStatus(w, http.StatusBadRequest, "no request")
			return
		}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(ctx, here, err)
			WriteStatus(w, http.StatusBadRequest, "failed to read request")
			return
		}
		if len(bytes) == 0 {
			WriteStatus(w, http.StatusBadRequest, "bad request")
			return
		}

		bytes, err = dispatch.Route(ctx, bytes)
		if err != nil {
			handleRouteException(ctx, w, err)
			return
		}
		WriteBinary(w, bytes)
	}

	withoutArchive := http.HandlerFunc(f)
	withArchive := ArchiveHandler(withoutArchive)
	return withArchive
}
