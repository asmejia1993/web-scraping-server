package handler

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	body        []byte
	wroteHeader bool
	wroteBody   bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteBody {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.wroteBody {
		return 0, nil
	}
	i, err := rw.ResponseWriter.Write(body)
	if err != nil {
		return 0, err
	}
	rw.body = make([]byte, len(body))
	copy(rw.body, body)
	rw.wroteBody = true

	return i, err
}

func (rw *responseWriter) Body() []byte {
	return rw.body
}

func (hf handlerFranchises) MiddlewareLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			requestBody, err := hf.readRequestBody(r)
			if err != nil {
				hf.respond(w, err, 0)
				return
			}
			hf.restoreRequestBody(r, requestBody)

			logMessage := fmt.Sprintf("path:%s, method: %s", r.URL.EscapedPath(), r.Method)

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(wrapped, r)

			logMessage = fmt.Sprintf("%s, status: %d", logMessage, wrapped.Status())

			hf.logger.Infof("%s, duration: %v", logMessage, time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}
