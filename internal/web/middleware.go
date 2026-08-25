package web

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"sync/atomic"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}
	written, err := r.ResponseWriter.Write(data)
	r.bytes += written
	return written, err
}

var requestSequence atomic.Uint64

func (s *Server) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		identifier := request.Header.Get("X-Request-ID")
		if identifier == "" {
			identifier = fmt.Sprintf("request-%08d", requestSequence.Add(1))
		}
		writer.Header().Set("X-Request-ID", identifier)
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("Referrer-Policy", "same-origin")
		next.ServeHTTP(writer, request)
	})
}

func (s *Server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				_ = debug.Stack()
				writeJSON(writer, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			}
		}()
		recorder := &statusRecorder{ResponseWriter: writer}
		next.ServeHTTP(recorder, request)
	})
}
