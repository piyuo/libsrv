package server

import (
	"compress/gzip"
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/piyuo/libsrv/log"
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

// HTTPEntry create http handler function
//
func HTTPEntry(httpHandler HTTPHandler) http.Handler {
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
	//withArchive := ArchiveHandler(withoutArchive)
	withArchive := gzipHandler(withoutArchive)

	//	withArchive := Gzip(withoutArchive)
	return withArchive
}

func gzipHandler(h http.Handler) http.Handler {
	wrapper, _ := gziphandler.NewGzipLevelAndMinSize(gzip.DefaultCompression, 150)
	return wrapper(h)
}

/*
// Gzip Compression
var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(ioutil.Discard)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		var b bytes.Buffer
		gz.Reset(&b)
		defer func() {
			gz.Close()
			l := len(b.Bytes())
			w.Header().Set("Content-Length", fmt.Sprint(l))
			w.Write(b.Bytes())
		}()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
*/
